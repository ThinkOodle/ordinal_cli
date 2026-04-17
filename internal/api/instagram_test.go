package api

import (
	"net/http"
	"testing"
)

func TestInstagramService_SearchLocations(t *testing.T) {
	svc := NewInstagramService(newTestClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/instagram/locations/search" {
			t.Errorf("expected /instagram/locations/search, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("query") != "Brooklyn" {
			t.Errorf("expected query=Brooklyn, got %s", r.URL.Query().Get("query"))
		}
		return jsonResponse(t, http.StatusOK, map[string]interface{}{"locations": []interface{}{}}), nil
	}))

	_, err := svc.SearchLocations("Brooklyn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
