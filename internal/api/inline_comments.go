package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

// InlineCommentService handles inline-comment API calls.
type InlineCommentService struct {
	client *client.Client
}

// NewInlineCommentService creates a new InlineCommentService.
func NewInlineCommentService(c *client.Client) *InlineCommentService {
	return &InlineCommentService{client: c}
}

// List returns the inline comments on a post.
func (s *InlineCommentService) List(postID string, params models.ListInlineCommentsParams) (*models.InlineCommentListResponse, error) {
	q := url.Values{}
	if params.Channel != "" {
		q.Set("channel", params.Channel)
	}
	if params.Resolved != nil {
		q.Set("resolved", strconv.FormatBool(*params.Resolved))
	}

	data, err := s.client.Get("/posts/"+postID+"/inline-comments", q)
	if err != nil {
		return nil, fmt.Errorf("listing inline comments: %w", err)
	}

	var resp models.InlineCommentListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing inline comments list: %w", err)
	}
	return &resp, nil
}
