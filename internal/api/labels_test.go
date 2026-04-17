package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestLabelService_List(t *testing.T) {
	svc := NewLabelService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/labels" {
			t.Errorf("expected /labels, got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		return jsonResponse(t, http.StatusOK, models.LabelListResponse{
			Labels: []models.Label{{ID: "1", Name: "A", Color: "red"}, {ID: "2", Name: "B", Color: "green"}},
		}), nil
	}))

	resp, err := svc.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Labels) != 2 {
		t.Fatalf("expected 2 labels, got %d", len(resp.Labels))
	}
}

func TestLabelService_Create(t *testing.T) {
	svc := NewLabelService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body models.CreateLabelRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding request: %v", err)
		}
		if body.Name != "Thought" || body.Color != "purple" {
			t.Errorf("unexpected body: %+v", body)
		}
		return jsonResponse(t, http.StatusOK, models.Label{ID: "lbl-1", Name: body.Name, Color: body.Color}), nil
	}))

	label, err := svc.Create(models.CreateLabelRequest{Name: "Thought", Color: "purple"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if label.ID != "lbl-1" {
		t.Errorf("expected lbl-1, got %s", label.ID)
	}
}

func TestLabelService_Delete(t *testing.T) {
	svc := NewLabelService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/labels/lbl-1" {
			t.Errorf("expected /labels/lbl-1, got %s", r.URL.Path)
		}
		// /labels/{id} DELETE returns the deleted Label per the OpenAPI
		// spec; the CLI must forward that body unchanged rather than
		// fabricating an acknowledgement.
		return jsonResponse(t, http.StatusOK, models.Label{ID: "lbl-1", Name: "Thought", Color: "purple"}), nil
	}))

	data, err := svc.Delete("lbl-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got models.Label
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("parse delete body: %v", err)
	}
	if got.ID != "lbl-1" {
		t.Errorf("expected delete body to carry the label; got %+v", got)
	}
}
