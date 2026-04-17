package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// TestPostCreate_RequiresTitlePublishAtStatus locks in the client-side
// validation of the fields the Ordinal API requires. Without this, an empty
// or partial --body-json quietly hit the API and failed server-side with a
// less-actionable error message.
func TestPostCreate_RequiresTitlePublishAtStatus(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	tests := []struct {
		name      string
		bodyJSON  string
		wantField string
	}{
		{"empty body", `{}`, "title"},
		{"only title", `{"title":"hi"}`, "publishAt"},
		{"title and publishAt", `{"title":"hi","publishAt":"2026-01-01T00:00:00Z"}`, "status"},
		// Typed-check regression cases: the pre-fix guard treated null,
		// whitespace-only, and non-string values as valid, sending them to
		// the API instead of failing locally.
		{"null title", `{"title":null,"publishAt":"2026-01-01T00:00:00Z","status":"ToDo"}`, "title"},
		{"whitespace title", `{"title":"   ","publishAt":"2026-01-01T00:00:00Z","status":"ToDo"}`, "title"},
		{"numeric status", `{"title":"hi","publishAt":"2026-01-01T00:00:00Z","status":1}`, "status"},
		{"null publishAt", `{"title":"hi","publishAt":null,"status":"ToDo"}`, "publishAt"},
		{"all null", `{"title":null,"publishAt":null,"status":null}`, "title"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetPostCreateFlags(t)

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{"post", "create", "--body-json", tc.bodyJSON})

			err := rootCmd.Execute()
			if err == nil {
				t.Fatalf("expected validation error")
			}
			if !strings.Contains(err.Error(), tc.wantField) {
				t.Errorf("expected error about %q field, got: %v", tc.wantField, err)
			}
		})
	}
}

// TestPostList_RejectsLimitOutOfRange locks in that the --limit flag's
// advertised 1-100 range is enforced locally; otherwise help text and
// runtime behavior drift and callers burn a round-trip on an obviously
// invalid value.
func TestPostList_RejectsLimitOutOfRange(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	for _, limit := range []string{"101", "1000", "-1", "0"} {
		t.Run("limit="+limit, func(t *testing.T) {
			postListLimit = 0

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{"post", "list", "--limit", limit})

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

// resetPostCreateFlags resets the package-level flag variables that
// postCreateCmd reads so successive tests don't leak state. Cobra's flag
// "changed" bits are reset automatically by re-binding via SetArgs, but the
// destination variables are global and must be cleared explicitly.
func resetPostCreateFlags(t *testing.T) {
	t.Helper()
	postCreateTitle = ""
	postCreatePublishAt = ""
	postCreateStatus = ""
	postCreateLabelIDs = ""
	postCreateCampaignID = ""
	postCreateNotes = ""
	postCreateBodyJSON = ""
	postCreateBodyFile = ""
}
