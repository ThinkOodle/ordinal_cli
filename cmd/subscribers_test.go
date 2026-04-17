package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// TestSubscriberCreate_RejectsEmptyUserIDs guards against Cobra's
// MarkFlagRequired being satisfied by presence alone. A value of "," or
// "  ,  " would pass that check but collapse to an empty slice through
// splitCSV, sending an empty userIds array to the API.
func TestSubscriberCreate_RejectsEmptyUserIDs(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	tests := []struct {
		name    string
		userIDs string
	}{
		{"just commas", ","},
		{"whitespace entries", "  ,  ,   "},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			subscriberPostID = ""
			subscriberUserIDs = ""
			cfgOutputFormat = ""

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs([]string{
				"subscriber", "create",
				"--post-id", "p-1",
				"--user-ids", tc.userIDs,
			})

			err := rootCmd.Execute()
			if err == nil {
				t.Fatalf("expected error for --user-ids=%q", tc.userIDs)
			}
			if !strings.Contains(err.Error(), "user-ids") {
				t.Errorf("expected error about user-ids, got: %v", err)
			}
		})
	}
}
