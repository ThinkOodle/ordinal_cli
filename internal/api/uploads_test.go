package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestUploadService_Create(t *testing.T) {
	svc := NewUploadService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/uploads" {
			t.Errorf("expected /uploads, got %s", r.URL.Path)
		}
		var body models.CreateUploadRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding: %v", err)
		}
		if body.URL != "https://example.com/a.png" {
			t.Errorf("unexpected url: %s", body.URL)
		}
		return jsonResponse(t, http.StatusOK, map[string]string{"id": "up1", "status": "pending"}), nil
	}))

	raw, err := svc.Create(models.CreateUploadRequest{URL: "https://example.com/a.png"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]string
	_ = json.Unmarshal(raw, &out)
	if out["status"] != "pending" {
		t.Errorf("expected status=pending, got %s", out["status"])
	}
}

func TestUploadService_Get(t *testing.T) {
	svc := NewUploadService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/uploads/up1" {
			t.Errorf("expected /uploads/up1, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, map[string]string{"id": "up1", "status": "ready"}), nil
	}))

	_, err := svc.Get("up1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
