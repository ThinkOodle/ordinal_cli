package cmd

import (
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/ordinal-cli/ordinal/internal/models"
	"github.com/spf13/cobra"
)

var (
	webhookID                string
	webhookCreateName        string
	webhookCreateURL         string
	webhookCreateDescription string
	webhookCreateTopics      string
	webhookCreateHeadersJSON string
	webhookUpdateBodyJSON    string
	webhookUpdateBodyFile    string
)

func init() {
	rootCmd.AddCommand(webhookCmd)
	webhookCmd.AddCommand(webhookListCmd)
	webhookCmd.AddCommand(webhookGetCmd)
	webhookCmd.AddCommand(webhookCreateCmd)
	webhookCmd.AddCommand(webhookUpdateCmd)
	webhookCmd.AddCommand(webhookDeleteCmd)

	webhookGetCmd.Flags().StringVar(&webhookID, "id", "", "Webhook ID (UUID)")
	webhookGetCmd.MarkFlagRequired("id")

	webhookCreateCmd.Flags().StringVar(&webhookCreateName, "name", "", "Webhook name")
	webhookCreateCmd.Flags().StringVar(&webhookCreateURL, "url", "", "URL that will receive webhook events")
	webhookCreateCmd.Flags().StringVar(&webhookCreateDescription, "description", "", "Optional description")
	webhookCreateCmd.Flags().StringVar(&webhookCreateTopics, "topics", "", "Comma-separated list of event topics (e.g. post.created,post.published)")
	webhookCreateCmd.Flags().StringVar(&webhookCreateHeadersJSON, "headers", "", "Optional custom headers as JSON object, e.g. '{\"X-Key\":\"value\"}'")
	webhookCreateCmd.MarkFlagRequired("name")
	webhookCreateCmd.MarkFlagRequired("url")
	webhookCreateCmd.MarkFlagRequired("topics")

	webhookUpdateCmd.Flags().StringVar(&webhookID, "id", "", "Webhook ID (UUID)")
	webhookUpdateCmd.Flags().StringVar(&webhookUpdateBodyJSON, "body-json", "", "Inline JSON body of fields to update")
	webhookUpdateCmd.Flags().StringVar(&webhookUpdateBodyFile, "body-file", "", "Path to JSON file (or - for stdin) with fields to update")
	webhookUpdateCmd.MarkFlagRequired("id")

	webhookDeleteCmd.Flags().StringVar(&webhookID, "id", "", "Webhook ID (UUID)")
	webhookDeleteCmd.MarkFlagRequired("id")
}

var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Manage webhook subscriptions",
}

var webhookListCmd = &cobra.Command{
	Use:   "list",
	Short: "List webhooks",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		webhooks, err := api.NewWebhookService(c).List()
		if err != nil {
			return err
		}
		return printResult(webhooks)
	},
}

var webhookGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a webhook by ID",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		w, err := api.NewWebhookService(c).Get(webhookID)
		if err != nil {
			return err
		}
		return printResult(w)
	},
}

var webhookCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a webhook",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Cobra's MarkFlagRequired only checks that --topics was passed,
		// not that it contains any non-empty entry after trimming. A value
		// like "," or "  ,  " collapses to an empty slice in splitCSV and
		// would otherwise send an empty topics array to the API.
		topics := splitCSV(webhookCreateTopics)
		if len(topics) == 0 {
			return fmt.Errorf("--topics must contain at least one non-empty event topic")
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		req := models.CreateWebhookRequest{
			Name:        webhookCreateName,
			URL:         webhookCreateURL,
			Description: webhookCreateDescription,
			Topics:      topics,
		}
		if webhookCreateHeadersJSON != "" {
			headers, err := parseBodyJSON(webhookCreateHeadersJSON, "")
			if err != nil {
				return fmt.Errorf("parsing --headers: %w", err)
			}
			if len(headers) > 0 {
				req.Headers = map[string]string{}
				for k, v := range headers {
					if s, ok := v.(string); ok {
						req.Headers[k] = s
					} else {
						req.Headers[k] = fmt.Sprintf("%v", v)
					}
				}
			}
		}
		w, err := api.NewWebhookService(c).Create(req)
		if err != nil {
			return err
		}
		return printResult(w)
	},
}

var webhookUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a webhook",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		body, err := parseBodyJSON(webhookUpdateBodyJSON, webhookUpdateBodyFile)
		if err != nil {
			return err
		}
		if body == nil {
			return fmt.Errorf("provide --body-json or --body-file")
		}
		w, err := api.NewWebhookService(c).Update(webhookID, body)
		if err != nil {
			return err
		}
		return printResult(w)
	},
}

var webhookDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a webhook",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		if err := api.NewWebhookService(c).Delete(webhookID); err != nil {
			return err
		}
		return printResult(deletedAck("webhook", webhookID))
	},
}
