package api

import (
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

const workspacesBasePath = "/workspace"

// WorkspaceService handles workspace-related API calls.
type WorkspaceService struct {
	client *client.Client
}

// NewWorkspaceService creates a new WorkspaceService.
func NewWorkspaceService(c *client.Client) *WorkspaceService {
	return &WorkspaceService{client: c}
}

// Get returns the current workspace.
func (s *WorkspaceService) Get() (*models.Workspace, error) {
	data, err := s.client.Get(workspacesBasePath, nil)
	if err != nil {
		return nil, fmt.Errorf("getting workspace: %w", err)
	}

	var ws models.Workspace
	if err := json.Unmarshal(data, &ws); err != nil {
		return nil, fmt.Errorf("parsing workspace: %w", err)
	}
	return &ws, nil
}
