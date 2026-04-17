package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// TestInviteCreate_RejectsEmptyEmail guards against Cobra's MarkFlagRequired
// being satisfied by presence alone. --email "" or whitespace would pass that
// check but reach the API and fail with a less-actionable error.
func TestInviteCreate_RejectsEmptyEmail(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	for _, email := range []string{"", "   ", "\t"} {
		t.Run("email="+email, func(t *testing.T) {
			inviteCreateEmail = ""

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{"invite", "create", "--email", email})

			err := rootCmd.Execute()
			if err == nil {
				t.Fatalf("expected error for --email=%q", email)
			}
			if !strings.Contains(err.Error(), "email") {
				t.Errorf("expected error about email, got: %v", err)
			}
		})
	}
}
