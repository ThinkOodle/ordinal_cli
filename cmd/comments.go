package cmd

import (
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/ordinal-cli/ordinal/internal/models"
	"github.com/spf13/cobra"
)

var (
	commentPostID  string
	commentID      string
	commentMessage string
)

func init() {
	rootCmd.AddCommand(commentCmd)
	commentCmd.AddCommand(commentListCmd)
	commentCmd.AddCommand(commentCreateCmd)
	commentCmd.AddCommand(commentDeleteCmd)

	commentListCmd.Flags().StringVar(&commentPostID, "post-id", "", "Post ID (UUID)")
	commentListCmd.MarkFlagRequired("post-id")

	commentCreateCmd.Flags().StringVar(&commentPostID, "post-id", "", "Post ID (UUID)")
	commentCreateCmd.Flags().StringVar(&commentMessage, "message", "", "Comment message (supports @[Display Name](userId) mentions)")
	commentCreateCmd.MarkFlagRequired("post-id")
	commentCreateCmd.MarkFlagRequired("message")

	commentDeleteCmd.Flags().StringVar(&commentID, "id", "", "Comment ID (UUID)")
	commentDeleteCmd.MarkFlagRequired("id")
}

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "Manage post comments",
}

var commentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List comments on a post",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		resp, err := api.NewCommentService(c).List(commentPostID)
		if err != nil {
			return err
		}
		return printResult(resp)
	},
}

var commentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a comment on a post",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		comment, err := api.NewCommentService(c).Create(commentPostID, models.CreateCommentRequest{Message: commentMessage})
		if err != nil {
			return err
		}
		return printResult(comment)
	},
}

var commentDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a comment (author only)",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		if err := api.NewCommentService(c).Delete(commentID); err != nil {
			return err
		}
		fmt.Printf("Comment %s deleted\n", commentID)
		return nil
	},
}
