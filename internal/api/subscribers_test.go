package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestSubscriberService_List(t *testing.T) {
	svc := NewSubscriberService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/posts/p1/subscribers" {
			t.Errorf("expected /posts/p1/subscribers, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, models.SubscriberListResponse{
			Subscribers: []models.Subscriber{{ID: "s1", User: models.SubscriberUser{Email: "a@example.com"}}},
		}), nil
	}))

	resp, err := svc.List("p1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Subscribers) != 1 {
		t.Errorf("expected 1 subscriber, got %d", len(resp.Subscribers))
	}
}

func TestSubscriberService_Create(t *testing.T) {
	svc := NewSubscriberService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/subscribers" {
			t.Errorf("expected /subscribers, got %s", r.URL.Path)
		}
		var body models.CreateSubscribersRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding body: %v", err)
		}
		if body.PostID != "p1" || len(body.UserIDs) != 2 {
			t.Errorf("unexpected body: %+v", body)
		}
		return jsonResponse(t, http.StatusOK, map[string]interface{}{"subscribers": []map[string]string{}}), nil
	}))

	if _, err := svc.Create(models.CreateSubscribersRequest{PostID: "p1", UserIDs: []string{"u1", "u2"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
