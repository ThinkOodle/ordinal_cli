package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestInviteService_CRUD(t *testing.T) {
	svc := NewInviteService(newTestClient(func(r *http.Request) (*http.Response, error) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/invites":
			return jsonResponse(t, http.StatusOK, []models.Invite{{ID: "i1", Email: "a@example.com"}}), nil
		case r.Method == http.MethodPost && r.URL.Path == "/invites":
			var body models.CreateInviteRequest
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Fatalf("decoding: %v", err)
			}
			return jsonResponse(t, http.StatusOK, models.Invite{ID: "i2", Email: body.Email}), nil
		case r.Method == http.MethodDelete && r.URL.Path == "/invites/i1":
			return jsonResponse(t, http.StatusOK, map[string]bool{"success": true}), nil
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		return jsonResponse(t, http.StatusBadRequest, nil), nil
	}))

	if list, err := svc.List(); err != nil || len(list) != 1 {
		t.Fatalf("list: %v", err)
	}
	created, err := svc.Create(models.CreateInviteRequest{Email: "new@example.com"})
	if err != nil || created.ID != "i2" {
		t.Fatalf("create: %v", err)
	}
	if err := svc.Delete("i1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}
