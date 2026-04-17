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
		return printMutationAck(data)
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

// printRawJSON parses a raw JSON payload and formats it via printResult. Use
// this for GET/read endpoints: a zero-length body — what the API returns for
// 204 No Content — is treated as a successful acknowledgement rather than a
// parse error, but a legitimate empty JSON object ({}) passes through
// unchanged so the caller sees whatever the API actually returned. Mutation
// endpoints that conventionally respond with {} on success should call
// printMutationAck instead.
func printRawJSON(data []byte) error {
	if len(bytes.TrimSpace(data)) == 0 {
		return printResult(map[string]interface{}{"success": true})
	}
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("parsing response: %w", err)
	}
	// A legitimate {} from a read endpoint must render as "{}" across
	// every output format. The table/CSV formatter deliberately collapses
	// empty objects to "No results" / empty body (mutation acks need
	// that), so intercept here to keep the read-path promise that {}
	// passes through unchanged — regardless of --output.
	if m, ok := v.(map[string]interface{}); ok && len(m) == 0 {
		fmt.Println("{}")
		return nil
	}
	return printResult(v)
}

// printMutationAck parses a raw JSON payload from a create/update/delete
// endpoint and formats it via printResult. Like printRawJSON it treats a
// zero-length body as success, but it also collapses an empty JSON object
// ({}) to the same acknowledgement — the Ordinal API returns {} from several
// mutation paths on success, and without this the table/csv formatters
// would render "No results" where the user expects a clear confirmation.
// Read endpoints must NOT use this helper: they may legitimately return {}
// as the real response, and we want to preserve that fidelity.
func printMutationAck(data []byte) error {
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
