package models

// CommentUser is the commenter's user info.
type CommentUser struct {
	ID              string `json:"id,omitempty"`
	Email           string `json:"email,omitempty"`
	FirstName       string `json:"firstName,omitempty"`
	LastName        string `json:"lastName,omitempty"`
	ProfileImageURL string `json:"profileImageUrl,omitempty"`
}

// Comment represents a comment on a post.
type Comment struct {
	ID        string       `json:"id"`
	Message   string       `json:"message"`
	PostID    string       `json:"postId,omitempty"`
	CreatedAt string       `json:"createdAt,omitempty"`
	UpdatedAt string       `json:"updatedAt,omitempty"`
	User      *CommentUser `json:"user,omitempty"`
}

// CreateCommentRequest is the request body for creating a comment.
type CreateCommentRequest struct {
	Message string `json:"message"`
}

// CommentListResponse is the response for listing comments.
type CommentListResponse struct {
	Comments []Comment `json:"comments"`
}
