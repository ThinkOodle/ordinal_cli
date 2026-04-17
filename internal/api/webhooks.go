package api

import (
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

const webhooksBasePath = "/webhooks"

// WebhookService handles webhook-related API calls.
type WebhookService struct {
	client *client.Client
}

// NewWebhookService creates a new WebhookService.
func NewWebhookService(c *client.Client) *WebhookService {
	return &WebhookService{client: c}
}

// List returns all webhooks.
func (s *WebhookService) List() ([]models.Webhook, error) {
	data, err := s.client.Get(webhooksBasePath, nil)
	if err != nil {
		return nil, fmt.Errorf("listing webhooks: %w", err)
	}

	var webhooks []models.Webhook
	if err := json.Unmarshal(data, &webhooks); err != nil {
		return nil, fmt.Errorf("parsing webhooks: %w", err)
	}
	return webhooks, nil
}

// Get returns a single webhook by ID.
func (s *WebhookService) Get(id string) (*models.Webhook, error) {
	data, err := s.client.Get(webhooksBasePath+"/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("getting webhook: %w", err)
	}

	var w models.Webhook
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, fmt.Errorf("parsing webhook: %w", err)
	}
	return &w, nil
}

// Create creates a new webhook.
func (s *WebhookService) Create(req models.CreateWebhookRequest) (*models.Webhook, error) {
	data, err := s.client.Post(webhooksBasePath, req)
	if err != nil {
		return nil, fmt.Errorf("creating webhook: %w", err)
	}

	var w models.Webhook
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, fmt.Errorf("parsing created webhook: %w", err)
	}
	return &w, nil
}

// Update updates an existing webhook.
func (s *WebhookService) Update(id string, req models.UpdateWebhookRequest) (*models.Webhook, error) {
	data, err := s.client.Patch(webhooksBasePath+"/"+id, req)
	if err != nil {
		return nil, fmt.Errorf("updating webhook: %w", err)
	}

	var w models.Webhook
	if err := json.Unmarshal(data, &w); err != nil {
		return nil, fmt.Errorf("parsing updated webhook: %w", err)
	}
	return &w, nil
}

// Delete deletes a webhook by ID. The API returns `{"success": true}`;
// forwarding it keeps CLI output honest about the real API response.
func (s *WebhookService) Delete(id string) (json.RawMessage, error) {
	data, err := s.client.Delete(webhooksBasePath + "/" + id)
	if err != nil {
		return nil, fmt.Errorf("deleting webhook: %w", err)
	}
	return data, nil
}
