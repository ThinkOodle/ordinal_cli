package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/spf13/pflag"
)

// TestAnalyticsCpmUpdate_BodyAndFlagsMerge locks in that `cpm-update` parses
// --body-json first, then overlays any individual --linkedin/--x/etc. flags
// the caller set. The pre-fix `if body != nil { ... } else { flags ... }`
// split meant passing both silently dropped every flag value, which clashed
// with the help text's "overrides individual flags" phrasing in both
// directions (flags were in fact ignored, not the other way around).
func TestAnalyticsCpmUpdate_BodyAndFlagsMerge(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	var captured map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &captured)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"linkedIn":0,"x":0,"instagram":0,"facebook":0,"threads":0}`))
	}))
	defer server.Close()

	prev := testClientOpts
	testClientOpts = []client.Option{
		client.WithBaseURL(server.URL),
		client.WithHTTPClient(server.Client()),
	}
	defer func() { testClientOpts = prev }()

	tests := []struct {
		name string
		args []string
		want map[string]float64
	}{
		{
			name: "flag overrides matching body key",
			args: []string{"analytics", "cpm-update", "--body-json", `{"x":1,"linkedIn":5}`, "--linkedin", "2"},
			want: map[string]float64{"x": 1, "linkedIn": 2},
		},
		{
			name: "flag adds key missing from body",
			args: []string{"analytics", "cpm-update", "--body-json", `{"x":1}`, "--linkedin", "2"},
			want: map[string]float64{"x": 1, "linkedIn": 2},
		},
		{
			name: "body-only request still works",
			args: []string{"analytics", "cpm-update", "--body-json", `{"facebook":3.5}`},
			want: map[string]float64{"facebook": 3.5},
		},
		{
			name: "flags-only request still works",
			args: []string{"analytics", "cpm-update", "--threads", "4"},
			want: map[string]float64{"threads": 4},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			captured = nil
			resetAnalyticsCpmUpdateFlags(t)

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs(tc.args)

			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("cpm-update: %v", err)
			}
			if captured == nil {
				t.Fatal("expected captured request body")
			}
			for k, want := range tc.want {
				got, ok := captured[k].(float64)
				if !ok {
					t.Errorf("missing %q in request body: %+v", k, captured)
					continue
				}
				if got != want {
					t.Errorf("%q = %v, want %v", k, got, want)
				}
			}
			for _, k := range []string{"linkedIn", "x", "instagram", "facebook", "threads"} {
				if _, expected := tc.want[k]; expected {
					continue
				}
				if _, present := captured[k]; present {
					t.Errorf("unexpected key %q in request body (omitempty should suppress unset fields): %+v", k, captured)
				}
			}
		})
	}
}

// TestAnalyticsCpmUpdate_RejectsEmptyUpdate guards the local fast-fail when
// neither flags nor body provide a platform value. Without this, the API
// would be called with an empty body and respond with a less-actionable
// error.
func TestAnalyticsCpmUpdate_RejectsEmptyUpdate(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	for _, tc := range []struct {
		name string
		args []string
	}{
		{"no flags, no body", []string{"analytics", "cpm-update"}},
		{"empty body object", []string{"analytics", "cpm-update", "--body-json", `{}`}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resetAnalyticsCpmUpdateFlags(t)

			var buf bytes.Buffer
			rootCmd.SetOut(&buf)
			rootCmd.SetErr(&buf)
			rootCmd.SetArgs(tc.args)

			err := rootCmd.Execute()
			if err == nil {
				t.Fatal("expected error for empty cpm update")
			}
		})
	}
}

func resetAnalyticsCpmUpdateFlags(t *testing.T) {
	t.Helper()
	analyticsCpmLinkedIn = 0
	analyticsCpmX = 0
	analyticsCpmInstagram = 0
	analyticsCpmFacebook = 0
	analyticsCpmThreads = 0
	analyticsCpmUpdateBodyJSON = ""
	analyticsCpmUpdateBodyFile = ""
	// Cobra does not clear the per-flag "changed" bit across Execute calls on
	// a shared rootCmd, so clear it here or flags set in a prior subtest
	// would leak into the next one via cmd.Flags().Changed(...).
	analyticsCpmUpdateCmd.Flags().VisitAll(func(f *pflag.Flag) {
		f.Changed = false
	})
}
