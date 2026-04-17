package api

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ordinal-cli/ordinal/internal/models"
)

func TestAnalyticsService_CpmGet(t *testing.T) {
	svc := NewAnalyticsService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/analytics/cpm" {
			t.Errorf("expected /analytics/cpm, got %s", r.URL.Path)
		}
		return jsonResponse(t, http.StatusOK, models.CpmValues{LinkedIn: 30, X: 10, Instagram: 10, Facebook: 10, Threads: 10}), nil
	}))

	v, err := svc.GetCpm()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.LinkedIn != 30 {
		t.Errorf("expected 30, got %v", v.LinkedIn)
	}
}

func TestAnalyticsService_CpmUpdate(t *testing.T) {
	svc := NewAnalyticsService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		var body map[string]float64
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decoding: %v", err)
		}
		if body["linkedIn"] != 25 {
			t.Errorf("expected linkedIn=25, got %v", body["linkedIn"])
		}
		return jsonResponse(t, http.StatusOK, models.CpmValues{LinkedIn: 25, X: 10, Instagram: 10, Facebook: 10, Threads: 10}), nil
	}))

	li := 25.0
	v, err := svc.UpdateCpm(models.CpmUpdateRequest{LinkedIn: &li})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.LinkedIn != 25 {
		t.Errorf("expected 25, got %v", v.LinkedIn)
	}
}

func TestAnalyticsService_Followers(t *testing.T) {
	svc := NewAnalyticsService(newTestClient(func(r *http.Request) (*http.Response, error) {
		want := "/analytics/linkedin/prof-1/followers"
		if r.URL.Path != want {
			t.Errorf("expected %s, got %s", want, r.URL.Path)
		}
		if r.URL.Query().Get("startDate") != "2026-01-01" {
			t.Errorf("expected startDate=2026-01-01, got %s", r.URL.Query().Get("startDate"))
		}
		return jsonResponse(t, http.StatusOK, []models.FollowerDataPoint{{FollowerCount: 100, RecordedAt: "2026-01-01T00:00:00Z"}}), nil
	}))

	points, err := svc.LinkedInFollowers("prof-1", models.AnalyticsDateRange{StartDate: "2026-01-01"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(points) != 1 || points[0].FollowerCount != 100 {
		t.Errorf("unexpected points: %+v", points)
	}
}
