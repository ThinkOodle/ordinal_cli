package cmd

import (
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(authCmd)
}

var authCmd = &cobra.Command{
	Use:         "auth <api-key>",
	Short:       "Set the API key for authentication",
	Long:        "Save your Ordinal API key to ~/.config/ordinal/config.yaml for use in subsequent commands.",
	Args:        cobra.ExactArgs(1),
	Annotations: map[string]string{skipOutputValidationAnnotation: "true"},
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey := args[0]

		if err := config.SaveAPIKey(apiKey); err != nil {
			return fmt.Errorf("saving api key: %w", err)
		}

		configFile, _ := config.ConfigFilePath()
		fmt.Printf("API key saved to %s\n", configFile)
		return nil
	},
}
