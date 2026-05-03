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

// TestPostUpdate_ClearsLabelsWithEmptyFlag locks in that an explicit
// --label-ids "" on update serializes to an empty JSON array, not null.
// The API distinguishes "don't touch labels" (key absent) from "clear all
// labels" (key present with []); sending null depends on undocumented
// server handling and is rejected outright by many validators.
func TestPostUpdate_ClearsLabelsWithEmptyFlag(t *testing.T) {
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
			resetPostUpdateFlags(t)

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{"post", "update", "--id", "p-1", "--label-ids", tc.labelArg})

			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("post update: %v", err)
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

// TestPostCreate_PreservesCurrentAssetObjectShape guards the current API
// request shape for channel assets. The OpenAPI schema now uses
// assets:[{assetId:"..."}], not the older assetIds:["..."] shorthand. The
// CLI must pass this nested JSON through exactly as supplied.
func TestPostCreate_PreservesCurrentAssetObjectShape(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	var captured []byte
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, _ = io.ReadAll(r.Body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"p-1"}`))
	}))
	defer server.Close()

	prev := testClientOpts
	testClientOpts = []client.Option{
		client.WithBaseURL(server.URL),
		client.WithHTTPClient(server.Client()),
	}
	defer func() { testClientOpts = prev }()

	resetPostCreateFlags(t)
	bodyJSON := `{
		"title":"Launch",
		"publishAt":"2026-01-15T14:00:00Z",
		"status":"Scheduled",
		"linkedIn":{"profileId":"li-1","copy":"hello","assets":[{"assetId":"asset-li"}]},
		"x":{"profileId":"x-1","tweets":[{"copy":"hello","assets":[{"assetId":"asset-x"}]}]},
		"instagram":{"profileId":"ig-1","type":"Feed","copy":"hello","assets":[{"assetId":"asset-ig","tags":[{"username":"tryordinal","x":0.5,"y":0.5}]}]},
		"tikTok":{"profileId":"tt-1","copy":"hello","assets":[{"assetId":"asset-tt"}]},
		"youTubeShorts":{"profileId":"yt-1","title":"hello","assets":[{"assetId":"asset-yt"}]}
	}`

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"post", "create", "--body-json", bodyJSON})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("post create: %v", err)
	}

	body := string(captured)
	if strings.Contains(body, `"assetIds"`) {
		t.Fatalf("post body must not use old assetIds shape: %s", body)
	}
	for _, want := range []string{
		`"assets":[{"assetId":"asset-li"}]`,
		`"assets":[{"assetId":"asset-x"}]`,
		`"assetId":"asset-ig"`,
		`"tags":[{"username":"tryordinal","x":0.5,"y":0.5}]`,
		`"assets":[{"assetId":"asset-tt"}]`,
		`"assets":[{"assetId":"asset-yt"}]`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("expected body to contain %s; got %s", want, body)
		}
	}
}

// TestPostUpdate_NoOpValidatesBeforeAuth locks in the ordering: when no
// fields are provided, the local "no fields to update" error must surface
// ahead of newClient()'s auth error. Otherwise a user mistyping the command
// on a machine without a configured API key sees an irrelevant auth
// complaint instead of the actionable validation message.
func TestPostUpdate_NoOpValidatesBeforeAuth(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	resetPostUpdateFlags(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"post", "update", "--id", "p-1"})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "api key") {
		t.Errorf("validation should run before auth; got auth error: %v", err)
	}
	if !strings.Contains(err.Error(), "no fields to update") {
		t.Errorf("expected no-op error, got: %v", err)
	}
}

// resetPostUpdateFlags clears postUpdateCmd's package-level flag variables
// and the pflag "changed" bits so consecutive Execute calls on the shared
// rootCmd don't leak state between tests.
func resetPostUpdateFlags(t *testing.T) {
	t.Helper()
	postID = ""
	postUpdateTitle = ""
	postUpdatePublishAt = ""
	postUpdateStatus = ""
	postUpdateLabelIDs = ""
	postUpdateCampaignID = ""
	postUpdateNotes = ""
	postUpdateBodyJSON = ""
	postUpdateBodyFile = ""
	postUpdateCmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
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
