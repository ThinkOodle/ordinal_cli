package api

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/ordinal-cli/ordinal/internal/client"
	"github.com/ordinal-cli/ordinal/internal/models"
)

const postsBasePath = "/posts"

// PostService handles post-related API calls.
type PostService struct {
	client *client.Client
}

// NewPostService creates a new PostService.
func NewPostService(c *client.Client) *PostService {
	return &PostService{client: c}
}

// List returns a paginated list of posts.
func (s *PostService) List(params models.ListPostsParams) (*models.PostListResponse, error) {
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
	if params.Status != "" {
		q.Set("status", params.Status)
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
	if params.InstagramProfileID != "" {
		q.Set("instagramProfileId", params.InstagramProfileID)
	}
	if params.LabelIDs != "" {
		q.Set("labelIds", params.LabelIDs)
	}
	if params.PublishDateMin != "" {
		q.Set("publishDateMin", params.PublishDateMin)
	}
	if params.PublishDateMax != "" {
		q.Set("publishDateMax", params.PublishDateMax)
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

	data, err := s.client.Get(postsBasePath, q)
	if err != nil {
		return nil, fmt.Errorf("listing posts: %w", err)
	}

	var resp models.PostListResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing posts list: %w", err)
	}
	return &resp, nil
}

// ListAll fetches all posts by auto-paginating. Fails fast on inconsistent
// cursor metadata (hasMore=true with empty nextCursor, or any cursor the
// server has already handed out in this run) rather than spinning forever
// or silently truncating results. Honors the caller's params.Limit as the
// per-request page size, defaulting to 100 when unset, so --all --limit N
// behaves consistently with what the flag advertises.
func (s *PostService) ListAll(params models.ListPostsParams) ([]models.Post, error) {
	var all []models.Post
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
		all = append(all, resp.Posts...)
		if !resp.HasMore {
			break
		}
		if resp.NextCursor == "" {
			return nil, fmt.Errorf("paginating posts: server reported hasMore=true with empty nextCursor")
		}
		if _, ok := seen[resp.NextCursor]; ok {
			return nil, fmt.Errorf("paginating posts: server returned repeated cursor %q", resp.NextCursor)
		}
		seen[resp.NextCursor] = struct{}{}
		cursor = resp.NextCursor
	}
	return all, nil
}

// Get returns a single post by ID. The API wraps the post in a
// {"post": ...} envelope; unwrap it so callers receive the post itself.
func (s *PostService) Get(id string) (*models.Post, error) {
	data, err := s.client.Get(postsBasePath+"/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("getting post: %w", err)
	}

	var wrapper struct {
		Post models.Post `json:"post"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("parsing post: %w", err)
	}
	return &wrapper.Post, nil
}

// Create creates a new post. Body is a map to accommodate the complex nested
// channel configs (linkedIn, x, instagram). Minimum required fields: title,
// publishAt, status.
func (s *PostService) Create(body map[string]interface{}) (json.RawMessage, error) {
	data, err := s.client.Post(postsBasePath, body)
	if err != nil {
		return nil, fmt.Errorf("creating post: %w", err)
	}
	return data, nil
}

// Update updates a post by ID. All fields are optional.
func (s *PostService) Update(id string, body map[string]interface{}) (json.RawMessage, error) {
	data, err := s.client.Patch(postsBasePath+"/"+id, body)
	if err != nil {
		return nil, fmt.Errorf("updating post: %w", err)
	}
	return data, nil
}

// Archive archives a post by ID.
func (s *PostService) Archive(id string) (json.RawMessage, error) {
	data, err := s.client.Post(postsBasePath+"/"+id+"/archive", nil)
	if err != nil {
		return nil, fmt.Errorf("archiving post: %w", err)
	}
	return data, nil
}

// Unarchive restores a post from the trash.
func (s *PostService) Unarchive(id string) (json.RawMessage, error) {
	data, err := s.client.Post(postsBasePath+"/"+id+"/unarchive", nil)
	if err != nil {
		return nil, fmt.Errorf("unarchiving post: %w", err)
	}
	return data, nil
}

// Schedule schedules or reschedules a post.
func (s *PostService) Schedule(id string, req models.SchedulePostRequest) (json.RawMessage, error) {
	data, err := s.client.Post(postsBasePath+"/"+id+"/schedule", req)
	if err != nil {
		return nil, fmt.Errorf("scheduling post: %w", err)
	}
	return data, nil
}

// Unschedule cancels a scheduled publish for a post.
func (s *PostService) Unschedule(id string) (json.RawMessage, error) {
	data, err := s.client.Post(postsBasePath+"/"+id+"/unschedule", nil)
	if err != nil {
		return nil, fmt.Errorf("unscheduling post: %w", err)
	}
	return data, nil
}

// Delete permanently deletes a post by ID.
func (s *PostService) Delete(id string) (json.RawMessage, error) {
	data, err := s.client.Delete(postsBasePath + "/" + id)
	if err != nil {
		return nil, fmt.Errorf("deleting post: %w", err)
	}
	return data, nil
}
