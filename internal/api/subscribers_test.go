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

// /subscribers/{id} DELETE returns `{"deletedSubscriber": Subscriber}` per
// the OpenAPI spec. The service must forward that body verbatim so the CLI
// can format the real server response rather than a fabricated
// acknowledgement.
func TestSubscriberService_Delete(t *testing.T) {
	svc := NewSubscriberService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/subscribers/s1" {
			t.Errorf("expected /subscribers/s1, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, map[string]interface{}{
			"deletedSubscriber": models.Subscriber{ID: "s1"},
		}), nil
	}))

	data, err := svc.Delete("s1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got struct {
		DeletedSubscriber models.Subscriber `json:"deletedSubscriber"`
	}
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("parse delete body: %v", err)
	}
	if got.DeletedSubscriber.ID != "s1" {
		t.Errorf("expected deleted subscriber id=s1; got %+v", got.DeletedSubscriber)
	}
}
