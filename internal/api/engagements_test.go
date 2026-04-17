package api

import (
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestEngagementService_List(t *testing.T) {
	svc := NewEngagementService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/posts/p1/engagements" {
			t.Errorf("expected /posts/p1/engagements, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, models.EngagementListResponse{
			Engagements: []models.Engagement{{ID: "e1", Channel: "LinkedIn", Type: "Like"}},
		}), nil
	}))

	resp, err := svc.List("p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Engagements) != 1 {
		t.Errorf("expected 1 engagement, got %d", len(resp.Engagements))
	}
}

func TestEngagementService_Create(t *testing.T) {
	svc := NewEngagementService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/posts/p1/engagements" {
			t.Errorf("expected /posts/p1/engagements, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, map[string]string{"ok": "true"}), nil
	}))

	_, err := svc.Create("p1", models.CreateEngagementsRequest{
		Channel:     "LinkedIn",
		Engagements: []models.EngagementInput{{"type": "Like", "profileId": "pr1"}},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEngagementService_UpdateDelete(t *testing.T) {
	paths := map[string]string{}
	svc := NewEngagementService(newTestClient(func(r *http.Request) (*http.Response, error) {
		paths[r.URL.Path] = r.Method
		return jsonResponse(t, http.StatusOK, map[string]string{"id": "e1"}), nil
	}))

	if _, err := svc.Update("e1", map[string]interface{}{"copy": "hi"}); err != nil {
		t.Fatalf("update failed: %v", err)
	}
	if err := svc.Delete("e1"); err != nil {
		t.Fatalf("delete failed: %v", err)
	}
	if paths["/engagements/e1"] == "" {
		t.Errorf("expected /engagements/e1, got %v", paths)
	}
}
