package models

// UserSummary is a minimal user representation used in nested fields.
type UserSummary struct {
	ID        string `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
}

// Profile represents a social profile (engagement or scheduling).
type Profile struct {
	ID                      string       `json:"id"`
	Name                    string       `json:"name,omitempty"`
	Detail                  string       `json:"detail,omitempty"`
	Channel                 string       `json:"channel"`
	ProfileImageURL         string       `json:"profileImageUrl,omitempty"`
	IsReintegrationRequired bool         `json:"isReintegrationRequired"`
	CreatedAt               string       `json:"createdAt,omitempty"`
	CreatedBy               *UserSummary `json:"createdBy,omitempty"`
}

// SchedulingProfile extends Profile with scheduling-only fields.
type SchedulingProfile struct {
	Profile
	IsLeadsScrapingEnabled bool `json:"isLeadsScrapingEnabled"`
}
