package cmd

import (
	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userListCmd)
}

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Workspace users",
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspace users",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		users, err := api.NewUserService(c).List()
		if err != nil {
			return err
		}
		return printResult(users)
	},
}
