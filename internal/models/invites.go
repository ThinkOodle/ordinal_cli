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

// CreateInviteResponse is the envelope returned by POST /invites. Invite is
// nil when the invited user already had an Ordinal account and was added to
// the workspace directly (in which case SentEmail is false).
type CreateInviteResponse struct {
	Invite    *Invite `json:"invite"`
	SentEmail bool    `json:"sentEmail"`
}
