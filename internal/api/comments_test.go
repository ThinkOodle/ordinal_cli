package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestCommentService_List(t *testing.T) {
	svc := NewCommentService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/posts/p1/comments" {
			t.Errorf("expected /posts/p1/comments, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, models.CommentListResponse{
			Comments: []models.Comment{{ID: "c1", Message: "hi"}},
		}), nil
	}))

	resp, err := svc.List("p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Comments) != 1 || resp.Comments[0].Message != "hi" {
		t.Errorf("unexpected comments: %+v", resp.Comments)
	}
}

func TestCommentService_Create(t *testing.T) {
	svc := NewCommentService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/posts/p1/comments" {
			t.Errorf("expected /posts/p1/comments, got %s", r.URL.Path)
		}
		var body models.CreateCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		return jsonResponse(t, http.StatusOK, models.Comment{ID: "c2", Message: body.Message}), nil
	}))

	c, err := svc.Create("p1", models.CreateCommentRequest{Message: "ok"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Message != "ok" {
		t.Errorf("expected 'ok', got %s", c.Message)
	}
}

func TestCommentService_Delete(t *testing.T) {
	svc := NewCommentService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/comments/c1" {
			t.Errorf("expected /comments/c1, got %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		// /comments/{id} DELETE returns the deleted Comment.
		return jsonResponse(t, http.StatusOK, models.Comment{ID: "c1", Message: "deleted"}), nil
	}))

	data, err := svc.Delete("c1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got models.Comment
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("parse delete body: %v", err)
	}
	if got.ID != "c1" {
		t.Errorf("expected real comment body to be forwarded; got %+v", got)
	}
}
