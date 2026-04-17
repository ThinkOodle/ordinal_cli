package api

import (
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

// SubscriberService handles post-subscriber API calls.
type SubscriberService struct {
	client *client.Client
}

// NewSubscriberService creates a new SubscriberService.
func NewSubscriberService(c *client.Client) *SubscriberService {
	return &SubscriberService{client: c}
}

// List returns the subscribers on a post.
func (s *SubscriberService) List(postID string) (*models.SubscriberListResponse, error) {
	data, err := s.client.Get("/posts/"+postID+"/subscribers", nil)
	if err != nil {
		return nil, fmt.Errorf("listing subscribers: %w", err)
	}

	var resp models.SubscriberListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing subscribers list: %w", err)
	}
	return &resp, nil
}

// Create adds subscribers to a post.
func (s *SubscriberService) Create(req models.CreateSubscribersRequest) (json.RawMessage, error) {
	data, err := s.client.Post("/subscribers", req)
	if err != nil {
		return nil, fmt.Errorf("adding subscribers: %w", err)
	}
	return data, nil
}

// Delete removes a subscriber by ID. The API returns a
// `{"deletedSubscriber": Subscriber}` envelope; we forward it verbatim so
// callers see the real server response instead of a fabricated acknowledgement.
func (s *SubscriberService) Delete(id string) (json.RawMessage, error) {
	data, err := s.client.Delete("/subscribers/" + id)
	if err != nil {
		return nil, fmt.Errorf("deleting subscriber: %w", err)
	}
	return data, nil
}
