package models

// Workspace represents the current workspace.
type Workspace struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Slug      string `json:"slug,omitempty"`
	Timezone  string `json:"timezone,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
}
