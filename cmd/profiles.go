package cmd

import (
	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(profileCmd)
	profileCmd.AddCommand(profileListSchedulingCmd)
	profileCmd.AddCommand(profileListEngagementCmd)
}

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage connected social profiles",
}

var profileListSchedulingCmd = &cobra.Command{
	Use:   "list-scheduling",
	Short: "List scheduling profiles (full capabilities)",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		profiles, err := api.NewProfileService(c).ListScheduling()
		if err != nil {
			return err
		}
		return printResult(profiles)
	},
}

var profileListEngagementCmd = &cobra.Command{
	Use:   "list-engagement",
	Short: "List engagement profiles (auto-engagement only)",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		profiles, err := api.NewProfileService(c).ListEngagement()
		if err != nil {
			return err
		}
		return printResult(profiles)
	},
}
