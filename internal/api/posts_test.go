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

// TestPostService_Get uses the documented {"post": ...} envelope shape
// (see https://docs.tryordinal.com/api/openapi.json, GET /posts/{id}) so
// the test fails if the service ever regresses to unwrapped decoding.
func TestPostService_Get(t *testing.T) {
	svc := NewPostService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/posts/abc" {
			t.Errorf("expected /posts/abc, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, map[string]interface{}{
			"post": models.Post{ID: "abc", Title: "hello", Status: "ToDo"},
		}), nil
	}))

	p, err := svc.Get("abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != "abc" {
		t.Errorf("expected id abc, got %s", p.ID)
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

// TestPostService_ListAll_EmptyNextCursor locks in that a server response
// of hasMore=true with an empty nextCursor fails fast, instead of silently
// truncating results or looping with the same cursor forever.
func TestPostService_ListAll_EmptyNextCursor(t *testing.T) {
	svc := NewPostService(newTestClient(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(t, http.StatusOK, models.PostListResponse{
			Posts:      []models.Post{{ID: "p1"}},
			NextCursor: "",
			HasMore:    true,
		}), nil
	}))

	_, err := svc.ListAll(models.ListPostsParams{})
	if err == nil {
		t.Fatal("expected error for hasMore=true with empty nextCursor")
	}
}

// TestPostService_ListAll_RepeatedCursor locks in that a server returning
// the same cursor twice breaks out with an error instead of spinning forever.
func TestPostService_ListAll_RepeatedCursor(t *testing.T) {
	var calls int
	svc := NewPostService(newTestClient(func(r *http.Request) (*http.Response, error) {
		calls++
		// First page hands out cursor "c1"; second page echoes the same
		// cursor. The client should stop on the second page.
		if calls == 1 {
			return jsonResponse(t, http.StatusOK, models.PostListResponse{
				Posts:      []models.Post{{ID: "p1"}},
				NextCursor: "c1",
				HasMore:    true,
			}), nil
		}
		return jsonResponse(t, http.StatusOK, models.PostListResponse{
			Posts:      []models.Post{{ID: "p2"}},
			NextCursor: "c1",
			HasMore:    true,
		}), nil
	}))

	_, err := svc.ListAll(models.ListPostsParams{})
	if err == nil {
		t.Fatal("expected error for repeated cursor")
	}
	if calls != 2 {
		t.Errorf("expected 2 calls before bailing, got %d", calls)
	}
}

// TestPostService_ListAll_HonorsLimit locks in that ListAll uses the caller's
// params.Limit as the page size rather than silently forcing 100. Prior
// behavior made --all --limit N behave identically to --all --limit 100.
func TestPostService_ListAll_HonorsLimit(t *testing.T) {
	var gotLimits []string
	svc := NewPostService(newTestClient(func(r *http.Request) (*http.Response, error) {
		gotLimits = append(gotLimits, r.URL.Query().Get("limit"))
		return jsonResponse(t, http.StatusOK, models.PostListResponse{
			Posts:   []models.Post{{ID: "p1"}},
			HasMore: false,
		}), nil
	}))

	if _, err := svc.ListAll(models.ListPostsParams{Limit: 10}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(gotLimits) != 1 || gotLimits[0] != "10" {
		t.Errorf("expected limit=10 on request, got %v", gotLimits)
	}
}

// TestPostService_ListAll_TwoCursorCycle locks in that a server alternating
// between two cursors (A -> B -> A -> B ...) with hasMore=true is detected
// and aborted instead of looping forever. The adjacent-cursor guard alone
// would not catch this.
func TestPostService_ListAll_TwoCursorCycle(t *testing.T) {
	var calls int
	svc := NewPostService(newTestClient(func(r *http.Request) (*http.Response, error) {
		calls++
		next := "B"
		if calls%2 == 0 {
			next = "A"
		}
		return jsonResponse(t, http.StatusOK, models.PostListResponse{
			Posts:      []models.Post{{ID: "p"}},
			NextCursor: next,
			HasMore:    true,
		}), nil
	}))

	_, err := svc.ListAll(models.ListPostsParams{})
	if err == nil {
		t.Fatal("expected error for two-cursor cycle")
	}
	// Calls: 1 -> next=B (seen), 2 -> next=A (seen), 3 -> next=B (repeat).
	if calls != 3 {
		t.Errorf("expected 3 calls before detecting cycle, got %d", calls)
	}
}

func TestPostService_Delete(t *testing.T) {
	var method string
	svc := NewPostService(newTestClient(func(r *http.Request) (*http.Response, error) {
		method = r.Method
		if r.URL.Path != "/posts/abc" {
			t.Errorf("expected /posts/abc, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, map[string]string{"id": "abc"}), nil
	}))

	if _, err := svc.Delete("abc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if method != http.MethodDelete {
		t.Errorf("expected DELETE, got %s", method)
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
