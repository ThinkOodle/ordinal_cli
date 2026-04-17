package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/ordinal-cli/ordinal/internal/models"
	"github.com/spf13/cobra"
)

var (
	engagementPostID         string
	engagementID             string
	engagementCreateChannel  string
	engagementCreateBodyJSON string
	engagementCreateBodyFile string
	engagementUpdateBodyJSON string
	engagementUpdateBodyFile string
)

func init() {
	rootCmd.AddCommand(engagementCmd)
	engagementCmd.AddCommand(engagementListCmd)
	engagementCmd.AddCommand(engagementCreateCmd)
	engagementCmd.AddCommand(engagementUpdateCmd)
	engagementCmd.AddCommand(engagementDeleteCmd)

	engagementListCmd.Flags().StringVar(&engagementPostID, "post-id", "", "Post ID (UUID)")
	engagementListCmd.MarkFlagRequired("post-id")

	engagementCreateCmd.Flags().StringVar(&engagementPostID, "post-id", "", "Post ID (UUID)")
	engagementCreateCmd.Flags().StringVar(&engagementCreateChannel, "channel", "", "Channel (LinkedIn or Twitter)")
	engagementCreateCmd.Flags().StringVar(&engagementCreateBodyJSON, "body-json", "", "JSON array of engagement inputs, e.g. '[{\"type\":\"Like\",\"profileId\":\"...\"}]'")
	engagementCreateCmd.Flags().StringVar(&engagementCreateBodyFile, "body-file", "", "Path to JSON file (or - for stdin) with engagements array")
	engagementCreateCmd.MarkFlagRequired("post-id")
	engagementCreateCmd.MarkFlagRequired("channel")

	engagementUpdateCmd.Flags().StringVar(&engagementID, "id", "", "Engagement ID (UUID)")
	engagementUpdateCmd.Flags().StringVar(&engagementUpdateBodyJSON, "body-json", "", "Inline JSON body with fields to update")
	engagementUpdateCmd.Flags().StringVar(&engagementUpdateBodyFile, "body-file", "", "Path to JSON file (or - for stdin) with update body")
	engagementUpdateCmd.MarkFlagRequired("id")

	engagementDeleteCmd.Flags().StringVar(&engagementID, "id", "", "Engagement ID (UUID)")
	engagementDeleteCmd.MarkFlagRequired("id")
}

var engagementCmd = &cobra.Command{
	Use:   "engagement",
	Short: "Manage auto-engagements (likes, comments, reposts) on posts",
}

var engagementListCmd = &cobra.Command{
	Use:   "list",
	Short: "List engagements on a post",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		resp, err := api.NewEngagementService(c).List(engagementPostID)
		if err != nil {
			return err
		}
		return printResult(resp)
	},
}

var engagementCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create engagements on a post",
	RunE: func(cmd *cobra.Command, args []string) error {
		var raw []byte
		var err2 error
		switch {
		case engagementCreateBodyJSON != "" && engagementCreateBodyFile != "":
			return fmt.Errorf("only one of --body-json and --body-file may be set")
		case engagementCreateBodyJSON != "":
			raw = []byte(engagementCreateBodyJSON)
		case engagementCreateBodyFile != "":
			if raw, err2 = readBodyFile(engagementCreateBodyFile); err2 != nil {
				return err2
			}
		default:
			return fmt.Errorf("provide --body-json or --body-file with an engagements array")
		}

		var engagements []map[string]interface{}
		if err := json.Unmarshal(raw, &engagements); err != nil {
			return fmt.Errorf("parsing engagements array: %w", err)
		}
		// Parity with subscriber/webhook create: a presence-only check on the
		// body flag lets [] through and sends an empty array to the API,
		// which wastes a round-trip for an obviously-invalid request.
		if len(engagements) == 0 {
			return fmt.Errorf("engagements array must contain at least one engagement")
		}

		c, err := newClient()
		if err != nil {
			return err
		}
		req := models.CreateEngagementsRequest{
			Channel:     engagementCreateChannel,
			Engagements: engagements,
		}
		data, err := api.NewEngagementService(c).Create(engagementPostID, req)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var engagementUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an engagement",
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := parseBodyJSON(engagementUpdateBodyJSON, engagementUpdateBodyFile)
		if err != nil {
			return err
		}
		if len(body) == 0 {
			return fmt.Errorf("no fields to update; provide --body-json or --body-file")
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewEngagementService(c).Update(engagementID, body)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var engagementDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an engagement",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		if err := api.NewEngagementService(c).Delete(engagementID); err != nil {
			return err
		}
		return printResult(deletedAck("engagement", engagementID))
	},
}
