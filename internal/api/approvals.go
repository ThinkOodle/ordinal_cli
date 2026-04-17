package api

import (
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

// ApprovalService handles approval-related API calls.
type ApprovalService struct {
	client *client.Client
}

// NewApprovalService creates a new ApprovalService.
func NewApprovalService(c *client.Client) *ApprovalService {
	return &ApprovalService{client: c}
}

// List returns all approvals for a post.
func (s *ApprovalService) List(postID string) (*models.ApprovalListResponse, error) {
	data, err := s.client.Get("/posts/"+postID+"/approvals", nil)
	if err != nil {
		return nil, fmt.Errorf("listing approvals: %w", err)
	}

	var resp models.ApprovalListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing approvals list: %w", err)
	}
	return &resp, nil
}

// Create creates one or more approval requests for a post.
func (s *ApprovalService) Create(req models.CreateApprovalsRequest) (json.RawMessage, error) {
	data, err := s.client.Post("/approvals", req)
	if err != nil {
		return nil, fmt.Errorf("creating approvals: %w", err)
	}
	return data, nil
}

// Delete deletes an approval by ID. The API returns a
// `{"deletedApproval": Approval}` envelope; we forward it verbatim.
func (s *ApprovalService) Delete(id string) (json.RawMessage, error) {
	data, err := s.client.Delete("/approvals/" + id)
	if err != nil {
		return nil, fmt.Errorf("deleting approval: %w", err)
	}
	return data, nil
}
