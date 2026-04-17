package models

// Webhook represents a webhook subscription.
type Webhook struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Description string            `json:"description,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Topics      []string          `json:"topics"`
	CreatedAt   string            `json:"createdAt,omitempty"`
	CreatedBy   *UserSummary      `json:"createdBy,omitempty"`
}

// CreateWebhookRequest is the request body for creating a webhook.
type CreateWebhookRequest struct {
	Name        string            `json:"name"`
	URL         string            `json:"url"`
	Description string            `json:"description,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Topics      []string          `json:"topics"`
}

// UpdateWebhookRequest is the request body for updating a webhook.
// All fields are optional — merge arbitrary updates via a map[string]interface{}.
type UpdateWebhookRequest = map[string]interface{}
