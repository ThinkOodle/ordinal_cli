package api

import (
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

const uploadsBasePath = "/uploads"

// UploadService handles file upload API calls.
type UploadService struct {
	client *client.Client
}

// NewUploadService creates a new UploadService.
func NewUploadService(c *client.Client) *UploadService {
	return &UploadService{client: c}
}

// Create starts a new upload from a URL.
func (s *UploadService) Create(req models.CreateUploadRequest) (json.RawMessage, error) {
	data, err := s.client.Post(uploadsBasePath, req)
	if err != nil {
		return nil, fmt.Errorf("creating upload: %w", err)
	}
	return data, nil
}

// Get returns the status of an upload job.
func (s *UploadService) Get(id string) (json.RawMessage, error) {
	data, err := s.client.Get(uploadsBasePath+"/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("getting upload: %w", err)
	}
	return data, nil
}
