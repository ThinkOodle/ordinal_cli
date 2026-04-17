package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/pflag"
)

// TestSlackBoostUpdate_NoOpValidatesBeforeAuth locks in the ordering: when
// no fields are provided, the local "no fields to update" error must surface
// ahead of newClient()'s auth error so users with no API key configured
// still see the actionable local complaint.
func TestSlackBoostUpdate_NoOpValidatesBeforeAuth(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	resetSlackBoostUpdateFlags(t)

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"slack-boost", "update", "--id", "b-1"})

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

func resetSlackBoostUpdateFlags(t *testing.T) {
	t.Helper()
	slackBoostID = ""
	slackBoostUpdateCopy = ""
	slackBoostUpdateBodyJSON = ""
	slackBoostUpdateBodyFile = ""
	slackBoostUpdateCmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}
