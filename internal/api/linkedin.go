package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ordinal-cli/ordinal/internal/client"
)

// LinkedInService handles LinkedIn profile/mention lookups.
type LinkedInService struct {
	client *client.Client
}

// NewLinkedInService creates a new LinkedInService.
func NewLinkedInService(c *client.Client) *LinkedInService {
	return &LinkedInService{client: c}
}

// GetProfile looks up a LinkedIn profile by URN.
func (s *LinkedInService) GetProfile(urn string) (json.RawMessage, error) {
	data, err := s.client.Get("/linkedin/profile/"+url.PathEscape(urn), nil)
	if err != nil {
		return nil, fmt.Errorf("getting linkedin profile: %w", err)
	}
	return data, nil
}

// GetMention returns the mention format for a LinkedIn username.
func (s *LinkedInService) GetMention(username string) (json.RawMessage, error) {
	data, err := s.client.Get("/linkedin/"+url.PathEscape(username)+"/mentions", nil)
	if err != nil {
		return nil, fmt.Errorf("getting linkedin mention: %w", err)
	}
	return data, nil
}

// LinkedInLeadsService handles LinkedIn leads endpoints.
type LinkedInLeadsService struct {
	client *client.Client
}

// NewLinkedInLeadsService creates a new LinkedInLeadsService.
func NewLinkedInLeadsService(c *client.Client) *LinkedInLeadsService {
	return &LinkedInLeadsService{client: c}
}

// ListPostsParams holds query parameters for listing LinkedIn lead posts.
type LinkedInLeadsListPostsParams struct {
	StartDate string
	EndDate   string
	Limit     int
	Cursor    string
}

// ListPosts lists LinkedIn posts for a profile's leads scraping.
func (s *LinkedInLeadsService) ListPosts(profileID string, p LinkedInLeadsListPostsParams) (json.RawMessage, error) {
	q := url.Values{}
	if p.StartDate != "" {
		q.Set("startDate", p.StartDate)
	}
	if p.EndDate != "" {
		q.Set("endDate", p.EndDate)
	}
	if p.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}

	data, err := s.client.Get("/linkedin/leads/"+profileID+"/posts", q)
	if err != nil {
		return nil, fmt.Errorf("listing linkedin lead posts: %w", err)
	}
	return data, nil
}

// GetLeadsParams holds query parameters for listing leads on a LinkedIn post.
type LinkedInLeadsGetLeadsParams struct {
	Types            string
	MinFollowerCount int
	Limit            int
	Cursor           string
}

// GetLeadsByPost returns leads (engagers) for a LinkedIn post.
func (s *LinkedInLeadsService) GetLeadsByPost(profileID, postID string, p LinkedInLeadsGetLeadsParams) (json.RawMessage, error) {
	q := url.Values{}
	if p.Types != "" {
		q.Set("types", p.Types)
	}
	if p.MinFollowerCount > 0 {
		q.Set("minFollowerCount", fmt.Sprintf("%d", p.MinFollowerCount))
	}
	if p.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", p.Limit))
	}
	if p.Cursor != "" {
		q.Set("cursor", p.Cursor)
	}

	data, err := s.client.Get("/linkedin/leads/"+profileID+"/posts/"+postID, q)
	if err != nil {
		return nil, fmt.Errorf("getting linkedin post leads: %w", err)
	}
	return data, nil
}
