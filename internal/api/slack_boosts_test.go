package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestSlackBoostService_ListByPost(t *testing.T) {
	svc := NewSlackBoostService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/posts/p1/slack-boosts" {
			t.Errorf("expected /posts/p1/slack-boosts, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, models.SlackBoostListResponse{
			SlackBoosts: []models.SlackBoost{{ID: "sb1", PostID: "p1"}},
		}), nil
	}))

	resp, err := svc.ListByPost("p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.SlackBoosts) != 1 {
		t.Errorf("expected 1 slack boost, got %d", len(resp.SlackBoosts))
	}
}

// TestSlackBoostService_CreateGetUpdateDelete uses the documented
// {"slackBoost": ...} envelope shape for create, get, and update
// (see https://docs.tryordinal.com/api-reference/slack-boosts/update-slack-boost
// and the matching create/get operations in the OpenAPI spec) so the test
// fails if any of those services regresses to unwrapped decoding.
func TestSlackBoostService_CreateGetUpdateDelete(t *testing.T) {
	svc := NewSlackBoostService(newTestClient(func(r *http.Request) (*http.Response, error) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/slack-boosts":
			var body models.CreateSlackBoostRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decoding: %v", err)
			}
			return jsonResponse(t, http.StatusOK, map[string]interface{}{
				"slackBoost": models.SlackBoost{ID: "sb1", PostID: body.PostID},
			}), nil
		case r.Method == http.MethodGet && r.URL.Path == "/slack-boosts/sb1":
			return jsonResponse(t, http.StatusOK, map[string]interface{}{
				"slackBoost": models.SlackBoost{ID: "sb1"},
			}), nil
		case r.Method == http.MethodPatch && r.URL.Path == "/slack-boosts/sb1":
			return jsonResponse(t, http.StatusOK, map[string]interface{}{
				"slackBoost": models.SlackBoost{ID: "sb1", Copy: "updated"},
			}), nil
		case r.Method == http.MethodDelete && r.URL.Path == "/slack-boosts/sb1":
			return jsonResponse(t, http.StatusOK, map[string]bool{"success": true}), nil
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		return jsonResponse(t, http.StatusBadRequest, nil), nil
	}))

	created, err := svc.Create(models.CreateSlackBoostRequest{PostID: "p1", SlackWebhookID: "sw1"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID != "sb1" || created.PostID != "p1" {
		t.Errorf("unexpected created boost: %+v", created)
	}
	got, err := svc.Get("sb1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.ID != "sb1" {
		t.Errorf("unexpected get boost: %+v", got)
	}
	updated, err := svc.Update("sb1", map[string]interface{}{"copy": "updated"})
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.ID != "sb1" || updated.Copy != "updated" {
		t.Errorf("unexpected updated boost: %+v", updated)
	}
	data, err := svc.Delete("sb1")
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	var ack map[string]bool
	if err := json.Unmarshal(data, &ack); err != nil {
		t.Fatalf("parse delete body: %v", err)
	}
	if !ack["success"] {
		t.Errorf("expected real API body with success=true; got %v", ack)
	}
}

func TestSlackWebhookService_List(t *testing.T) {
	svc := NewSlackWebhookService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/slack-webhooks" {
			t.Errorf("expected /slack-webhooks, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, models.SlackWebhookListResponse{
			SlackWebhooks: []models.SlackWebhook{{ID: "sw1", ChannelName: "#marketing"}},
		}), nil
	}))

	resp, err := svc.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.SlackWebhooks) != 1 {
		t.Errorf("expected 1 slack webhook, got %d", len(resp.SlackWebhooks))
	}
}
