package api

import (
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestInlineCommentService_List(t *testing.T) {
	svc := NewInlineCommentService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/posts/p1/inline-comments" {
			t.Errorf("expected /posts/p1/inline-comments, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("channel") != "LinkedIn" {
			t.Errorf("expected channel=LinkedIn, got %s", r.URL.Query().Get("channel"))
		}
		if r.URL.Query().Get("resolved") != "true" {
			t.Errorf("expected resolved=true, got %s", r.URL.Query().Get("resolved"))
		}
		return jsonResponse(t, http.StatusOK, models.InlineCommentListResponse{
			InlineComments: []models.InlineComment{{ID: "ic1", Resolved: true, Channel: "LinkedIn"}},
		}), nil
	}))

	resolved := true
	resp, err := svc.List("p1", models.ListInlineCommentsParams{Channel: "LinkedIn", Resolved: &resolved})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.InlineComments) != 1 {
		t.Errorf("expected 1 inline comment, got %d", len(resp.InlineComments))
	}
}
