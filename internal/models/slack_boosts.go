package models

// SlackWebhookSummary is a nested reference to a Slack webhook.
type SlackWebhookSummary struct {
	ID          string `json:"id,omitempty"`
	ChannelName string `json:"channelName,omitempty"`
	TeamName    string `json:"teamName,omitempty"`
}

// SlackBoost represents a Slack boost attached to a post.
type SlackBoost struct {
	ID            string              `json:"id"`
	Copy          string              `json:"copy,omitempty"`
	IsAutoCreated bool                `json:"isAutoCreated"`
	PostID        string              `json:"postId,omitempty"`
	SlackWebhook  SlackWebhookSummary `json:"slackWebhook,omitempty"`
	CreatedAt     string              `json:"createdAt,omitempty"`
	UpdatedAt     string              `json:"updatedAt,omitempty"`
}

// CreateSlackBoostRequest is the request body for creating a Slack boost.
type CreateSlackBoostRequest struct {
	PostID         string `json:"postId"`
	SlackWebhookID string `json:"slackWebhookId"`
	Copy           string `json:"copy,omitempty"`
}

// UpdateSlackBoostRequest is the request body for updating a Slack boost.
type UpdateSlackBoostRequest = map[string]interface{}

// SlackBoostListResponse is the response for listing post Slack boosts.
type SlackBoostListResponse struct {
	SlackBoosts []SlackBoost `json:"slackBoosts"`
}

// SlackWebhook represents a connected Slack marketing boost channel.
type SlackWebhook struct {
	ID              string `json:"id"`
	ChannelName     string `json:"channelName,omitempty"`
	TeamName        string `json:"teamName,omitempty"`
	TeamIconURL     string `json:"teamIconUrl,omitempty"`
	NotifyByDefault bool   `json:"notifyByDefault"`
	CreatedAt       string `json:"createdAt,omitempty"`
	UpdatedAt       string `json:"updatedAt,omitempty"`
}

// SlackWebhookListResponse is the response for listing Slack webhooks.
type SlackWebhookListResponse struct {
	SlackWebhooks []SlackWebhook `json:"slackWebhooks"`
}
