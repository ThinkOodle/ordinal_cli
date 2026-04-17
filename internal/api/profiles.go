package api

import (
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

// ProfileService handles profile-related API calls.
type ProfileService struct {
	client *client.Client
}

// NewProfileService creates a new ProfileService.
func NewProfileService(c *client.Client) *ProfileService {
	return &ProfileService{client: c}
}

// ListScheduling returns all scheduling profiles.
func (s *ProfileService) ListScheduling() ([]models.SchedulingProfile, error) {
	data, err := s.client.Get("/profiles/scheduling", nil)
	if err != nil {
		return nil, fmt.Errorf("listing scheduling profiles: %w", err)
	}

	var profiles []models.SchedulingProfile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return nil, fmt.Errorf("parsing scheduling profiles: %w", err)
	}
	return profiles, nil
}

// ListEngagement returns all engagement profiles.
func (s *ProfileService) ListEngagement() ([]models.Profile, error) {
	data, err := s.client.Get("/profiles/engagement", nil)
	if err != nil {
		return nil, fmt.Errorf("listing engagement profiles: %w", err)
	}

	var profiles []models.Profile
	if err := json.Unmarshal(data, &profiles); err != nil {
		return nil, fmt.Errorf("parsing engagement profiles: %w", err)
	}
	return profiles, nil
}
