package cmd

import (
	"fmt"
	"strings"

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
		// cobra.ExactArgs(1) only checks argument count, so an empty or
		// whitespace-only value (e.g. `ordinal auth ""` or `ordinal auth "   "`)
		// would otherwise write a blank key to the config file and then be
		// sent in the Authorization header on every subsequent request. Trim
		// and require a non-empty value so the failure surfaces here with an
		// actionable message instead of as a confusing 401 later.
		apiKey := strings.TrimSpace(args[0])
		if apiKey == "" {
			return fmt.Errorf("api key must not be empty")
		}

		if err := config.SaveAPIKey(apiKey); err != nil {
			return fmt.Errorf("saving api key: %w", err)
		}

		configFile, _ := config.ConfigFilePath()
		fmt.Printf("API key saved to %s\n", configFile)
		return nil
	},
}
