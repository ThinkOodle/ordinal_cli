package api

import (
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestIdeaService_List(t *testing.T) {
	svc := NewIdeaService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/ideas" {
			t.Errorf("expected /ideas, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, models.IdeaListResponse{
			Ideas:   []models.Idea{{ID: "i1", Title: "I1"}},
			HasMore: false,
		}), nil
	}))

	resp, err := svc.List(models.ListIdeasParams{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Ideas) != 1 {
		t.Errorf("expected 1 idea, got %d", len(resp.Ideas))
	}
}

func TestIdeaService_AddToCalendar(t *testing.T) {
	svc := NewIdeaService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/ideas/abc/add-to-calendar" {
			t.Errorf("expected /ideas/abc/add-to-calendar, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, map[string]string{"id": "abc"}), nil
	}))

	if _, err := svc.AddToCalendar("abc", models.AddIdeaToCalendarRequest{PublishDate: "2026-05-01"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
