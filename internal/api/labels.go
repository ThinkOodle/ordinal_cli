// Package api provides typed methods for each Ordinal API resource.
package api

import (
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

const labelsBasePath = "/labels"

// LabelService handles label-related API calls.
type LabelService struct {
	client *client.Client
}

// NewLabelService creates a new LabelService.
func NewLabelService(c *client.Client) *LabelService {
	return &LabelService{client: c}
}

// List returns all labels in the workspace.
func (s *LabelService) List() (*models.LabelListResponse, error) {
	data, err := s.client.Get(labelsBasePath, nil)
	if err != nil {
		return nil, fmt.Errorf("listing labels: %w", err)
	}

	var resp models.LabelListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing labels list: %w", err)
	}
	return &resp, nil
}

// Create creates a new label.
func (s *LabelService) Create(req models.CreateLabelRequest) (*models.Label, error) {
	data, err := s.client.Post(labelsBasePath, req)
	if err != nil {
		return nil, fmt.Errorf("creating label: %w", err)
	}

	var label models.Label
	if err := json.Unmarshal(data, &label); err != nil {
		return nil, fmt.Errorf("parsing created label: %w", err)
	}
	return &label, nil
}

// Delete deletes a label by ID.
func (s *LabelService) Delete(id string) error {
	if _, err := s.client.Delete(labelsBasePath + "/" + id); err != nil {
		return fmt.Errorf("deleting label: %w", err)
	}
	return nil
}
