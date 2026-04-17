package cmd

import (
	"bytes"
	"strings"
	"testing"
)

// TestEngagementCreate_RejectsEmptyArray locks in that an empty engagements
// array fails locally instead of wasting an API round-trip on an obviously
// invalid request, matching the fast-fail already in subscriber/webhook create.
func TestEngagementCreate_RejectsEmptyArray(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	resetEngagementCreateFlags(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{
		"engagement", "create",
		"--post-id", "p-1",
		"--channel", "LinkedIn",
		"--body-json", `[]`,
	})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for empty engagements array")
	}
	if !strings.Contains(err.Error(), "engagement") {
		t.Errorf("expected error about engagements, got: %v", err)
	}
}

// TestEngagementUpdate_RejectsEmptyBody guards parity with post/idea update:
// an empty {} body is a no-op and must fail locally rather than reach the API.
func TestEngagementUpdate_RejectsEmptyBody(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	for _, tc := range []struct {
		name string
		args []string
	}{
		{"no body flag", []string{"engagement", "update", "--id", "e-1"}},
		{"empty object", []string{"engagement", "update", "--id", "e-1", "--body-json", `{}`}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			engagementID = ""
			engagementUpdateBodyJSON = ""
			engagementUpdateBodyFile = ""

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

func resetEngagementCreateFlags(t *testing.T) {
	t.Helper()
	engagementPostID = ""
	engagementCreateChannel = ""
	engagementCreateBodyJSON = ""
	engagementCreateBodyFile = ""
}
