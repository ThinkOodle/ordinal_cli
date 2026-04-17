package cmd

import (
	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/ordinal-cli/ordinal/internal/models"
	"github.com/spf13/cobra"
)

var (
	labelCreateName  string
	labelCreateColor string
	labelID          string
)

func init() {
	rootCmd.AddCommand(labelCmd)
	labelCmd.AddCommand(labelListCmd)
	labelCmd.AddCommand(labelCreateCmd)
	labelCmd.AddCommand(labelDeleteCmd)

	labelCreateCmd.Flags().StringVar(&labelCreateName, "name", "", "Label name")
	labelCreateCmd.Flags().StringVar(&labelCreateColor, "color", "", "Label color: yellow, purple, orange, red, brown, green")
	labelCreateCmd.MarkFlagRequired("name")
	labelCreateCmd.MarkFlagRequired("color")

	labelDeleteCmd.Flags().StringVar(&labelID, "id", "", "Label ID (UUID)")
	labelDeleteCmd.MarkFlagRequired("id")
}

var labelCmd = &cobra.Command{
	Use:   "label",
	Short: "Manage labels",
	Long:  "Create, list, and delete labels used to organize posts and ideas.",
}

var labelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all labels",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		resp, err := api.NewLabelService(c).List()
		if err != nil {
			return err
		}
		return printResult(resp)
	},
}

var labelCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a label",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		label, err := api.NewLabelService(c).Create(models.CreateLabelRequest{
			Name:  labelCreateName,
			Color: labelCreateColor,
		})
		if err != nil {
			return err
		}
		return printResult(label)
	},
}

var labelDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a label",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := newClient()
		if err != nil {
			return err
		}
		if err := api.NewLabelService(c).Delete(labelID); err != nil {
			return err
		}
		return printResult(deletedAck("label", labelID))
	},
}
