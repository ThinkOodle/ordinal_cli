package cmd

import (
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/ordinal-cli/ordinal/internal/models"
	"github.com/spf13/cobra"
)

var (
	slackBoostPostID         string
	slackBoostID             string
	slackBoostCreateWebhookID string
	slackBoostCreateCopy      string
	slackBoostUpdateCopy      string
	slackBoostUpdateBodyJSON  string
	slackBoostUpdateBodyFile  string
)

func init() {
	rootCmd.AddCommand(slackBoostCmd)
	slackBoostCmd.AddCommand(slackBoostListCmd)
	slackBoostCmd.AddCommand(slackBoostGetCmd)
	slackBoostCmd.AddCommand(slackBoostCreateCmd)
	slackBoostCmd.AddCommand(slackBoostUpdateCmd)
	slackBoostCmd.AddCommand(slackBoostDeleteCmd)

	rootCmd.AddCommand(slackWebhookCmd)
	slackWebhookCmd.AddCommand(slackWebhookListCmd)

	slackBoostListCmd.Flags().StringVar(&slackBoostPostID, "post-id", "", "Post ID (UUID)")
	slackBoostListCmd.MarkFlagRequired("post-id")

	slackBoostGetCmd.Flags().StringVar(&slackBoostID, "id", "", "Slack boost ID (UUID)")
	slackBoostGetCmd.MarkFlagRequired("id")

	slackBoostCreateCmd.Flags().StringVar(&slackBoostPostID, "post-id", "", "Post ID (UUID)")
	slackBoostCreateCmd.Flags().StringVar(&slackBoostCreateWebhookID, "slack-webhook-id", "", "Slack webhook ID (UUID) to send to")
	slackBoostCreateCmd.Flags().StringVar(&slackBoostCreateCopy, "copy", "", "Optional custom message")
	slackBoostCreateCmd.MarkFlagRequired("post-id")
	slackBoostCreateCmd.MarkFlagRequired("slack-webhook-id")

	slackBoostUpdateCmd.Flags().StringVar(&slackBoostID, "id", "", "Slack boost ID (UUID)")
	slackBoostUpdateCmd.Flags().StringVar(&slackBoostUpdateCopy, "copy", "", "Updated custom message")
	slackBoostUpdateCmd.Flags().StringVar(&slackBoostUpdateBodyJSON, "body-json", "", "Full JSON body (individual flags override matching keys when set)")
	slackBoostUpdateCmd.Flags().StringVar(&slackBoostUpdateBodyFile, "body-file", "", "Path to JSON body file (or - for stdin)")
	slackBoostUpdateCmd.MarkFlagRequired("id")

	slackBoostDeleteCmd.Flags().StringVar(&slackBoostID, "id", "", "Slack boost ID (UUID)")
	slackBoostDeleteCmd.MarkFlagRequired("id")
}

var slackBoostCmd = &cobra.Command{
	Use:   "slack-boost",
	Short: "Manage Slack boosts attached to posts",
}

var slackBoostListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Slack boosts on a post",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		resp, err := api.NewSlackBoostService(c).ListByPost(slackBoostPostID)
		if err != nil {
			return err
		}
		return printResult(resp)
	},
}

var slackBoostGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a Slack boost by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		b, err := api.NewSlackBoostService(c).Get(slackBoostID)
		if err != nil {
			return err
		}
		return printResult(b)
	},
}

var slackBoostCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a Slack boost",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		b, err := api.NewSlackBoostService(c).Create(models.CreateSlackBoostRequest{
			PostID:         slackBoostPostID,
			SlackWebhookID: slackBoostCreateWebhookID,
			Copy:           slackBoostCreateCopy,
		})
		if err != nil {
			return err
		}
		return printResult(b)
	},
}

var slackBoostUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a Slack boost",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		body, err := parseBodyJSON(slackBoostUpdateBodyJSON, slackBoostUpdateBodyFile)
		if err != nil {
			return err
		}
		if body == nil {
			body = map[string]interface{}{}
		}
		if cmd.Flags().Changed("copy") {
			body["copy"] = slackBoostUpdateCopy
		}
		if len(body) == 0 {
			return fmt.Errorf("no fields to update; provide flags or --body-json/--body-file")
		}
		b, err := api.NewSlackBoostService(c).Update(slackBoostID, body)
		if err != nil {
			return err
		}
		return printResult(b)
	},
}

var slackBoostDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Slack boost",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		if err := api.NewSlackBoostService(c).Delete(slackBoostID); err != nil {
			return err
		}
		fmt.Printf("Slack boost %s deleted\n", slackBoostID)
		return nil
	},
}

var slackWebhookCmd = &cobra.Command{
	Use:   "slack-webhook",
	Short: "Slack boost channels (connected Slack webhooks)",
}

var slackWebhookListCmd = &cobra.Command{
	Use:   "list",
	Short: "List connected Slack boost channels",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		resp, err := api.NewSlackWebhookService(c).List()
		if err != nil {
			return err
		}
		return printResult(resp)
	},
}
