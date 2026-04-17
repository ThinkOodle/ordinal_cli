package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

// AnalyticsService handles analytics API calls.
type AnalyticsService struct {
	client *client.Client
}

// NewAnalyticsService creates a new AnalyticsService.
func NewAnalyticsService(c *client.Client) *AnalyticsService {
	return &AnalyticsService{client: c}
}

// GetCpm returns the current CPM values.
func (s *AnalyticsService) GetCpm() (*models.CpmValues, error) {
	data, err := s.client.Get("/analytics/cpm", nil)
	if err != nil {
		return nil, fmt.Errorf("getting cpm: %w", err)
	}

	var v models.CpmValues
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("parsing cpm: %w", err)
	}
	return &v, nil
}

// UpdateCpm updates the CPM values.
func (s *AnalyticsService) UpdateCpm(req models.CpmUpdateRequest) (*models.CpmValues, error) {
	data, err := s.client.Put("/analytics/cpm", req)
	if err != nil {
		return nil, fmt.Errorf("updating cpm: %w", err)
	}

	var v models.CpmValues
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("parsing updated cpm: %w", err)
	}
	return &v, nil
}

func dateRangeQuery(r models.AnalyticsDateRange) url.Values {
	q := url.Values{}
	if r.StartDate != "" {
		q.Set("startDate", r.StartDate)
	}
	if r.EndDate != "" {
		q.Set("endDate", r.EndDate)
	}
	return q
}

// LinkedInFollowers returns LinkedIn follower history for a profile.
func (s *AnalyticsService) LinkedInFollowers(profileID string, r models.AnalyticsDateRange) ([]models.FollowerDataPoint, error) {
	data, err := s.client.Get("/analytics/linkedin/"+profileID+"/followers", dateRangeQuery(r))
	if err != nil {
		return nil, fmt.Errorf("getting linkedin followers: %w", err)
	}

	var points []models.FollowerDataPoint
	if err := json.Unmarshal(data, &points); err != nil {
		return nil, fmt.Errorf("parsing linkedin followers: %w", err)
	}
	return points, nil
}

// LinkedInPosts returns LinkedIn post analytics for a profile.
func (s *AnalyticsService) LinkedInPosts(profileID string, r models.AnalyticsDateRange) (json.RawMessage, error) {
	data, err := s.client.Get("/analytics/linkedin/"+profileID+"/posts", dateRangeQuery(r))
	if err != nil {
		return nil, fmt.Errorf("getting linkedin post analytics: %w", err)
	}
	return data, nil
}

// XFollowers returns X follower history for a profile.
func (s *AnalyticsService) XFollowers(profileID string, r models.AnalyticsDateRange) ([]models.FollowerDataPoint, error) {
	data, err := s.client.Get("/analytics/x/"+profileID+"/followers", dateRangeQuery(r))
	if err != nil {
		return nil, fmt.Errorf("getting x followers: %w", err)
	}

	var points []models.FollowerDataPoint
	if err := json.Unmarshal(data, &points); err != nil {
		return nil, fmt.Errorf("parsing x followers: %w", err)
	}
	return points, nil
}

// XPosts returns X post analytics for a profile.
func (s *AnalyticsService) XPosts(profileID string, r models.AnalyticsDateRange) (json.RawMessage, error) {
	data, err := s.client.Get("/analytics/x/"+profileID+"/posts", dateRangeQuery(r))
	if err != nil {
		return nil, fmt.Errorf("getting x post analytics: %w", err)
	}
	return data, nil
}
