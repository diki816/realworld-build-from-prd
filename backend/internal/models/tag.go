package models

// Tag represents a tag in the system
type Tag struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// TagsResponse represents the response format for tags
type TagsResponse struct {
	Tags []string `json:"tags"`
}