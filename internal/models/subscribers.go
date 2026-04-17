package models

// SubscriberUser is a user referenced in a subscriber.
type SubscriberUser struct {
	ID              string `json:"id,omitempty"`
	Email           string `json:"email,omitempty"`
	FirstName       string `json:"firstName,omitempty"`
	LastName        string `json:"lastName,omitempty"`
	ProfileImageURL string `json:"profileImageUrl,omitempty"`
}

// Subscriber represents a post subscriber.
type Subscriber struct {
	ID        string          `json:"id"`
	User      SubscriberUser  `json:"user"`
	CreatedBy *SubscriberUser `json:"createdBy,omitempty"`
	CreatedAt string          `json:"createdAt,omitempty"`
}

// CreateSubscribersRequest is the request body for adding subscribers to a post.
type CreateSubscribersRequest struct {
	PostID  string   `json:"postId"`
	UserIDs []string `json:"userIds"`
}

// SubscriberListResponse is the response for listing post subscribers.
type SubscriberListResponse struct {
	Subscribers []Subscriber `json:"subscribers"`
}
