package cmd

import (
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/ordinal-cli/ordinal/internal/models"
	"github.com/spf13/cobra"
)

var (
	subscriberPostID  string
	subscriberID      string
	subscriberUserIDs string
)

func init() {
	rootCmd.AddCommand(subscriberCmd)
	subscriberCmd.AddCommand(subscriberListCmd)
	subscriberCmd.AddCommand(subscriberCreateCmd)
	subscriberCmd.AddCommand(subscriberDeleteCmd)

	subscriberListCmd.Flags().StringVar(&subscriberPostID, "post-id", "", "Post ID (UUID)")
	subscriberListCmd.MarkFlagRequired("post-id")

	subscriberCreateCmd.Flags().StringVar(&subscriberPostID, "post-id", "", "Post ID (UUID)")
	subscriberCreateCmd.Flags().StringVar(&subscriberUserIDs, "user-ids", "", "Comma-separated user IDs (UUIDs) to subscribe")
	subscriberCreateCmd.MarkFlagRequired("post-id")
	subscriberCreateCmd.MarkFlagRequired("user-ids")

	subscriberDeleteCmd.Flags().StringVar(&subscriberID, "id", "", "Subscriber ID (UUID)")
	subscriberDeleteCmd.MarkFlagRequired("id")
}

var subscriberCmd = &cobra.Command{
	Use:   "subscriber",
	Short: "Manage post subscribers",
}

var subscriberListCmd = &cobra.Command{
	Use:   "list",
	Short: "List subscribers on a post",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		resp, err := api.NewSubscriberService(c).List(subscriberPostID)
		if err != nil {
			return err
		}
		return printResult(resp)
	},
}

var subscriberCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Add subscribers to a post",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewSubscriberService(c).Create(models.CreateSubscribersRequest{
			PostID:  subscriberPostID,
			UserIDs: splitCSV(subscriberUserIDs),
		})
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var subscriberDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove a subscriber",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		if err := api.NewSubscriberService(c).Delete(subscriberID); err != nil {
			return err
		}
		fmt.Printf("Subscriber %s deleted\n", subscriberID)
		return nil
	},
}
