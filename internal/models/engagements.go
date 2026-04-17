package models

import "encoding/json"

// Engagement represents an auto-engagement on a post.
type Engagement struct {
	ID           string          `json:"id"`
	Channel      string          `json:"channel"`
	Type         string          `json:"type"`
	ProfileID    string          `json:"profileId,omitempty"`
	Copy         string          `json:"copy,omitempty"`
	DelaySeconds int             `json:"delaySeconds,omitempty"`
	ReactionType json.RawMessage `json:"reactionType,omitempty"`
}

// EngagementInput is a single engagement input for create requests.
// Shape varies by Type; callers provide this as raw JSON via --body-json or
// --body-file for create.
type EngagementInput = map[string]interface{}

// CreateEngagementsRequest is the request body for creating engagements.
type CreateEngagementsRequest struct {
	Channel     string            `json:"channel"`
	Engagements []EngagementInput `json:"engagements"`
}

// EngagementListResponse is the response for listing engagements.
type EngagementListResponse struct {
	Engagements []Engagement `json:"engagements"`
}

// UpdateEngagementRequest is the request body for updating a single engagement.
// Only include fields to change.
type UpdateEngagementRequest = map[string]interface{}
