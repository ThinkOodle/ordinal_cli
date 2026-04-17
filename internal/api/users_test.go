package api

import (
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestUserService_List(t *testing.T) {
	svc := NewUserService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/users" {
			t.Errorf("expected /users, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, []models.User{{ID: "u1", Email: "a@example.com"}}), nil
	}))

	users, err := svc.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("expected 1 user, got %d", len(users))
	}
}
