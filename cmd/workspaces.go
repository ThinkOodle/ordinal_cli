package cmd

import (
	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workspaceCmd)
	workspaceCmd.AddCommand(workspaceGetCmd)
}

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage the current workspace",
}

var workspaceGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the current workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		ws, err := api.NewWorkspaceService(c).Get()
		if err != nil {
			return err
		}
		return printResult(ws)
	},
}
