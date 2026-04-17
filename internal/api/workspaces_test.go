package api

import (
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestWorkspaceService_Get(t *testing.T) {
	svc := NewWorkspaceService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/workspace" {
			t.Errorf("expected /workspace, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, models.Workspace{ID: "ws1", Name: "My Team"}), nil
	}))

	ws, err := svc.Get()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ws.Name != "My Team" {
		t.Errorf("expected 'My Team', got %s", ws.Name)
	}
}
