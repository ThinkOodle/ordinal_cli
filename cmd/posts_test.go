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
