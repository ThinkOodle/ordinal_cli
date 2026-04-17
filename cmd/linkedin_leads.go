package cmd

import (
	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/spf13/cobra"
)

var (
	llProfileID        string
	llPostID           string
	llStartDate        string
	llEndDate          string
	llLimit            int
	llCursor           string
	llTypes            string
	llMinFollowerCount int
)

func init() {
	rootCmd.AddCommand(linkedInLeadsCmd)
	linkedInLeadsCmd.AddCommand(linkedInLeadsListPostsCmd)
	linkedInLeadsCmd.AddCommand(linkedInLeadsGetLeadsCmd)

	linkedInLeadsListPostsCmd.Flags().StringVar(&llProfileID, "profile-id", "", "LinkedIn profile ID (UUID)")
	linkedInLeadsListPostsCmd.Flags().StringVar(&llStartDate, "start-date", "", "Filter posts published on or after this date")
	linkedInLeadsListPostsCmd.Flags().StringVar(&llEndDate, "end-date", "", "Filter posts published on or before this date")
	linkedInLeadsListPostsCmd.Flags().IntVar(&llLimit, "limit", 0, "Max posts to return (1-100)")
	linkedInLeadsListPostsCmd.Flags().StringVar(&llCursor, "cursor", "", "Pagination cursor")
	linkedInLeadsListPostsCmd.MarkFlagRequired("profile-id")

	linkedInLeadsGetLeadsCmd.Flags().StringVar(&llProfileID, "profile-id", "", "LinkedIn profile ID (UUID)")
	linkedInLeadsGetLeadsCmd.Flags().StringVar(&llPostID, "post-id", "", "LinkedIn post ID (UUID)")
	linkedInLeadsGetLeadsCmd.Flags().StringVar(&llTypes, "types", "", "Comma-separated engagement types (LIKE,COMMENT,RESHARE)")
	linkedInLeadsGetLeadsCmd.Flags().IntVar(&llMinFollowerCount, "min-follower-count", 0, "Minimum follower count filter")
	linkedInLeadsGetLeadsCmd.Flags().IntVar(&llLimit, "limit", 0, "Max leads to return (1-250)")
	linkedInLeadsGetLeadsCmd.Flags().StringVar(&llCursor, "cursor", "", "Pagination cursor")
	linkedInLeadsGetLeadsCmd.MarkFlagRequired("profile-id")
	linkedInLeadsGetLeadsCmd.MarkFlagRequired("post-id")
}

var linkedInLeadsCmd = &cobra.Command{
	Use:   "linkedin-leads",
	Short: "LinkedIn leads scraping: posts and their engagers",
}

var linkedInLeadsListPostsCmd = &cobra.Command{
	Use:   "list-posts",
	Short: "List LinkedIn posts available for leads scraping on a profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewLinkedInLeadsService(c).ListPosts(llProfileID, api.LinkedInLeadsListPostsParams{
			StartDate: llStartDate,
			EndDate:   llEndDate,
			Limit:     llLimit,
			Cursor:    llCursor,
		})
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var linkedInLeadsGetLeadsCmd = &cobra.Command{
	Use:   "get-leads",
	Short: "Get leads (engagers) for a LinkedIn post",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewLinkedInLeadsService(c).GetLeadsByPost(llProfileID, llPostID, api.LinkedInLeadsGetLeadsParams{
			Types:            llTypes,
			MinFollowerCount: llMinFollowerCount,
			Limit:            llLimit,
			Cursor:           llCursor,
		})
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}
