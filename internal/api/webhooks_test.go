package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestWebhookService_CRUD(t *testing.T) {
	svc := NewWebhookService(newTestClient(func(r *http.Request) (*http.Response, error) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/webhooks":
			return jsonResponse(t, http.StatusOK, []models.Webhook{{ID: "w1", Name: "n", URL: "https://x", Topics: []string{"post.created"}}}), nil
		case r.Method == http.MethodGet && r.URL.Path == "/webhooks/w1":
			return jsonResponse(t, http.StatusOK, models.Webhook{ID: "w1", Name: "n"}), nil
		case r.Method == http.MethodPost && r.URL.Path == "/webhooks":
			var body models.CreateWebhookRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decoding: %v", err)
			}
			return jsonResponse(t, http.StatusOK, models.Webhook{ID: "w2", Name: body.Name, URL: body.URL, Topics: body.Topics}), nil
		case r.Method == http.MethodPatch && r.URL.Path == "/webhooks/w1":
			return jsonResponse(t, http.StatusOK, models.Webhook{ID: "w1", Name: "updated"}), nil
		case r.Method == http.MethodDelete && r.URL.Path == "/webhooks/w1":
			return jsonResponse(t, http.StatusOK, map[string]bool{"success": true}), nil
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		return jsonResponse(t, http.StatusBadRequest, nil), nil
	}))

	if list, err := svc.List(); err != nil || len(list) != 1 {
		t.Fatalf("list failed: %v", err)
	}
	if _, err := svc.Get("w1"); err != nil {
		t.Fatalf("get: %v", err)
	}
	created, err := svc.Create(models.CreateWebhookRequest{Name: "n2", URL: "https://y", Topics: []string{"post.published"}})
	if err != nil || created.ID != "w2" {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Update("w1", map[string]interface{}{"name": "updated"}); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := svc.Delete("w1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
