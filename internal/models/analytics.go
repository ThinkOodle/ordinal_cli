package models

// FollowerDataPoint is a single follower count observation.
type FollowerDataPoint struct {
	FollowerCount int    `json:"followerCount"`
	RecordedAt    string `json:"recordedAt"`
}

// CpmValues is the workspace CPM configuration for EMV calculation.
type CpmValues struct {
	LinkedIn  float64 `json:"linkedIn"`
	X         float64 `json:"x"`
	Instagram float64 `json:"instagram"`
	Facebook  float64 `json:"facebook"`
	Threads   float64 `json:"threads"`
}

// CpmUpdateRequest is the partial update for CPM values.
type CpmUpdateRequest struct {
	LinkedIn  *float64 `json:"linkedIn,omitempty"`
	X         *float64 `json:"x,omitempty"`
	Instagram *float64 `json:"instagram,omitempty"`
	Facebook  *float64 `json:"facebook,omitempty"`
	Threads   *float64 `json:"threads,omitempty"`
}

// AnalyticsDateRange holds date-range filters used by analytics endpoints.
type AnalyticsDateRange struct {
	StartDate string
	EndDate   string
}
