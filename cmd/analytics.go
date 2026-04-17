package cmd

import (
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/ordinal-cli/ordinal/internal/models"
	"github.com/spf13/cobra"
)

var (
	analyticsProfileID         string
	analyticsStartDate         string
	analyticsEndDate           string
	analyticsCpmUpdateBodyJSON string
	analyticsCpmUpdateBodyFile string
	analyticsCpmLinkedIn       float64
	analyticsCpmX              float64
	analyticsCpmInstagram      float64
	analyticsCpmFacebook       float64
	analyticsCpmThreads        float64
)

func init() {
	rootCmd.AddCommand(analyticsCmd)
	analyticsCmd.AddCommand(analyticsCpmGetCmd)
	analyticsCmd.AddCommand(analyticsCpmUpdateCmd)
	analyticsCmd.AddCommand(analyticsLinkedInFollowersCmd)
	analyticsCmd.AddCommand(analyticsLinkedInPostsCmd)
	analyticsCmd.AddCommand(analyticsXFollowersCmd)
	analyticsCmd.AddCommand(analyticsXPostsCmd)

	analyticsCpmUpdateCmd.Flags().Float64Var(&analyticsCpmLinkedIn, "linkedin", 0, "LinkedIn CPM value")
	analyticsCpmUpdateCmd.Flags().Float64Var(&analyticsCpmX, "x", 0, "X (Twitter) CPM value")
	analyticsCpmUpdateCmd.Flags().Float64Var(&analyticsCpmInstagram, "instagram", 0, "Instagram CPM value")
	analyticsCpmUpdateCmd.Flags().Float64Var(&analyticsCpmFacebook, "facebook", 0, "Facebook CPM value")
	analyticsCpmUpdateCmd.Flags().Float64Var(&analyticsCpmThreads, "threads", 0, "Threads CPM value")
	analyticsCpmUpdateCmd.Flags().StringVar(&analyticsCpmUpdateBodyJSON, "body-json", "", "Full JSON body (overrides individual flags)")
	analyticsCpmUpdateCmd.Flags().StringVar(&analyticsCpmUpdateBodyFile, "body-file", "", "Path to JSON body file (or - for stdin)")

	for _, c := range []*cobra.Command{analyticsLinkedInFollowersCmd, analyticsLinkedInPostsCmd, analyticsXFollowersCmd, analyticsXPostsCmd} {
		c.Flags().StringVar(&analyticsProfileID, "profile-id", "", "Profile ID (UUID)")
		c.Flags().StringVar(&analyticsStartDate, "start", "", "Start date (defaults to 30 days ago)")
		c.Flags().StringVar(&analyticsEndDate, "end", "", "End date (defaults to today)")
		c.MarkFlagRequired("profile-id")
	}
}

var analyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "Analytics: follower growth, post performance, and CPM configuration",
}

var analyticsCpmGetCmd = &cobra.Command{
	Use:   "cpm-get",
	Short: "Get CPM values for EMV calculations",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		v, err := api.NewAnalyticsService(c).GetCpm()
		if err != nil {
			return err
		}
		return printResult(v)
	},
}

var analyticsCpmUpdateCmd = &cobra.Command{
	Use:   "cpm-update",
	Short: "Update CPM values",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}

		body, err := parseBodyJSON(analyticsCpmUpdateBodyJSON, analyticsCpmUpdateBodyFile)
		if err != nil {
			return err
		}

		req := models.CpmUpdateRequest{}
		if body != nil {
			if v, ok := body["linkedIn"].(float64); ok {
				req.LinkedIn = &v
			}
			if v, ok := body["x"].(float64); ok {
				req.X = &v
			}
			if v, ok := body["instagram"].(float64); ok {
				req.Instagram = &v
			}
			if v, ok := body["facebook"].(float64); ok {
				req.Facebook = &v
			}
			if v, ok := body["threads"].(float64); ok {
				req.Threads = &v
			}
		} else {
			if cmd.Flags().Changed("linkedin") {
				req.LinkedIn = &analyticsCpmLinkedIn
			}
			if cmd.Flags().Changed("x") {
				req.X = &analyticsCpmX
			}
			if cmd.Flags().Changed("instagram") {
				req.Instagram = &analyticsCpmInstagram
			}
			if cmd.Flags().Changed("facebook") {
				req.Facebook = &analyticsCpmFacebook
			}
			if cmd.Flags().Changed("threads") {
				req.Threads = &analyticsCpmThreads
			}
		}

		if req.LinkedIn == nil && req.X == nil && req.Instagram == nil && req.Facebook == nil && req.Threads == nil {
			return fmt.Errorf("provide at least one platform CPM flag or a --body-json/--body-file")
		}

		v, err := api.NewAnalyticsService(c).UpdateCpm(req)
		if err != nil {
			return err
		}
		return printResult(v)
	},
}

func analyticsDateRange() models.AnalyticsDateRange {
	return models.AnalyticsDateRange{StartDate: analyticsStartDate, EndDate: analyticsEndDate}
}

var analyticsLinkedInFollowersCmd = &cobra.Command{
	Use:   "linkedin-followers",
	Short: "LinkedIn follower growth for a profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		points, err := api.NewAnalyticsService(c).LinkedInFollowers(analyticsProfileID, analyticsDateRange())
		if err != nil {
			return err
		}
		return printResult(points)
	},
}

var analyticsLinkedInPostsCmd = &cobra.Command{
	Use:   "linkedin-posts",
	Short: "LinkedIn post analytics for a profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewAnalyticsService(c).LinkedInPosts(analyticsProfileID, analyticsDateRange())
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var analyticsXFollowersCmd = &cobra.Command{
	Use:   "x-followers",
	Short: "X (Twitter) follower growth for a profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		points, err := api.NewAnalyticsService(c).XFollowers(analyticsProfileID, analyticsDateRange())
		if err != nil {
			return err
		}
		return printResult(points)
	},
}

var analyticsXPostsCmd = &cobra.Command{
	Use:   "x-posts",
	Short: "X (Twitter) post analytics for a profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewAnalyticsService(c).XPosts(analyticsProfileID, analyticsDateRange())
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}
