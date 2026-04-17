// Package models defines request and response types for the Ordinal API.
package models

// Label represents a label that can be attached to posts or ideas.
type Label struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Color           string `json:"color"`
	BackgroundColor string `json:"backgroundColor,omitempty"`
	CreatedAt       string `json:"createdAt,omitempty"`
	UpdatedAt       string `json:"updatedAt,omitempty"`
}

// CreateLabelRequest is the request body for creating a label.
type CreateLabelRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// LabelListResponse is the response for listing labels.
type LabelListResponse struct {
	Labels []Label `json:"labels"`
}
