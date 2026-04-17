package cmd

import (
	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/spf13/cobra"
)

var (
	linkedInURN      string
	linkedInUsername string
)

func init() {
	rootCmd.AddCommand(linkedInCmd)
	linkedInCmd.AddCommand(linkedInGetProfileCmd)
	linkedInCmd.AddCommand(linkedInGetMentionCmd)

	linkedInGetProfileCmd.Flags().StringVar(&linkedInURN, "urn", "", "LinkedIn URN (e.g. urn:li:person:abc123 or urn:li:organization:xyz789)")
	linkedInGetProfileCmd.MarkFlagRequired("urn")

	linkedInGetMentionCmd.Flags().StringVar(&linkedInUsername, "username", "", "LinkedIn vanity username")
	linkedInGetMentionCmd.MarkFlagRequired("username")
}

var linkedInCmd = &cobra.Command{
	Use:   "linkedin",
	Short: "LinkedIn profile lookups and mention formats",
}

var linkedInGetProfileCmd = &cobra.Command{
	Use:   "get-profile",
	Short: "Look up a LinkedIn profile by URN",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewLinkedInService(c).GetProfile(linkedInURN)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var linkedInGetMentionCmd = &cobra.Command{
	Use:   "get-mention",
	Short: "Get the mention format for a LinkedIn username",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewLinkedInService(c).GetMention(linkedInUsername)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}
