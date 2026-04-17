package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

const ideasBasePath = "/ideas"

// IdeaService handles idea-related API calls.
type IdeaService struct {
	client *client.Client
}

// NewIdeaService creates a new IdeaService.
func NewIdeaService(c *client.Client) *IdeaService {
	return &IdeaService{client: c}
}

// List returns a paginated list of ideas.
func (s *IdeaService) List(params models.ListIdeasParams) (*models.IdeaListResponse, error) {
	q := url.Values{}
	if params.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", params.Limit))
	}
	if params.Cursor != "" {
		q.Set("cursor", params.Cursor)
	}
	if params.IDs != "" {
		q.Set("ids", params.IDs)
	}
	if params.Channel != "" {
		q.Set("channel", params.Channel)
	}
	if params.LinkedInProfileID != "" {
		q.Set("linkedInProfileId", params.LinkedInProfileID)
	}
	if params.XProfileID != "" {
		q.Set("xProfileId", params.XProfileID)
	}
	if params.LabelIDs != "" {
		q.Set("labelIds", params.LabelIDs)
	}
	if params.CreatedAtMin != "" {
		q.Set("createdAtMin", params.CreatedAtMin)
	}
	if params.CreatedAtMax != "" {
		q.Set("createdAtMax", params.CreatedAtMax)
	}
	if params.SortBy != "" {
		q.Set("sortBy", params.SortBy)
	}
	if params.SortOrder != "" {
		q.Set("sortOrder", params.SortOrder)
	}

	data, err := s.client.Get(ideasBasePath, q)
	if err != nil {
		return nil, fmt.Errorf("listing ideas: %w", err)
	}

	var resp models.IdeaListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing ideas list: %w", err)
	}
	return &resp, nil
}

// ListAll fetches all ideas by auto-paginating. Fails fast on inconsistent
// cursor metadata (hasMore=true with empty nextCursor, or any cursor the
// server has already handed out in this run) rather than spinning forever
// or silently truncating results. Honors the caller's params.Limit as the
// per-request page size, defaulting to 100 when unset, so --all --limit N
// behaves consistently with what the flag advertises.
func (s *IdeaService) ListAll(params models.ListIdeasParams) ([]models.Idea, error) {
	var all []models.Idea
	cursor := params.Cursor
	seen := make(map[string]struct{})
	if params.Limit <= 0 {
		params.Limit = 100
	}

	for {
		params.Cursor = cursor

		resp, err := s.List(params)
		if err != nil {
			return nil, err
		}
		all = append(all, resp.Ideas...)
		if !resp.HasMore {
			break
		}
		if resp.NextCursor == "" {
			return nil, fmt.Errorf("paginating ideas: server reported hasMore=true with empty nextCursor")
		}
		if _, ok := seen[resp.NextCursor]; ok {
			return nil, fmt.Errorf("paginating ideas: server returned repeated cursor %q", resp.NextCursor)
		}
		seen[resp.NextCursor] = struct{}{}
		cursor = resp.NextCursor
	}
	return all, nil
}

// Get returns a single idea by ID. The API wraps the idea in an
// {"idea": ...} envelope; unwrap it so callers receive the idea itself.
func (s *IdeaService) Get(id string) (*models.Idea, error) {
	data, err := s.client.Get(ideasBasePath+"/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("getting idea: %w", err)
	}

	var wrapper struct {
		Idea models.Idea `json:"idea"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing idea: %w", err)
	}
	return &wrapper.Idea, nil
}

// Create creates a new idea.
func (s *IdeaService) Create(body map[string]interface{}) (json.RawMessage, error) {
	data, err := s.client.Post(ideasBasePath, body)
	if err != nil {
		return nil, fmt.Errorf("creating idea: %w", err)
	}
	return data, nil
}

// Update updates an idea by ID.
func (s *IdeaService) Update(id string, body map[string]interface{}) (json.RawMessage, error) {
	data, err := s.client.Patch(ideasBasePath+"/"+id, body)
	if err != nil {
		return nil, fmt.Errorf("updating idea: %w", err)
	}
	return data, nil
}

// Archive archives an idea.
func (s *IdeaService) Archive(id string) (json.RawMessage, error) {
	data, err := s.client.Post(ideasBasePath+"/"+id+"/archive", nil)
	if err != nil {
		return nil, fmt.Errorf("archiving idea: %w", err)
	}
	return data, nil
}

// Unarchive restores an archived idea.
func (s *IdeaService) Unarchive(id string) (json.RawMessage, error) {
	data, err := s.client.Post(ideasBasePath+"/"+id+"/unarchive", nil)
	if err != nil {
		return nil, fmt.Errorf("unarchiving idea: %w", err)
	}
	return data, nil
}

// AddToCalendar converts an idea to a scheduled post.
func (s *IdeaService) AddToCalendar(id string, req models.AddIdeaToCalendarRequest) (json.RawMessage, error) {
	data, err := s.client.Post(ideasBasePath+"/"+id+"/add-to-calendar", req)
	if err != nil {
		return nil, fmt.Errorf("adding idea to calendar: %w", err)
	}
	return data, nil
}
