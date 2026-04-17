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

// TestIdeaService_ListAll_EmptyNextCursor locks in that a server response
// of hasMore=true with an empty nextCursor fails fast.
func TestIdeaService_ListAll_EmptyNextCursor(t *testing.T) {
	svc := NewIdeaService(newTestClient(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(t, http.StatusOK, models.IdeaListResponse{
			Ideas:      []models.Idea{{ID: "i1"}},
			NextCursor: "",
			HasMore:    true,
		}), nil
	}))

	_, err := svc.ListAll(models.ListIdeasParams{})
	if err == nil {
		t.Fatal("expected error for hasMore=true with empty nextCursor")
	}
}

// TestIdeaService_ListAll_RepeatedCursor locks in that a server returning
// the same cursor twice breaks out with an error instead of spinning forever.
func TestIdeaService_ListAll_RepeatedCursor(t *testing.T) {
	var calls int
	svc := NewIdeaService(newTestClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			return jsonResponse(t, http.StatusOK, models.IdeaListResponse{
				Ideas:      []models.Idea{{ID: "i1"}},
				NextCursor: "c1",
				HasMore:    true,
			}), nil
		}
		return jsonResponse(t, http.StatusOK, models.IdeaListResponse{
			Ideas:      []models.Idea{{ID: "i2"}},
			NextCursor: "c1",
			HasMore:    true,
		}), nil
	}))

	_, err := svc.ListAll(models.ListIdeasParams{})
	if err == nil {
		t.Fatal("expected error for repeated cursor")
	}
	if calls != 2 {
		t.Errorf("expected 2 calls before bailing, got %d", calls)
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
