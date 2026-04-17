package cmd

import (
	"fmt"
	"strings"

	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/ordinal-cli/ordinal/internal/models"
	"github.com/spf13/cobra"
)

var (
	inviteCreateEmail string
	inviteID          string
)

func init() {
	rootCmd.AddCommand(inviteCmd)
	inviteCmd.AddCommand(inviteListCmd)
	inviteCmd.AddCommand(inviteCreateCmd)
	inviteCmd.AddCommand(inviteDeleteCmd)

	inviteCreateCmd.Flags().StringVar(&inviteCreateEmail, "email", "", "Email address to invite")
	inviteCreateCmd.MarkFlagRequired("email")

	inviteDeleteCmd.Flags().StringVar(&inviteID, "id", "", "Invite ID (UUID)")
	inviteDeleteCmd.MarkFlagRequired("id")
}

var inviteCmd = &cobra.Command{
	Use:   "invite",
	Short: "Manage workspace invites",
}

var inviteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List pending invites",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		invites, err := api.NewInviteService(c).List()
		if err != nil {
			return err
		}
		return printResult(invites)
	},
}

var inviteCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an invite",
	RunE: func(cmd *cobra.Command, args []string) error {
		// MarkFlagRequired only checks presence; --email "" or whitespace
		// would otherwise reach the API and fail with a less-actionable error.
		if strings.TrimSpace(inviteCreateEmail) == "" {
			return fmt.Errorf("--email must not be empty")
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		invite, err := api.NewInviteService(c).Create(models.CreateInviteRequest{Email: inviteCreateEmail})
		if err != nil {
			return err
		}
		return printResult(invite)
	},
}

var inviteDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an invite",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		if err := api.NewInviteService(c).Delete(inviteID); err != nil {
			return err
		}
		return printResult(deletedAck("invite", inviteID))
	},
}
