package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// TestLinkedInLeadsListPosts_RejectsLimitOutOfRange locks in that the help
// text's advertised 1-100 range is enforced locally. The previous behavior
// silently dropped invalid values from the query and used server defaults,
// masking the bad flag value entirely.
func TestLinkedInLeadsListPosts_RejectsLimitOutOfRange(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	for _, limit := range []string{"101", "500", "-1", "0"} {
		t.Run("limit="+limit, func(t *testing.T) {
			llLimit = 0

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{"linkedin-leads", "list-posts", "--profile-id", "p1", "--limit", limit})

			err := rootCmd.Execute()
			if err == nil {
				t.Fatalf("expected error for --limit=%s", limit)
			}
			if !strings.Contains(err.Error(), "limit") {
				t.Errorf("expected limit error, got: %v", err)
			}
		})
	}
}

// TestLinkedInLeadsGetLeads_RejectsLimitOutOfRange locks in that the help
// text's advertised 1-250 range is enforced locally.
func TestLinkedInLeadsGetLeads_RejectsLimitOutOfRange(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	for _, limit := range []string{"251", "1000", "-1", "0"} {
		t.Run("limit="+limit, func(t *testing.T) {
			llLimit = 0

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{"linkedin-leads", "get-leads", "--profile-id", "p1", "--post-id", "po1", "--limit", limit})

			err := rootCmd.Execute()
			if err == nil {
				t.Fatalf("expected error for --limit=%s", limit)
			}
			if !strings.Contains(err.Error(), "limit") {
				t.Errorf("expected limit error, got: %v", err)
			}
		})
	}
}
