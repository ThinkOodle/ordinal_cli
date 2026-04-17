package cmd

import (
	"fmt"
	"strings"

	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/spf13/cobra"
)

var instagramSearchQuery string

func init() {
	rootCmd.AddCommand(instagramCmd)
	instagramCmd.AddCommand(instagramSearchLocationsCmd)

	instagramSearchLocationsCmd.Flags().StringVar(&instagramSearchQuery, "query", "", "Search query for location name")
	instagramSearchLocationsCmd.MarkFlagRequired("query")
}

var instagramCmd = &cobra.Command{
	Use:   "instagram",
	Short: "Instagram helpers",
}

var instagramSearchLocationsCmd = &cobra.Command{
	Use:   "search-locations",
	Short: "Search Instagram locations",
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(instagramSearchQuery) == "" {
			return fmt.Errorf("--query must not be empty")
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewInstagramService(c).SearchLocations(instagramSearchQuery)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}
