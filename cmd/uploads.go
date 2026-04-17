package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ordinal-cli/ordinal/internal/api"
	"github.com/ordinal-cli/ordinal/internal/models"
	"github.com/spf13/cobra"
)

var (
	uploadCreateURL string
	uploadID        string
)

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.AddCommand(uploadCreateCmd)
	uploadCmd.AddCommand(uploadGetCmd)

	uploadCreateCmd.Flags().StringVar(&uploadCreateURL, "url", "", "Publicly accessible URL of the file to upload")
	uploadCreateCmd.MarkFlagRequired("url")

	uploadGetCmd.Flags().StringVar(&uploadID, "id", "", "Upload job ID (UUID)")
	uploadGetCmd.MarkFlagRequired("id")
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload files from URLs and check upload status",
}

var uploadCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Upload a file from a URL",
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(uploadCreateURL) == "" {
			return fmt.Errorf("--url must not be empty")
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewUploadService(c).Create(models.CreateUploadRequest{URL: uploadCreateURL})
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

var uploadGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the status of an upload job",
	RunE: func(cmd *cobra.Command, args []string) error {
		if strings.TrimSpace(uploadID) == "" {
			return fmt.Errorf("--id must not be empty")
		}
		c, err := newClient()
		if err != nil {
			return err
		}
		data, err := api.NewUploadService(c).Get(uploadID)
		if err != nil {
			return err
		}
		return printRawJSON(data)
	},
}

// printRawJSON parses a raw JSON payload and formats it via printResult.
func printRawJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}
	return printResult(v)
}
