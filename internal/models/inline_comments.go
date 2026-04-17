package models

import "encoding/json"

// InlineComment represents a text-anchored inline comment thread.
type InlineComment struct {
	ID              string            `json:"id"`
	Resolved        bool              `json:"resolved"`
	Channel         string            `json:"channel"`
	HighlightedText string            `json:"highlightedText,omitempty"`
	Replies         []json.RawMessage `json:"replies,omitempty"`
	CreatedAt       string            `json:"createdAt,omitempty"`
	UpdatedAt       string            `json:"updatedAt,omitempty"`
}

// ListInlineCommentsParams holds query parameters for listing inline comments.
type ListInlineCommentsParams struct {
	Channel  string
	Resolved *bool
}

// InlineCommentListResponse is the response for listing inline comments.
type InlineCommentListResponse struct {
	InlineComments []InlineComment `json:"inlineComments"`
}
