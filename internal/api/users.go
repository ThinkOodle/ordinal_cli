package api

import (
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

const usersBasePath = "/users"

// UserService handles user-related API calls.
type UserService struct {
	client *client.Client
}

// NewUserService creates a new UserService.
func NewUserService(c *client.Client) *UserService {
	return &UserService{client: c}
}

// List returns the workspace users.
func (s *UserService) List() ([]models.User, error) {
	data, err := s.client.Get(usersBasePath, nil)
	if err != nil {
		return nil, fmt.Errorf("listing users: %w", err)
	}

	var users []models.User
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("parsing users: %w", err)
	}
	return users, nil
}
