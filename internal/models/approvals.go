package models

// ApprovalUser is a user referenced in an approval.
type ApprovalUser struct {
	ID              string `json:"id,omitempty"`
	Email           string `json:"email,omitempty"`
	FirstName       string `json:"firstName,omitempty"`
	LastName        string `json:"lastName,omitempty"`
	ProfileImageURL string `json:"profileImageUrl,omitempty"`
}

// Approval represents a post approval request.
type Approval struct {
	ID          string        `json:"id"`
	Status      string        `json:"status"`
	IsBlocking  bool          `json:"isBlocking"`
	Message     string        `json:"message,omitempty"`
	DueDate     string        `json:"dueDate,omitempty"`
	ApprovedAt  string        `json:"approvedAt,omitempty"`
	CreatedAt   string        `json:"createdAt,omitempty"`
	User        *ApprovalUser `json:"user,omitempty"`
	RequestedBy *ApprovalUser `json:"requestedBy,omitempty"`
}

// ApprovalRequestInput is a single approval request for a post.
//
// IsBlocking is a *bool rather than bool so the CLI can distinguish three
// cases the API treats differently: the caller explicitly requested a
// blocking approval (true), explicitly requested a non-blocking approval
// (false), or left the field off entirely and expects server defaulting.
// A plain bool with omitempty collapses the explicit-false case into the
// "omitted" case on the wire.
type ApprovalRequestInput struct {
	UserID     string `json:"userId"`
	Message    string `json:"message,omitempty"`
	DueDate    string `json:"dueDate,omitempty"`
	IsBlocking *bool  `json:"isBlocking,omitempty"`
}

// CreateApprovalsRequest is the request body for creating post approvals.
type CreateApprovalsRequest struct {
	PostID    string                 `json:"postId"`
	Approvals []ApprovalRequestInput `json:"approvals"`
}

// ApprovalListResponse is the response for listing post approvals.
type ApprovalListResponse struct {
	Approvals []Approval `json:"approvals"`
}
