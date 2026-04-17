package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/client"
)

// TestNormalizeCSV locks in the trim-and-rejoin contract of the helper so the
// cmd layer can forward CSV flag values into query strings without embedded
// whitespace. Keep this in step with splitCSV's semantics: empty entries are
// dropped, surviving entries are trimmed.
func TestNormalizeCSV(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"", ""},
		{"a", "a"},
		{"a,b,c", "a,b,c"},
		{"a, b, c", "a,b,c"},
		{"  a  ,  b  ,  c  ", "a,b,c"},
		{"a,,b", "a,b"},
		{" , , ", ""},
		{",a,", "a"},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			if got := normalizeCSV(tc.in); got != tc.want {
				t.Errorf("normalizeCSV(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

// TestPostList_NormalizesCSVFilters verifies that --ids and --label-ids are
// trimmed and re-joined before being placed in the query string. Before the
// fix, "a, b" went out as "a, b" verbatim and the server silently returned
// nothing matching the space-prefixed entry.
func TestPostList_NormalizesCSVFilters(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	postListLimit = 0
	postListIDs = ""
	postListLabelIDs = ""

	var req *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r.Clone(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"posts":[],"hasMore":false}`))
	}))
	defer server.Close()

	prev := testClientOpts
	testClientOpts = []client.Option{
		client.WithBaseURL(server.URL),
		client.WithHTTPClient(server.Client()),
	}
	defer func() { testClientOpts = prev }()

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{
		"post", "list",
		"--ids", "a, b, ,c",
		"--label-ids", "  x  ,y  ,  z",
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("post list: %v", err)
	}
	if req == nil {
		t.Fatal("expected captured request")
	}
	if got := req.URL.Query().Get("ids"); got != "a,b,c" {
		t.Errorf("ids query = %q, want %q", got, "a,b,c")
	}
	if got := req.URL.Query().Get("labelIds"); got != "x,y,z" {
		t.Errorf("labelIds query = %q, want %q", got, "x,y,z")
	}
}

// TestIdeaList_NormalizesCSVFilters is the idea-list parity assertion for
// the same --ids / --label-ids whitespace handling.
func TestIdeaList_NormalizesCSVFilters(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	ideaListLimit = 0
	ideaListIDs = ""
	ideaListLabelIDs = ""

	var req *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r.Clone(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ideas":[],"hasMore":false}`))
	}))
	defer server.Close()

	prev := testClientOpts
	testClientOpts = []client.Option{
		client.WithBaseURL(server.URL),
		client.WithHTTPClient(server.Client()),
	}
	defer func() { testClientOpts = prev }()

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{
		"idea", "list",
		"--ids", "id-1, id-2 ,id-3",
		"--label-ids", "lbl-1 , lbl-2",
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("idea list: %v", err)
	}
	if req == nil {
		t.Fatal("expected captured request")
	}
	if got := req.URL.Query().Get("ids"); got != "id-1,id-2,id-3" {
		t.Errorf("ids query = %q, want %q", got, "id-1,id-2,id-3")
	}
	if got := req.URL.Query().Get("labelIds"); got != "lbl-1,lbl-2" {
		t.Errorf("labelIds query = %q, want %q", got, "lbl-1,lbl-2")
	}
}

// TestLinkedInLeadsGetLeads_NormalizesTypes verifies the --types flag is
// trimmed and re-joined before being placed in the query string. LinkedIn's
// lead types are a fixed enum set (LIKE, COMMENT, RESHARE) and the server
// does no whitespace tolerance.
func TestLinkedInLeadsGetLeads_NormalizesTypes(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("ORDINAL_API_KEY", "test-key")
	t.Setenv("ORDINAL_OUTPUT_FORMAT", "")
	t.Setenv("ORDINAL_VERBOSE", "")

	llLimit = 0
	llTypes = ""
	llMinFollowerCount = 0

	var req *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req = r.Clone(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"leads":[]}`))
	}))
	defer server.Close()

	prev := testClientOpts
	testClientOpts = []client.Option{
		client.WithBaseURL(server.URL),
		client.WithHTTPClient(server.Client()),
	}
	defer func() { testClientOpts = prev }()

	var buf bytes.Buffer
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{
		"linkedin-leads", "get-leads",
		"--profile-id", "p1",
		"--post-id", "po1",
		"--types", "LIKE, COMMENT , RESHARE",
	})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("linkedin-leads get-leads: %v", err)
	}
	if req == nil {
		t.Fatal("expected captured request")
	}
	if got := req.URL.Query().Get("types"); got != "LIKE,COMMENT,RESHARE" {
		t.Errorf("types query = %q, want %q", got, "LIKE,COMMENT,RESHARE")
	}
}
