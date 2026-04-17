package api

import (
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

// EngagementService handles engagement-related API calls.
type EngagementService struct {
	client *client.Client
}

// NewEngagementService creates a new EngagementService.
func NewEngagementService(c *client.Client) *EngagementService {
	return &EngagementService{client: c}
}

// List returns the auto-engagements configured for a post.
func (s *EngagementService) List(postID string) (*models.EngagementListResponse, error) {
	data, err := s.client.Get("/posts/"+postID+"/engagements", nil)
	if err != nil {
		return nil, fmt.Errorf("listing engagements: %w", err)
	}

	var resp models.EngagementListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing engagements list: %w", err)
	}
	return &resp, nil
}

// Create creates engagements on a post.
func (s *EngagementService) Create(postID string, req models.CreateEngagementsRequest) (json.RawMessage, error) {
	data, err := s.client.Post("/posts/"+postID+"/engagements", req)
	if err != nil {
		return nil, fmt.Errorf("creating engagements: %w", err)
	}
	return data, nil
}

// Update updates a single engagement by ID.
func (s *EngagementService) Update(id string, req models.UpdateEngagementRequest) (json.RawMessage, error) {
	data, err := s.client.Patch("/engagements/"+id, req)
	if err != nil {
		return nil, fmt.Errorf("updating engagement: %w", err)
	}
	return data, nil
}

// Delete deletes an engagement by ID.
func (s *EngagementService) Delete(id string) error {
	if _, err := s.client.Delete("/engagements/" + id); err != nil {
		return fmt.Errorf("deleting engagement: %w", err)
	}
	return nil
}
