package models

import "encoding/json"

// CreateUploadRequest is the request body for creating an upload job.
type CreateUploadRequest struct {
	URL string `json:"url"`
}

// Upload is the response for an upload job. The underlying API returns a
// discriminated oneOf by status ("pending", "processing", "ready", "failed",
// "expired"), so we keep the raw payload alongside the common fields.
type Upload struct {
	ID     string          `json:"id"`
	Status string          `json:"status"`
	Asset  json.RawMessage `json:"asset,omitempty"`
	Error  json.RawMessage `json:"error,omitempty"`
	Raw    json.RawMessage `json:"-"`
}
