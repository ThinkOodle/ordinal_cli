package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestPostService_List(t *testing.T) {
	svc := NewPostService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/posts" {
			t.Errorf("expected /posts, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("status") != "Scheduled" {
			t.Errorf("expected status=Scheduled, got %s", r.URL.Query().Get("status"))
		}
		return jsonResponse(t, http.StatusOK, models.PostListResponse{
			Posts:      []models.Post{{ID: "p1", Title: "t1", Status: "Scheduled"}},
			NextCursor: "cursor-2",
			HasMore:    true,
		}), nil
	}))

	resp, err := svc.List(models.ListPostsParams{Limit: 10, Status: "Scheduled"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Posts) != 1 || resp.Posts[0].ID != "p1" {
		t.Errorf("unexpected posts: %+v", resp.Posts)
	}
	if resp.NextCursor != "cursor-2" || !resp.HasMore {
		t.Errorf("unexpected pagination: %+v", resp)
	}
}

func TestPostService_Get(t *testing.T) {
	svc := NewPostService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/posts/abc" {
			t.Errorf("expected /posts/abc, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, models.Post{ID: "abc", Title: "hello", Status: "ToDo"}), nil
	}))

	p, err := svc.Get("abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Title != "hello" {
		t.Errorf("expected hello, got %s", p.Title)
	}
}

func TestPostService_Create(t *testing.T) {
	svc := NewPostService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body["title"] != "Launch" {
			t.Errorf("unexpected title: %v", body["title"])
		}
		return jsonResponse(t, http.StatusOK, map[string]string{"id": "new-id", "title": "Launch"}), nil
	}))

	raw, err := svc.Create(map[string]interface{}{"title": "Launch", "publishAt": "2026-05-01T10:00:00Z", "status": "ToDo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]string
	_ = json.Unmarshal(raw, &out)
	if out["id"] != "new-id" {
		t.Errorf("expected new-id, got %s", out["id"])
	}
}

func TestPostService_Schedule(t *testing.T) {
	svc := NewPostService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/posts/abc/schedule" {
			t.Errorf("expected /posts/abc/schedule, got %s", r.URL.Path)
		}
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body["publishAt"] != "2026-05-01T10:00:00Z" {
			t.Errorf("unexpected publishAt: %v", body["publishAt"])
		}
		return jsonResponse(t, http.StatusOK, map[string]string{"id": "abc"}), nil
	}))

	if _, err := svc.Schedule("abc", models.SchedulePostRequest{PublishAt: "2026-05-01T10:00:00Z"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostService_ArchiveUnarchiveUnschedule(t *testing.T) {
	paths := map[string]string{}
	svc := NewPostService(newTestClient(func(r *http.Request) (*http.Response, error) {
		paths[r.URL.Path] = r.Method
		return jsonResponse(t, http.StatusOK, map[string]string{"id": "abc"}), nil
	}))

	for _, fn := range []func(string) (json.RawMessage, error){svc.Archive, svc.Unarchive, svc.Unschedule} {
		if _, err := fn("abc"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	for _, want := range []string{"/posts/abc/archive", "/posts/abc/unarchive", "/posts/abc/unschedule"} {
		if paths[want] != http.MethodPost {
			t.Errorf("missing POST %s, got paths=%v", want, paths)
		}
	}
}
