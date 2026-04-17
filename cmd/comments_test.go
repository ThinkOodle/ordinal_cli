package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// TestCommentCreate_RejectsEmptyMessage guards against Cobra's MarkFlagRequired
// being satisfied by presence alone. --message "" or whitespace would pass that
// check but send an empty message the API rejects with a less-actionable error.
func TestCommentCreate_RejectsEmptyMessage(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	for _, msg := range []string{"", "   ", "\t\n"} {
		t.Run("message="+msg, func(t *testing.T) {
			commentPostID = ""
			commentMessage = ""

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{
				"comment", "create",
				"--post-id", "p-1",
				"--message", msg,
			})

			err := rootCmd.Execute()
			if err == nil {
				t.Fatalf("expected error for --message=%q", msg)
			}
			if !strings.Contains(err.Error(), "message") {
				t.Errorf("expected error about message, got: %v", err)
			}
		})
	}
}
