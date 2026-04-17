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
			return jsonResponse(t, http.StatusOK, models.CreateInviteResponse{
				Invite:    &models.Invite{ID: "i2", Email: body.Email},
				SentEmail: true,
			}), nil
		case r.Method == http.MethodDelete && r.URL.Path == "/invites/i1":
			// /invites/{id} DELETE returns `{"deleted": true}` in the spec.
			return jsonResponse(t, http.StatusOK, map[string]bool{"deleted": true}), nil
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		return jsonResponse(t, http.StatusBadRequest, nil), nil
	}))

	if list, err := svc.List(); err != nil || len(list) != 1 {
		t.Fatalf("list: %v", err)
	}
	created, err := svc.Create(models.CreateInviteRequest{Email: "new@example.com"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if !created.SentEmail {
		t.Fatalf("create: expected sentEmail=true, got false")
	}
	if created.Invite == nil || created.Invite.ID != "i2" || created.Invite.Email != "new@example.com" {
		t.Fatalf("create: unexpected invite %+v", created.Invite)
	}
	data, err := svc.Delete("i1")
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	var got map[string]bool
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("parse delete body: %v", err)
	}
	if !got["deleted"] {
		t.Errorf("delete: expected real API body with deleted=true; got %v", got)
	}
}

// TestInviteService_CreateExistingUser covers the documented success path
// where the invited email already belongs to an Ordinal user: the API adds
// them to the workspace directly, so invite is null and sentEmail is false.
func TestInviteService_CreateExistingUser(t *testing.T) {
	svc := NewInviteService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.Method == http.MethodPost && r.URL.Path == "/invites" {
			return jsonResponse(t, http.StatusOK, map[string]interface{}{
				"invite":    nil,
				"sentEmail": false,
			}), nil
		}
		t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		return jsonResponse(t, http.StatusBadRequest, nil), nil
	}))

	resp, err := svc.Create(models.CreateInviteRequest{Email: "existing@example.com"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if resp.Invite != nil {
		t.Fatalf("create: expected nil invite, got %+v", resp.Invite)
	}
	if resp.SentEmail {
		t.Fatalf("create: expected sentEmail=false, got true")
	}
}
