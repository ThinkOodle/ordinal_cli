package api

import (
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestProfileService_ListScheduling(t *testing.T) {
	svc := NewProfileService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/profiles/scheduling" {
			t.Errorf("expected /profiles/scheduling, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, []models.SchedulingProfile{
			{Profile: models.Profile{ID: "p1", Channel: "LinkedIn"}, IsLeadsScrapingEnabled: true},
		}), nil
	}))

	profiles, err := svc.ListScheduling()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 || profiles[0].ID != "p1" {
		t.Errorf("unexpected profiles: %+v", profiles)
	}
}

func TestProfileService_ListEngagement(t *testing.T) {
	svc := NewProfileService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/profiles/engagement" {
			t.Errorf("expected /profiles/engagement, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, []models.Profile{{ID: "p2", Channel: "Twitter"}}), nil
	}))

	profiles, err := svc.ListEngagement()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Errorf("expected 1 profile, got %d", len(profiles))
	}
}
