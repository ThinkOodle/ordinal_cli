package api

import (
	"encoding/json"
	"fmt"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

// CommentService handles comment-related API calls.
type CommentService struct {
	client *client.Client
}

// NewCommentService creates a new CommentService.
func NewCommentService(c *client.Client) *CommentService {
	return &CommentService{client: c}
}

// List returns comments for a post.
func (s *CommentService) List(postID string) (*models.CommentListResponse, error) {
	data, err := s.client.Get("/posts/"+postID+"/comments", nil)
	if err != nil {
		return nil, fmt.Errorf("listing comments: %w", err)
	}

	var resp models.CommentListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing comments list: %w", err)
	}
	return &resp, nil
}

// Create creates a comment on a post.
func (s *CommentService) Create(postID string, req models.CreateCommentRequest) (*models.Comment, error) {
	data, err := s.client.Post("/posts/"+postID+"/comments", req)
	if err != nil {
		return nil, fmt.Errorf("creating comment: %w", err)
	}

	var c models.Comment
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("parsing created comment: %w", err)
	}
	return &c, nil
}

// Delete deletes a comment by ID. Only the author can delete their own
// comment. The API returns the deleted Comment as the response body.
func (s *CommentService) Delete(commentID string) (json.RawMessage, error) {
	data, err := s.client.Delete("/comments/" + commentID)
	if err != nil {
		return nil, fmt.Errorf("deleting comment: %w", err)
	}
	return data, nil
}
