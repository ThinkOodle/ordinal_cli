// Package cmd defines all CLI commands for the Ordinal CLI.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/config"
	"github.com/ordinal-cli/ordinal/internal/output"
	"github.com/spf13/cobra"
)

var (
	cfgAPIKey       string
	cfgOutputFormat string
	cfgVerbose      bool
	appConfig       *config.Config
)

// skipOutputValidationAnnotation marks commands (or command groups) that do
// not render API output and therefore should not fail when the saved
// output_format is invalid. Without this, a bad config value would brick the
// very commands meant to repair config state (`auth`) or perform local
// utility work (`skill`, shell completion).
const skipOutputValidationAnnotation = "ordinal/skipOutputValidation"

// skipsOutputValidation reports whether cmd or any of its ancestors carry the
// skip-output-validation annotation. Walking the parent chain lets the
// annotation apply to an entire subtree (e.g. `skill` covers `skill install`).
func skipsOutputValidation(cmd *cobra.Command) bool {
	for c := cmd; c != nil; c = c.Parent() {
		if c.Annotations[skipOutputValidationAnnotation] == "true" {
			return true
		}
		// Cobra auto-generates the `completion` command; we can't set
		// annotations on it from here, so match by name.
		if c.Name() == "completion" {
			return true
		}
	}
	return false
}

// rootCmd is the base command for the CLI.
var rootCmd = &cobra.Command{
	Use:   "ordinal",
	Short: "CLI for the Ordinal (tryordinal.com) API",
	Long:  "A command-line interface for managing Ordinal social posts, ideas, approvals, comments, analytics, webhooks, and related workspace resources.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			// A broken config file must not block higher-priority sources
			// (flags, env vars) or the `auth` command that can repair it.
			// Load still returns env-derived values alongside the error, so
			// we warn but keep the partial config rather than dropping it.
			fmt.Fprintf(os.Stderr, "warning: existing config could not be loaded (%v); continuing with env and defaults\n", err)
		}
		if cfg == nil {
			cfg = &config.Config{}
		}
		appConfig = cfg

		if cmd.Flags().Changed("output") {
			appConfig.OutputFormat = cfgOutputFormat
		}
		if cmd.Flags().Changed("verbose") {
			appConfig.Verbose = cfgVerbose
		}

		if appConfig.OutputFormat == "" {
			appConfig.OutputFormat = config.DefaultOutputFormat
		}
		if skipsOutputValidation(cmd) {
			return nil
		}
		if !output.IsValidFormat(output.Format(appConfig.OutputFormat)) {
			return fmt.Errorf("invalid output format %q: must be one of json, table, csv", appConfig.OutputFormat)
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgAPIKey, "api-key", "k", "", "API key (env: ORDINAL_API_KEY)")
	rootCmd.PersistentFlags().StringVarP(&cfgOutputFormat, "output", "o", "", "Output format: json, table, csv (default: json)")
	rootCmd.PersistentFlags().BoolVarP(&cfgVerbose, "verbose", "v", false, "Verbose output (shows request/response details)")
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// newClient creates an authenticated API client from the current configuration.
func newClient() (*client.Client, error) {
	apiKey, err := config.GetAPIKey(cfgAPIKey, appConfig)
	if err != nil {
		return nil, err
	}

	opts := []client.Option{}
	if appConfig != nil && appConfig.Verbose {
		opts = append(opts, client.WithVerbose(true))
	}

	return client.New(apiKey, opts...), nil
}

// getOutputFormat returns the current output format.
func getOutputFormat() output.Format {
	if appConfig != nil && appConfig.OutputFormat != "" {
		return output.Format(appConfig.OutputFormat)
	}
	return output.FormatJSON
}

// printResult formats and prints the result according to the current output format.
// For CSV, the pagination footer goes to stderr so stdout stays strictly
// parseable. For table output the footer goes to stdout beneath the rows so it
// isn't separated from the table when stdout is piped through a pager or
// redirected (stdout/stderr are buffered independently, so stderr footers can
// arrive out of order or bypass the pager entirely).
func printResult(data interface{}) error {
	format := getOutputFormat()
	out, footer, err := output.FormatOutput(data, format)
	if err != nil {
		return fmt.Errorf("formatting output: %w", err)
	}
	fmt.Println(out)
	if footer != "" {
		if format == output.FormatCSV {
			fmt.Fprintln(os.Stderr, footer)
		} else {
			fmt.Println(footer)
		}
	}
	return nil
}

// deletedAck is the structured acknowledgement returned by delete subcommands
// so their output flows through printResult rather than raw fmt.Printf. This
// keeps --output json and --output csv machine-readable — a plain "X deleted"
// line would corrupt both.
func deletedAck(resource, id string) map[string]interface{} {
	return map[string]interface{}{
		"resource": resource,
		"id":       id,
		"deleted":  true,
	}
}

// splitCSV splits a comma-separated flag value and trims each entry.
// Returns nil if the input is empty.
func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// parseBodyJSON reads a JSON body from either an inline string or a file path.
// If both are empty, returns nil. If filePath is "-", reads from stdin.
func parseBodyJSON(bodyJSON, bodyFile string) (map[string]interface{}, error) {
	raw, err := readBodyJSONRaw(bodyJSON, bodyFile)
	if err != nil || raw == nil {
		return nil, err
	}

	var body map[string]interface{}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, fmt.Errorf("parsing body json: %w", err)
	}
	return body, nil
}

// readBodyJSONRaw returns the raw JSON bytes from either an inline string or a
// file path without unmarshaling. Useful when the caller needs to accept either
// an object or an array at the top level.
func readBodyJSONRaw(bodyJSON, bodyFile string) ([]byte, error) {
	switch {
	case bodyJSON != "" && bodyFile != "":
		return nil, fmt.Errorf("only one of --body-json and --body-file may be set")
	case bodyJSON != "":
		return []byte(bodyJSON), nil
	case bodyFile == "-":
		raw, err := readAllStdin()
		if err != nil {
			return nil, fmt.Errorf("reading stdin: %w", err)
		}
		return raw, nil
	case bodyFile != "":
		raw, err := os.ReadFile(bodyFile)
		if err != nil {
			return nil, fmt.Errorf("reading body file: %w", err)
		}
		return raw, nil
	}
	return nil, nil
}

// readBodyFile reads a JSON body file, or stdin if path is "-".
func readBodyFile(path string) ([]byte, error) {
	if path == "-" {
		raw, err := readAllStdin()
		if err != nil {
			return nil, fmt.Errorf("reading stdin: %w", err)
		}
		return raw, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading body file: %w", err)
	}
	return raw, nil
}

func readAllStdin() ([]byte, error) {
	var buf []byte
	tmp := make([]byte, 4096)
	for {
		n, err := os.Stdin.Read(tmp)
		if n > 0 {
			buf = append(buf, tmp[:n]...)
		}
		if err != nil {
			if err.Error() == "EOF" {
				return buf, nil
			}
			return buf, err
		}
	}
}
