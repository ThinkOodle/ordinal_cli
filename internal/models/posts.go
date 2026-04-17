package models

import "encoding/json"

// Post represents an Ordinal post. Channel-specific fields (linkedIn, x,
// instagram) are deeply nested and kept as raw JSON to preserve full fidelity.
type Post struct {
	ID          string          `json:"id"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title"`
	Channels    []string        `json:"channels,omitempty"`
	Status      string          `json:"status"`
	PublishDate string          `json:"publishDate,omitempty"`
	PublishAt   string          `json:"publishAt,omitempty"`
	CreatedAt   string          `json:"createdAt,omitempty"`
	UpdatedAt   string          `json:"updatedAt,omitempty"`
	ArchivedAt  string          `json:"archivedAt,omitempty"`
	Labels      []Label         `json:"labels,omitempty"`
	CampaignID  string          `json:"campaignId,omitempty"`
	Notes       string          `json:"notes,omitempty"`
	LinkedIn    json.RawMessage `json:"linkedIn,omitempty"`
	X           json.RawMessage `json:"x,omitempty"`
	Instagram   json.RawMessage `json:"instagram,omitempty"`
}

// ListPostsParams holds query parameters for listing posts.
type ListPostsParams struct {
	Limit              int
	Cursor             string
	IDs                string
	Status             string
	Channel            string
	LinkedInProfileID  string
	XProfileID         string
	InstagramProfileID string
	LabelIDs           string
	PublishDateMin     string
	PublishDateMax     string
	CreatedAtMin       string
	CreatedAtMax       string
	SortBy             string
	SortOrder          string
}

// PostListResponse is the response for listing posts.
type PostListResponse struct {
	Posts      []Post `json:"posts"`
	NextCursor string `json:"nextCursor,omitempty"`
	HasMore    bool   `json:"hasMore"`
}

// SchedulePostRequest is the request body for scheduling a post.
type SchedulePostRequest struct {
	PublishAt string `json:"publishAt,omitempty"`
}
