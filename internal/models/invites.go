package models

// Invite represents a pending workspace invite.
type Invite struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// CreateInviteRequest is the request body for creating an invite.
type CreateInviteRequest struct {
	Email string `json:"email"`
}
