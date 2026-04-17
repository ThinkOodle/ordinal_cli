package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ordinal-cli/ordinal/internal/client"
)

// InstagramService handles Instagram API calls.
type InstagramService struct {
	client *client.Client
}

// NewInstagramService creates a new InstagramService.
func NewInstagramService(c *client.Client) *InstagramService {
	return &InstagramService{client: c}
}

// SearchLocations searches Instagram locations by query.
func (s *InstagramService) SearchLocations(query string) (json.RawMessage, error) {
	q := url.Values{}
	q.Set("query", query)

	data, err := s.client.Get("/instagram/locations/search", q)
	if err != nil {
		return nil, fmt.Errorf("searching instagram locations: %w", err)
	}
	return data, nil
}
