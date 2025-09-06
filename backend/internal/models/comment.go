package models

import (
	"errors"
	"time"
)

// Comment represents a comment in the system
type Comment struct {
	ID        int       `json:"id" db:"id"`
	Body      string    `json:"body" db:"body"`
	AuthorID  int       `json:"-" db:"author_id"`
	ArticleID int       `json:"-" db:"article_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	Author    Profile   `json:"author"`
}

// CreateCommentRequest represents the request payload for creating a comment
type CreateCommentRequest struct {
	Comment struct {
		Body string `json:"body"`
	} `json:"comment"`
}

// CommentResponse represents the response format for a single comment
type CommentResponse struct {
	Comment Comment `json:"comment"`
}

// CommentsResponse represents the response format for multiple comments
type CommentsResponse struct {
	Comments []Comment `json:"comments"`
}

// Validate validates a CreateCommentRequest
func (r *CreateCommentRequest) Validate() ValidationErrors {
	var errors ValidationErrors

	if r.Comment.Body == "" {
		errors = append(errors, ValidationError{"body", "is required"})
	} else {
		if len(r.Comment.Body) > 2000 {
			errors = append(errors, ValidationError{"body", "must be less than 2000 characters"})
		}
	}

	return errors
}

// Common errors
var (
	ErrCommentNotFound = errors.New("comment not found")
)