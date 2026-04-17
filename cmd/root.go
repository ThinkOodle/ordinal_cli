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
	cfgNoColor      bool
	cfgVerbose      bool
	appConfig       *config.Config
)

// rootCmd is the base command for the CLI.
var rootCmd = &cobra.Command{
	Use:   "ordinal",
	Short: "CLI for the Ordinal (tryordinal.com) API",
	Long:  "A command-line interface for managing Ordinal social posts, ideas, approvals, comments, analytics, webhooks, and related workspace resources.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		appConfig = cfg

		if cfgOutputFormat != "" {
			appConfig.OutputFormat = cfgOutputFormat
		}
		if cfgNoColor {
			appConfig.NoColor = true
		}
		if cfgVerbose {
			appConfig.Verbose = true
		}

		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgAPIKey, "api-key", "k", "", "API key (env: ORDINAL_API_KEY)")
	rootCmd.PersistentFlags().StringVarP(&cfgOutputFormat, "output", "o", "", "Output format: json, table, csv (default: json)")
	rootCmd.PersistentFlags().BoolVar(&cfgNoColor, "no-color", false, "Disable colored output")
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
func printResult(data interface{}) error {
	format := getOutputFormat()
	out, err := output.FormatOutput(data, format)
	if err != nil {
		return fmt.Errorf("formatting output: %w", err)
	}
	fmt.Println(out)
	return nil
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
	var raw []byte
	var err error

	switch {
	case bodyJSON != "" && bodyFile != "":
		return nil, fmt.Errorf("only one of --body-json and --body-file may be set")
	case bodyJSON != "":
		raw = []byte(bodyJSON)
	case bodyFile == "-":
		raw, err = readAllStdin()
		if err != nil {
			return nil, fmt.Errorf("reading stdin: %w", err)
		}
	case bodyFile != "":
		raw, err = os.ReadFile(bodyFile)
		if err != nil {
			return nil, fmt.Errorf("reading body file: %w", err)
		}
	default:
		return nil, nil
	}

	var body map[string]interface{}
	if err := json.Unmarshal(raw, &body); err != nil {
		return nil, fmt.Errorf("parsing body json: %w", err)
	}
	return body, nil
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
