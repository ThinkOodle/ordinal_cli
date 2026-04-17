package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"

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
// A zero-length (or whitespace-only) body — what the API returns for 204
// No Content — is treated as a successful acknowledgement rather than a
// parse error, so delete/update endpoints that intentionally return no body
// still exit 0 with a machine-readable confirmation. An empty JSON object
// ({}) gets the same treatment: extractRows has no columns to render from
// a zero-field map, so --output table/csv would otherwise fail on a
// successful mutation that simply returns no body fields.
func printRawJSON(data []byte) error {
	if len(bytes.TrimSpace(data)) == 0 {
		return printResult(map[string]interface{}{"success": true})
	}
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}
	if m, ok := v.(map[string]interface{}); ok && len(m) == 0 {
		return printResult(map[string]interface{}{"success": true})
	}
	return printResult(v)
}
