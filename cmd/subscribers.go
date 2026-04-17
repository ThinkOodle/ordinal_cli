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
		// Cobra's MarkFlagRequired only checks that --user-ids was passed,
		// not that it contains any non-empty entry after trimming. A value
		// like "," or "  ,  " collapses to an empty slice in splitCSV and
		// would otherwise send an empty userIds array to the API.
		userIDs := splitCSV(subscriberUserIDs)
		if len(userIDs) == 0 {
			return fmt.Errorf("--user-ids must contain at least one non-empty UUID")
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewSubscriberService(c).Create(models.CreateSubscribersRequest{
			PostID:  subscriberPostID,
			UserIDs: userIDs,
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
		return printResult(deletedAck("subscriber", subscriberID))
	},
}
