package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// TestLabelCreate_RejectsEmptyName guards against Cobra's MarkFlagRequired
// being satisfied by presence alone. --name "" or whitespace would pass that
// check but reach the API and fail with a less-actionable error.
func TestLabelCreate_RejectsEmptyName(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	for _, name := range []string{"", "   "} {
		t.Run("name="+name, func(t *testing.T) {
			labelCreateName = ""
			labelCreateColor = ""

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{
				"label", "create",
				"--name", name,
				"--color", "green",
			})

			err := rootCmd.Execute()
			if err == nil {
				t.Fatalf("expected error for --name=%q", name)
			}
			if !strings.Contains(err.Error(), "name") {
				t.Errorf("expected error about name, got: %v", err)
			}
		})
	}
}
