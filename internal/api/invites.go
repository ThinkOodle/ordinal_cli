package api

import (
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

const invitesBasePath = "/invites"

// InviteService handles invite-related API calls.
type InviteService struct {
	client *client.Client
}

// NewInviteService creates a new InviteService.
func NewInviteService(c *client.Client) *InviteService {
	return &InviteService{client: c}
}

// List returns the pending invites.
func (s *InviteService) List() ([]models.Invite, error) {
	data, err := s.client.Get(invitesBasePath, nil)
	if err != nil {
		return nil, fmt.Errorf("listing invites: %w", err)
	}

	var invites []models.Invite
	if err := json.Unmarshal(data, &invites); err != nil {
		return nil, fmt.Errorf("parsing invites: %w", err)
	}
	return invites, nil
}

// Create creates a new invite. The response envelope carries either a new
// Invite (when an email was sent) or a null invite with sentEmail=false
// (when the user already existed and was added to the workspace directly).
func (s *InviteService) Create(req models.CreateInviteRequest) (*models.CreateInviteResponse, error) {
	data, err := s.client.Post(invitesBasePath, req)
	if err != nil {
		return nil, fmt.Errorf("creating invite: %w", err)
	}

	var resp models.CreateInviteResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing created invite: %w", err)
	}
	return &resp, nil
}

// Delete deletes an invite by ID.
func (s *InviteService) Delete(id string) error {
	if _, err := s.client.Delete(invitesBasePath + "/" + id); err != nil {
		return fmt.Errorf("deleting invite: %w", err)
	}
	return nil
}
