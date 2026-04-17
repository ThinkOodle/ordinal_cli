package cmd

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/config"
	"github.com/ordinal-cli/ordinal/internal/output"
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

// TestRoot_RejectsEmptyRequiredStringFlags locks in that every required string
// flag across the CLI is rejected when passed as "" or whitespace. The shared
// PersistentPreRunE validator replaces per-command TrimSpace checks, so one
// case per representative command here guards the whole set from regressing.
func TestRoot_RejectsEmptyRequiredStringFlags(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")
	// Clear any --output value bled in from a previous test on this
	// shared rootCmd so PersistentPreRunE falls through to the default.
	cfgOutputFormat = ""

	tests := []struct {
		name     string
		args     []string
		wantFlag string
	}{
		{"comment list --post-id ''", []string{"comment", "list", "--post-id", ""}, "post-id"},
		{"comment list --post-id whitespace", []string{"comment", "list", "--post-id", "   "}, "post-id"},
		{"comment delete --id ''", []string{"comment", "delete", "--id", ""}, "id"},
		{"subscriber list --post-id ''", []string{"subscriber", "list", "--post-id", ""}, "post-id"},
		{"slack-boost list --post-id ''", []string{"slack-boost", "list", "--post-id", ""}, "post-id"},
		{"webhook get --id ''", []string{"webhook", "get", "--id", ""}, "id"},
		{"post get --id ''", []string{"post", "get", "--id", ""}, "id"},
		{"analytics linkedin-followers --profile-id ''", []string{"analytics", "linkedin-followers", "--profile-id", ""}, "profile-id"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs(tc.args)

			err := rootCmd.Execute()
			if err == nil {
				t.Fatalf("expected error for empty required flag")
			}
			if !strings.Contains(err.Error(), tc.wantFlag) {
				t.Errorf("expected error about %q, got: %v", tc.wantFlag, err)
			}
			if !strings.Contains(err.Error(), "must not be empty") {
				t.Errorf("expected 'must not be empty' message, got: %v", err)
			}
		})
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

// TestPrintRawJSON_EmptyBodyIsSuccess locks in that a 204-style zero-length
// response — or a 200 "{}" response — is treated as a structured success
// rather than a parse error or a bare map that extractRows cannot render.
// Any mutation endpoint that legitimately returns no body content (common
// for DELETE and some idempotent POSTs) routes through printRawJSON, so one
// case per format here guards the whole set from regressing to "unexpected
// end of JSON input" (empty body) or "unsupported data type" (empty object
// under --output table/csv).
func TestPrintRawJSON_EmptyBodyIsSuccess(t *testing.T) {
	// printRawJSON is called directly without going through Execute, so
	// appConfig isn't re-populated by PersistentPreRunE. Bypass state bled
	// in from earlier tests on this shared rootCmd.
	prev := appConfig
	t.Cleanup(func() { appConfig = prev })

	bodies := map[string][]byte{
		"nil":        nil,
		"empty":      {},
		"spaces":     []byte("   "),
		"whitespace": []byte("\n\t"),
		"empty_obj":  []byte("{}"),
	}
	formats := []string{"json", "table", "csv"}
	for name, data := range bodies {
		for _, format := range formats {
			t.Run(name+"/"+format, func(t *testing.T) {
				appConfig = &config.Config{OutputFormat: format}
				if err := printRawJSON(data); err != nil {
					t.Errorf("empty body should succeed, got: %v", err)
				}
			})
		}
	}
}

// TestPrintResult_EmptyCSVHasNoBlankLine guards the "stdout stays strictly
// parseable" contract for CSV output when the result set is empty. A lone
// trailing newline was being emitted because fmt.Println always ran, even
// when the formatter returned an empty body; a downstream csv.Reader would
// then see a single empty record and misreport row counts.
func TestPrintResult_EmptyCSVHasNoBlankLine(t *testing.T) {
	prev := appConfig
	t.Cleanup(func() { appConfig = prev })
	appConfig = &config.Config{OutputFormat: "csv"}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	origStdout := os.Stdout
	os.Stdout = w
	t.Cleanup(func() { os.Stdout = origStdout })

	if err := printResult([]map[string]interface{}{}); err != nil {
		t.Fatalf("printResult: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		t.Fatalf("read: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty stdout for empty csv; got %q", buf.String())
	}
}

// TestDeletedAck_FormatsParseable guards against regressing delete
// subcommands to raw fmt.Printf. All three formats must remain machine-
// readable: JSON must decode, CSV must have a header+row, and table must
// show the two scalar fields. Every delete command routes through
// deletedAck + printResult, so this one test covers the whole set.
func TestDeletedAck_FormatsParseable(t *testing.T) {
	ack := deletedAck("comment", "c-123")

	t.Run("json", func(t *testing.T) {
		out, _, err := output.FormatOutput(ack, output.FormatJSON)
		if err != nil {
			t.Fatalf("json: %v", err)
		}
		var decoded map[string]interface{}
		if err := json.Unmarshal([]byte(out), &decoded); err != nil {
			t.Fatalf("json output not parseable: %v\n%s", err, out)
		}
		if decoded["deleted"] != true || decoded["id"] != "c-123" || decoded["resource"] != "comment" {
			t.Errorf("unexpected decoded json: %+v", decoded)
		}
	})

	t.Run("csv", func(t *testing.T) {
		out, _, err := output.FormatOutput(ack, output.FormatCSV)
		if err != nil {
			t.Fatalf("csv: %v", err)
		}
		records, err := csv.NewReader(strings.NewReader(out)).ReadAll()
		if err != nil {
			t.Fatalf("csv output not parseable: %v\n%s", err, out)
		}
		if len(records) != 2 {
			t.Fatalf("expected header + 1 row, got %d rows: %v", len(records), records)
		}
		// Headers are alphabetically sorted by mapSliceToRows.
		wantHeader := []string{"deleted", "id", "resource"}
		for i, h := range wantHeader {
			if records[0][i] != h {
				t.Errorf("header[%d] = %q, want %q", i, records[0][i], h)
			}
		}
		if records[1][0] != "true" || records[1][1] != "c-123" || records[1][2] != "comment" {
			t.Errorf("unexpected csv row: %v", records[1])
		}
	})

	t.Run("table", func(t *testing.T) {
		out, _, err := output.FormatOutput(ack, output.FormatTable)
		if err != nil {
			t.Fatalf("table: %v", err)
		}
		for _, want := range []string{"DELETED", "ID", "RESOURCE", "c-123", "comment", "true"} {
			if !strings.Contains(out, want) {
				t.Errorf("table output missing %q:\n%s", want, out)
			}
		}
	})
}
