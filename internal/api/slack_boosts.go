package api

import (
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

const slackBoostsBasePath = "/slack-boosts"
const slackWebhooksBasePath = "/slack-webhooks"

// SlackBoostService handles Slack boost API calls.
type SlackBoostService struct {
	client *client.Client
}

// NewSlackBoostService creates a new SlackBoostService.
func NewSlackBoostService(c *client.Client) *SlackBoostService {
	return &SlackBoostService{client: c}
}

// ListByPost returns all Slack boosts attached to a post.
func (s *SlackBoostService) ListByPost(postID string) (*models.SlackBoostListResponse, error) {
	data, err := s.client.Get("/posts/"+postID+"/slack-boosts", nil)
	if err != nil {
		return nil, fmt.Errorf("listing slack boosts: %w", err)
	}

	var resp models.SlackBoostListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing slack boosts list: %w", err)
	}
	return &resp, nil
}

// slackBoostEnvelope matches the {"slackBoost": ...} wrapper the API
// returns on get, create, and update. Unwrapping it here keeps callers
// from silently receiving a zero-valued SlackBoost on success.
type slackBoostEnvelope struct {
	SlackBoost models.SlackBoost `json:"slackBoost"`
}

// Get returns a Slack boost by ID.
func (s *SlackBoostService) Get(id string) (*models.SlackBoost, error) {
	data, err := s.client.Get(slackBoostsBasePath+"/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("getting slack boost: %w", err)
	}

	var env slackBoostEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parsing slack boost: %w", err)
	}
	return &env.SlackBoost, nil
}

// Create creates a Slack boost.
func (s *SlackBoostService) Create(req models.CreateSlackBoostRequest) (*models.SlackBoost, error) {
	data, err := s.client.Post(slackBoostsBasePath, req)
	if err != nil {
		return nil, fmt.Errorf("creating slack boost: %w", err)
	}

	var env slackBoostEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parsing created slack boost: %w", err)
	}
	return &env.SlackBoost, nil
}

// Update updates a Slack boost.
func (s *SlackBoostService) Update(id string, req models.UpdateSlackBoostRequest) (*models.SlackBoost, error) {
	data, err := s.client.Patch(slackBoostsBasePath+"/"+id, req)
	if err != nil {
		return nil, fmt.Errorf("updating slack boost: %w", err)
	}

	var env slackBoostEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("parsing updated slack boost: %w", err)
	}
	return &env.SlackBoost, nil
}

// Delete deletes a Slack boost by ID.
func (s *SlackBoostService) Delete(id string) error {
	if _, err := s.client.Delete(slackBoostsBasePath + "/" + id); err != nil {
		return fmt.Errorf("deleting slack boost: %w", err)
	}
	return nil
}

// SlackWebhookService handles Slack webhook (boost channel) API calls.
type SlackWebhookService struct {
	client *client.Client
}

// NewSlackWebhookService creates a new SlackWebhookService.
func NewSlackWebhookService(c *client.Client) *SlackWebhookService {
	return &SlackWebhookService{client: c}
}

// List returns all connected Slack boost channels in the workspace.
func (s *SlackWebhookService) List() (*models.SlackWebhookListResponse, error) {
	data, err := s.client.Get(slackWebhooksBasePath, nil)
	if err != nil {
		return nil, fmt.Errorf("listing slack webhooks: %w", err)
	}

	var resp models.SlackWebhookListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing slack webhooks list: %w", err)
	}
	return &resp, nil
}
