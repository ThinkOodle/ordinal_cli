package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// TestWebhookCreate_RejectsEmptyTopics guards against Cobra's
// MarkFlagRequired being satisfied by presence alone. A value of "," or
// "  ,  " would pass that check but collapse to an empty slice through
// splitCSV, sending an empty topics array to the API.
func TestWebhookCreate_RejectsEmptyTopics(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	tests := []struct {
		name   string
		topics string
	}{
		{"just commas", ","},
		{"whitespace entries", "  ,  ,   "},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			webhookCreateName = ""
			webhookCreateURL = ""
			webhookCreateDescription = ""
			webhookCreateTopics = ""
			webhookCreateHeadersJSON = ""
			cfgOutputFormat = ""

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{
				"webhook", "create",
				"--name", "n",
				"--url", "https://example.com/h",
				"--topics", tc.topics,
			})

			err := rootCmd.Execute()
			if err == nil {
				t.Fatalf("expected error for --topics=%q", tc.topics)
			}
			if !strings.Contains(err.Error(), "topics") {
				t.Errorf("expected error about topics, got: %v", err)
			}
		})
	}
}

// TestWebhookUpdate_RejectsEmptyBody guards parity with post/idea update: an
// empty {} body is a no-op and must fail locally rather than reach the API.
func TestWebhookUpdate_RejectsEmptyBody(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	for _, tc := range []struct {
		name string
		args []string
	}{
		{"no body flag", []string{"webhook", "update", "--id", "w-1"}},
		{"empty object", []string{"webhook", "update", "--id", "w-1", "--body-json", `{}`}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			webhookID = ""
			webhookUpdateBodyJSON = ""
			webhookUpdateBodyFile = ""

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs(tc.args)

			err := rootCmd.Execute()
			if err == nil {
				t.Fatal("expected error for empty update body")
			}
			if !strings.Contains(err.Error(), "update") && !strings.Contains(err.Error(), "body") {
				t.Errorf("expected update/body error, got: %v", err)
			}
		})
	}
}
