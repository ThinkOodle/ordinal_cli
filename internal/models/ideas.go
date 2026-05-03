package models

import "encoding/json"

// Idea represents a content idea.
type Idea struct {
	ID            string          `json:"id"`
	URL           string          `json:"url,omitempty"`
	Title         string          `json:"title"`
	Channels      []string        `json:"channels,omitempty"`
	Status        string          `json:"status,omitempty"`
	CreatedAt     string          `json:"createdAt,omitempty"`
	UpdatedAt     string          `json:"updatedAt,omitempty"`
	Labels        []Label         `json:"labels,omitempty"`
	CampaignID    string          `json:"campaignId,omitempty"`
	Notes         string          `json:"notes,omitempty"`
	LinkedIn      json.RawMessage `json:"linkedIn,omitempty"`
	X             json.RawMessage `json:"x,omitempty"`
	TikTok        json.RawMessage `json:"tikTok,omitempty"`
	YouTubeShorts json.RawMessage `json:"youTubeShorts,omitempty"`
}

// ListIdeasParams holds query parameters for listing ideas.
type ListIdeasParams struct {
	Limit             int
	Cursor            string
	IDs               string
	Channel           string
	LinkedInProfileID string
	XProfileID        string
	TikTokProfileID   string
	YouTubeProfileID  string
	LabelIDs          string
	CreatedAtMin      string
	CreatedAtMax      string
	SortBy            string
	SortOrder         string
}

// IdeaListResponse is the response for listing ideas.
type IdeaListResponse struct {
	Ideas      []Idea `json:"ideas"`
	NextCursor string `json:"nextCursor,omitempty"`
	HasMore    bool   `json:"hasMore"`
}

// AddIdeaToCalendarRequest is the request body for converting an idea to a post.
type AddIdeaToCalendarRequest struct {
	PublishDate string `json:"publishDate"`
}
