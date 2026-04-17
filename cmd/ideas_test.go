package cmd

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/spf13/pflag"
)

// TestIdeaCreate_RequiresNonEmptyTitle locks in the client-side title
// validation. Previously the command only checked that a "title" key existed,
// so {"title":""} and {"title":null} reached the API and failed there with a
// less-actionable error. Validation also runs before newClient() so a missing
// title surfaces ahead of auth errors.
func TestIdeaCreate_RequiresNonEmptyTitle(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	tests := []struct {
		name     string
		bodyJSON string
	}{
		{"empty body", `{}`},
		{"empty title", `{"title":""}`},
		{"whitespace title", `{"title":"   "}`},
		{"null title", `{"title":null}`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetIdeaCreateFlags(t)

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{"idea", "create", "--body-json", tc.bodyJSON})

			err := rootCmd.Execute()
			if err == nil {
				t.Fatalf("expected validation error for body %s", tc.bodyJSON)
			}
			if !strings.Contains(err.Error(), "title") {
				t.Errorf("expected error about title, got: %v", err)
			}
		})
	}
}

// TestIdeaCreate_ValidatesBeforeAuth guards the ordering: validation must run
// before newClient() so a user with no configured API key still gets the
// more-actionable "title required" message for a bad body, matching the post
// create behavior.
func TestIdeaCreate_ValidatesBeforeAuth(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	resetIdeaCreateFlags(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"idea", "create", "--body-json", `{}`})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "api key") {
		t.Errorf("validation should run before auth; got auth error: %v", err)
	}
	if !strings.Contains(err.Error(), "title") {
		t.Errorf("expected title error, got: %v", err)
	}
}

// TestIdeaList_RejectsLimitOutOfRange parity test for post list: help text
// advertises 1-100 but the API call previously accepted any positive value.
// Enforcing locally keeps help, behavior, and the API agreeing.
func TestIdeaList_RejectsLimitOutOfRange(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	for _, limit := range []string{"101", "1000", "-1", "0"} {
		t.Run("limit="+limit, func(t *testing.T) {
			ideaListLimit = 0

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{"idea", "list", "--limit", limit})

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

func resetIdeaCreateFlags(t *testing.T) {
	t.Helper()
	ideaCreateTitle = ""
	ideaCreateLabelIDs = ""
	ideaCreateCampaignID = ""
	ideaCreateBodyJSON = ""
	ideaCreateBodyFile = ""
}

// TestIdeaUpdate_ClearsLabelsWithEmptyFlag is the idea-update parity of
// TestPostUpdate_ClearsLabelsWithEmptyFlag. An explicit --label-ids "" must
// serialize as an empty JSON array so the API sees "clear all labels"
// rather than a null that some validators reject outright.
func TestIdeaUpdate_ClearsLabelsWithEmptyFlag(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	var captured []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer server.Close()

	prev := testClientOpts
	testClientOpts = []client.Option{
		client.WithBaseURL(server.URL),
		client.WithHTTPClient(server.Client()),
	}
	defer func() { testClientOpts = prev }()

	tests := []struct {
		name     string
		labelArg string
	}{
		{"empty string", ""},
		{"just commas", ","},
		{"whitespace entries", "  ,  "},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			captured = nil
			resetIdeaUpdateFlags(t)

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{"idea", "update", "--id", "i-1", "--label-ids", tc.labelArg})

			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("idea update: %v", err)
			}
			if !strings.Contains(string(captured), `"labelIds":[]`) {
				t.Errorf("expected labelIds:[] in body, got: %s", captured)
			}
			if strings.Contains(string(captured), `"labelIds":null`) {
				t.Errorf("labelIds must not marshal as null: %s", captured)
			}
		})
	}
}

// resetIdeaUpdateFlags clears ideaUpdateCmd's package-level flag variables
// and the pflag "changed" bits so consecutive Execute calls on the shared
// rootCmd don't leak state between tests.
func resetIdeaUpdateFlags(t *testing.T) {
	t.Helper()
	ideaID = ""
	ideaUpdateTitle = ""
	ideaUpdateLabelIDs = ""
	ideaUpdateCampaignID = ""
	ideaUpdateBodyJSON = ""
	ideaUpdateBodyFile = ""
	ideaUpdateCmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}
