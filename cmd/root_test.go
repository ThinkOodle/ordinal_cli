package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/config"
	"github.com/spf13/cobra"
)

// TestRoot_InvalidOutputFormatRejected exercises PersistentPreRunE and asserts
// that an unknown --output value fails fast instead of silently falling back
// to JSON. The pre-run hook runs before any command's RunE, so the attached
// dummy subcommand is never reached when validation trips.
func TestRoot_InvalidOutputFormatRejected(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	dummy := &cobra.Command{
		Use:  "dummy-probe",
		RunE: func(cmd *cobra.Command, args []string) error { return nil },
	}
	rootCmd.AddCommand(dummy)
	t.Cleanup(func() { rootCmd.RemoveCommand(dummy) })

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"dummy-probe", "--output", "yaml"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatalf("expected error for invalid --output")
	}
	if !strings.Contains(err.Error(), "invalid output format") {
		t.Errorf("expected invalid output format error, got %v", err)
	}
}

// TestRoot_AuthBypassesOutputFormatValidation guards against regressing to a
// state where a bad saved output_format blocks the very command meant to
// repair auth state. A user with a typo in config.yaml must still be able to
// run `ordinal auth <key>` to recover.
func TestRoot_AuthBypassesOutputFormatValidation(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("ORDINAL_API_KEY", "")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	dir := filepath.Join(home, ".config", "ordinal")
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte("output_format: yaml\n"), 0600); err != nil {
		t.Fatalf("write: %v", err)
	}

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"auth", "new-key"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("auth should succeed despite invalid saved output_format, got: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("post-auth Load: %v", err)
	}
	if cfg.APIKey != "new-key" {
		t.Errorf("expected saved api key new-key, got %q", cfg.APIKey)
	}
}

// TestBodyJSONFlagPrecedence_UpdatePattern locks in the precedence rule that
// the help text for --body-json on the *update commands* promises: the body
// is parsed first, then individual flags override matching top-level keys
// (using cmd.Flags().Changed so explicit empty overrides count).
func TestBodyJSONFlagPrecedence_UpdatePattern(t *testing.T) {
	var bodyJSONArg, titleArg, statusArg string
	var merged map[string]interface{}

	cmd := &cobra.Command{
		Use: "toy",
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseBodyJSON(bodyJSONArg, "")
			if err != nil {
				return err
			}
			if body == nil {
				body = map[string]interface{}{}
			}
			if cmd.Flags().Changed("title") {
				body["title"] = titleArg
			}
			if cmd.Flags().Changed("status") {
				body["status"] = statusArg
			}
			merged = body
			return nil
		},
	}
	cmd.Flags().StringVar(&bodyJSONArg, "body-json", "", "")
	cmd.Flags().StringVar(&titleArg, "title", "", "")
	cmd.Flags().StringVar(&statusArg, "status", "", "")

	cmd.SetArgs([]string{
		"--body-json", `{"title":"from-body","status":"ToDo","extra":"keep"}`,
		"--title", "from-flag",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	if merged["title"] != "from-flag" {
		t.Errorf("flag should override matching body key: got %v", merged["title"])
	}
	if merged["status"] != "ToDo" {
		t.Errorf("unset flag must leave body value intact: got %v", merged["status"])
	}
	if merged["extra"] != "keep" {
		t.Errorf("non-matching body key must survive: got %v", merged["extra"])
	}
}

// TestBodyJSONFlagPrecedence_CreatePattern covers the create-side pattern used
// by posts/ideas/slack-boost create commands: flags override body keys when
// the flag has a non-empty value (the !="" check rather than Changed()).
func TestBodyJSONFlagPrecedence_CreatePattern(t *testing.T) {
	var bodyJSONArg, titleArg string
	var merged map[string]interface{}

	cmd := &cobra.Command{
		Use: "toy",
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := parseBodyJSON(bodyJSONArg, "")
			if err != nil {
				return err
			}
			if body == nil {
				body = map[string]interface{}{}
			}
			if titleArg != "" {
				body["title"] = titleArg
			}
			merged = body
			return nil
		},
	}
	cmd.Flags().StringVar(&bodyJSONArg, "body-json", "", "")
	cmd.Flags().StringVar(&titleArg, "title", "", "")

	cmd.SetArgs([]string{
		"--body-json", `{"title":"from-body","extra":"keep"}`,
		"--title", "from-flag",
	})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute: %v", err)
	}

	if merged["title"] != "from-flag" {
		t.Errorf("flag should override matching body key: got %v", merged["title"])
	}
	if merged["extra"] != "keep" {
		t.Errorf("non-matching body key must survive: got %v", merged["extra"])
	}
}
