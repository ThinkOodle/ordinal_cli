package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/ordinal-cli/ordinal/internal/models"
	"github.com/spf13/cobra"
)

var (
	approvalPostID         string
	approvalID             string
	approvalCreateBodyJSON string
	approvalCreateBodyFile string
	approvalCreateUserIDs  string
	approvalCreateMessage  string
	approvalCreateDueDate  string
	approvalCreateBlocking bool
)

func init() {
	rootCmd.AddCommand(approvalCmd)
	approvalCmd.AddCommand(approvalListCmd)
	approvalCmd.AddCommand(approvalCreateCmd)
	approvalCmd.AddCommand(approvalDeleteCmd)

	approvalListCmd.Flags().StringVar(&approvalPostID, "post-id", "", "Post ID (UUID)")
	approvalListCmd.MarkFlagRequired("post-id")

	approvalCreateCmd.Flags().StringVar(&approvalPostID, "post-id", "", "Post ID (UUID)")
	approvalCreateCmd.Flags().StringVar(&approvalCreateUserIDs, "user-ids", "", "Comma-separated user IDs to request approval from (shortcut; all share the same message/due-date/blocking flag)")
	approvalCreateCmd.Flags().StringVar(&approvalCreateMessage, "message", "", "Optional message for each approval request")
	approvalCreateCmd.Flags().StringVar(&approvalCreateDueDate, "due-date", "", "Optional due date (ISO 8601)")
	approvalCreateCmd.Flags().BoolVar(&approvalCreateBlocking, "blocking", false, "Mark these approvals as blocking")
	approvalCreateCmd.Flags().StringVar(&approvalCreateBodyJSON, "body-json", "", "JSON approvals: either a top-level array or an object with an \"approvals\" array (overrides --user-ids)")
	approvalCreateCmd.Flags().StringVar(&approvalCreateBodyFile, "body-file", "", "Path to JSON file (or - for stdin) with approvals: array or {\"approvals\":[...]}")
	approvalCreateCmd.MarkFlagRequired("post-id")

	approvalDeleteCmd.Flags().StringVar(&approvalID, "id", "", "Approval ID (UUID)")
	approvalDeleteCmd.MarkFlagRequired("id")
}

var approvalCmd = &cobra.Command{
	Use:   "approval",
	Short: "Manage post approvals",
}

var approvalListCmd = &cobra.Command{
	Use:   "list",
	Short: "List approvals for a post",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		resp, err := api.NewApprovalService(c).List(approvalPostID)
		if err != nil {
			return err
		}
		return printResult(resp)
	},
}

var approvalCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create approval requests for a post",
	Long:  "Create one or more approval requests. Use --user-ids to request from multiple users with shared settings, or --body-json/--body-file for full control over each approval entry.",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		req := models.CreateApprovalsRequest{PostID: approvalPostID}

		raw, err := readBodyJSONRaw(approvalCreateBodyJSON, approvalCreateBodyFile)
		if err != nil {
			return err
		}
		if raw != nil {
			if err := parseApprovalsBody(raw, &req.Approvals); err != nil {
				return err
			}
		} else if approvalCreateUserIDs != "" {
			for _, uid := range splitCSV(approvalCreateUserIDs) {
				req.Approvals = append(req.Approvals, models.ApprovalRequestInput{
					UserID:     uid,
					Message:    approvalCreateMessage,
					DueDate:    approvalCreateDueDate,
					IsBlocking: approvalCreateBlocking,
				})
			}
		}

		if len(req.Approvals) == 0 {
			return fmt.Errorf("provide --user-ids or an approvals array via --body-json/--body-file")
		}

		data, err := api.NewApprovalService(c).Create(req)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

// parseApprovalsBody decodes a --body-json/--body-file payload into approvals.
// Accepts either a bare JSON array of approval entries or an object with an
// "approvals" key whose value is the array, as documented in --help.
func parseApprovalsBody(raw []byte, out *[]models.ApprovalRequestInput) error {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return nil
	}

	switch trimmed[0] {
	case '[':
		if err := json.Unmarshal(trimmed, out); err != nil {
			return fmt.Errorf("parsing approvals array: %w", err)
		}
		return nil
	case '{':
		var obj struct {
			Approvals json.RawMessage `json:"approvals"`
		}
		if err := json.Unmarshal(trimmed, &obj); err != nil {
			return fmt.Errorf("parsing approvals body: %w", err)
		}
		if len(obj.Approvals) == 0 {
			return fmt.Errorf(`approvals body object must contain an "approvals" array`)
		}
		if err := json.Unmarshal(obj.Approvals, out); err != nil {
			return fmt.Errorf("parsing approvals body: %w", err)
		}
		return nil
	}
	return fmt.Errorf("approvals body must be a JSON array or an object with an \"approvals\" array")
}

var approvalDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an approval",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		if err := api.NewApprovalService(c).Delete(approvalID); err != nil {
			return err
		}
		fmt.Printf("Approval %s deleted\n", approvalID)
		return nil
	},
}
