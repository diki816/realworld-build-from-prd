package models

import (
	"errors"
	"time"
)

// Article represents an article in the system
type Article struct {
	ID             int       `json:"id" db:"id"`
	Slug           string    `json:"slug" db:"slug"`
	Title          string    `json:"title" db:"title"`
	Description    string    `json:"description" db:"description"`
	Body           string    `json:"body" db:"body"`
	AuthorID       int       `json:"-" db:"author_id"`
	CreatedAt      time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" db:"updated_at"`
	Favorited      bool      `json:"favorited"`
	FavoritesCount int       `json:"favoritesCount"`
	TagList        []string  `json:"tagList"`
	Author         Profile   `json:"author"`
}

// CreateArticleRequest represents the request payload for creating an article
type CreateArticleRequest struct {
	Article struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Body        string   `json:"body"`
		TagList     []string `json:"tagList"`
	} `json:"article"`
}

// UpdateArticleRequest represents the request payload for updating an article
type UpdateArticleRequest struct {
	Article struct {
		Title       string   `json:"title,omitempty"`
		Description string   `json:"description,omitempty"`
		Body        string   `json:"body,omitempty"`
		TagList     []string `json:"tagList,omitempty"`
	} `json:"article"`
}

// ArticleResponse represents the response format for a single article
type ArticleResponse struct {
	Article Article `json:"article"`
}

// ArticlesResponse represents the response format for multiple articles
type ArticlesResponse struct {
	Articles      []Article `json:"articles"`
	ArticlesCount int       `json:"articlesCount"`
}

// ArticleFilters represents filters for querying articles
type ArticleFilters struct {
	Tag        string `json:"tag"`
	Author     string `json:"author"`
	Favorited  string `json:"favorited"`
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

// Validate validates a CreateArticleRequest
func (r *CreateArticleRequest) Validate() ValidationErrors {
	var errors ValidationErrors

	if r.Article.Title == "" {
		errors = append(errors, ValidationError{"title", "is required"})
	} else {
		if len(r.Article.Title) > 255 {
			errors = append(errors, ValidationError{"title", "must be less than 255 characters"})
		}
	}

	if r.Article.Description == "" {
		errors = append(errors, ValidationError{"description", "is required"})
	} else {
		if len(r.Article.Description) > 500 {
			errors = append(errors, ValidationError{"description", "must be less than 500 characters"})
		}
	}

	if r.Article.Body == "" {
		errors = append(errors, ValidationError{"body", "is required"})
	}

	// Validate tags
	if len(r.Article.TagList) > 10 {
		errors = append(errors, ValidationError{"tagList", "cannot have more than 10 tags"})
	}

	for _, tag := range r.Article.TagList {
		if len(tag) > 50 {
			errors = append(errors, ValidationError{"tagList", "each tag must be less than 50 characters"})
		}
		if tag == "" {
			errors = append(errors, ValidationError{"tagList", "tags cannot be empty"})
		}
	}

	return errors
}

// Validate validates an UpdateArticleRequest
func (r *UpdateArticleRequest) Validate() ValidationErrors {
	var errors ValidationErrors

	if r.Article.Title != "" && len(r.Article.Title) > 255 {
		errors = append(errors, ValidationError{"title", "must be less than 255 characters"})
	}

	if r.Article.Description != "" && len(r.Article.Description) > 500 {
		errors = append(errors, ValidationError{"description", "must be less than 500 characters"})
	}

	// Validate tags if provided
	if len(r.Article.TagList) > 10 {
		errors = append(errors, ValidationError{"tagList", "cannot have more than 10 tags"})
	}

	for _, tag := range r.Article.TagList {
		if len(tag) > 50 {
			errors = append(errors, ValidationError{"tagList", "each tag must be less than 50 characters"})
		}
		if tag == "" {
			errors = append(errors, ValidationError{"tagList", "tags cannot be empty"})
		}
	}

	return errors
}

// Common errors
var (
	ErrArticleNotFound = errors.New("article not found")
	ErrSlugExists      = errors.New("article with this slug already exists")
	ErrNotAuthorized   = errors.New("not authorized to perform this action")
)