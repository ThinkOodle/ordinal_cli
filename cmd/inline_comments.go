package cmd

import (
	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/ordinal-cli/ordinal/internal/models"
	"github.com/spf13/cobra"
)

var (
	inlineCommentPostID   string
	inlineCommentChannel  string
	inlineCommentResolved bool
)

func init() {
	rootCmd.AddCommand(inlineCommentCmd)
	inlineCommentCmd.AddCommand(inlineCommentListCmd)

	inlineCommentListCmd.Flags().StringVar(&inlineCommentPostID, "post-id", "", "Post ID (UUID)")
	inlineCommentListCmd.Flags().StringVar(&inlineCommentChannel, "channel", "", "Filter by channel")
	inlineCommentListCmd.Flags().BoolVar(&inlineCommentResolved, "resolved", false, "Filter by resolution status")
	inlineCommentListCmd.MarkFlagRequired("post-id")
}

var inlineCommentCmd = &cobra.Command{
	Use:   "inline-comment",
	Short: "List inline (text-anchored) comments on a post",
}

var inlineCommentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List inline comments on a post",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		params := models.ListInlineCommentsParams{
			Channel: inlineCommentChannel,
		}
		if cmd.Flags().Changed("resolved") {
			params.Resolved = &inlineCommentResolved
		}
		resp, err := api.NewInlineCommentService(c).List(inlineCommentPostID, params)
		if err != nil {
			return err
		}
		return printResult(resp)
	},
}
