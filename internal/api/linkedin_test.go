package api

import (
	"net/http"
	"testing"
)

func TestLinkedInService_GetProfile(t *testing.T) {
	svc := NewLinkedInService(newTestClient(func(r *http.Request) (*http.Response, error) {
		want := "/linkedin/profile/urn:li:person:abc"
		if r.URL.Path != want {
			t.Errorf("expected %s, got %s", want, r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, map[string]string{"urn": "urn:li:person:abc", "name": "A"}), nil
	}))

	data, err := svc.GetProfile("urn:li:person:abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) == 0 {
		t.Errorf("expected payload, got empty")
	}
}

func TestLinkedInService_GetMention(t *testing.T) {
	svc := NewLinkedInService(newTestClient(func(r *http.Request) (*http.Response, error) {
		want := "/linkedin/johndoe/mentions"
		if r.URL.Path != want {
			t.Errorf("expected %s, got %s", want, r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, map[string]interface{}{"user": map[string]string{"urn": "urn:li:person:abc"}}), nil
	}))

	_, err := svc.GetMention("johndoe")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLinkedInLeadsService_ListPosts(t *testing.T) {
	svc := NewLinkedInLeadsService(newTestClient(func(r *http.Request) (*http.Response, error) {
		want := "/linkedin/leads/prof1/posts"
		if r.URL.Path != want {
			t.Errorf("expected %s, got %s", want, r.URL.Path)
		}
		if r.URL.Query().Get("limit") != "50" {
			t.Errorf("expected limit=50, got %s", r.URL.Query().Get("limit"))
		}
		return jsonResponse(t, http.StatusOK, map[string]interface{}{"posts": []interface{}{}, "hasMore": false}), nil
	}))

	_, err := svc.ListPosts("prof1", LinkedInLeadsListPostsParams{Limit: 50})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLinkedInLeadsService_GetLeadsByPost(t *testing.T) {
	svc := NewLinkedInLeadsService(newTestClient(func(r *http.Request) (*http.Response, error) {
		want := "/linkedin/leads/prof1/posts/post1"
		if r.URL.Path != want {
			t.Errorf("expected %s, got %s", want, r.URL.Path)
		}
		if r.URL.Query().Get("types") != "LIKE,COMMENT" {
			t.Errorf("expected types=LIKE,COMMENT, got %s", r.URL.Query().Get("types"))
		}
		return jsonResponse(t, http.StatusOK, map[string]interface{}{"leads": []interface{}{}, "hasMore": false}), nil
	}))

	_, err := svc.GetLeadsByPost("prof1", "post1", LinkedInLeadsGetLeadsParams{Types: "LIKE,COMMENT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
