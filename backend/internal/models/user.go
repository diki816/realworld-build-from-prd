package models

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Bio       string    `json:"bio" db:"bio"`
	Image     string    `json:"image" db:"image"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// Profile represents a user profile (public view)
type Profile struct {
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

// RegisterRequest represents the request payload for user registration
type RegisterRequest struct {
	User struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"user"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	User struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	} `json:"user"`
}

// UpdateUserRequest represents the request payload for updating user profile
type UpdateUserRequest struct {
	User struct {
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
		Password string `json:"password,omitempty"`
		Bio      string `json:"bio,omitempty"`
		Image    string `json:"image,omitempty"`
	} `json:"user"`
}

// UserResponse represents the response format for user data
type UserResponse struct {
	User UserData `json:"user"`
}

// UserData represents the user data in responses
type UserData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
	Token    string `json:"token"`
}

// ProfileResponse represents the response format for profile data
type ProfileResponse struct {
	Profile Profile `json:"profile"`
}

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string
	Message string
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	var messages []string
	for _, err := range ve {
		messages = append(messages, err.Field+": "+err.Message)
	}
	return strings.Join(messages, ", ")
}

// Validate validates a RegisterRequest
func (r *RegisterRequest) Validate() ValidationErrors {
	var errors ValidationErrors

	// Username validation
	if r.User.Username == "" {
		errors = append(errors, ValidationError{"username", "is required"})
	} else {
		if len(r.User.Username) < 3 {
			errors = append(errors, ValidationError{"username", "must be at least 3 characters long"})
		}
		if len(r.User.Username) > 50 {
			errors = append(errors, ValidationError{"username", "must be less than 50 characters"})
		}
		// Check for valid characters (alphanumeric, underscore, hyphen)
		if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, r.User.Username); !matched {
			errors = append(errors, ValidationError{"username", "can only contain letters, numbers, underscores, and hyphens"})
		}
	}

	// Email validation
	if r.User.Email == "" {
		errors = append(errors, ValidationError{"email", "is required"})
	} else {
		if !isValidEmail(r.User.Email) {
			errors = append(errors, ValidationError{"email", "is invalid"})
		}
	}

	// Password validation
	if r.User.Password == "" {
		errors = append(errors, ValidationError{"password", "is required"})
	} else {
		if len(r.User.Password) < 6 {
			errors = append(errors, ValidationError{"password", "must be at least 6 characters long"})
		}
		if len(r.User.Password) > 128 {
			errors = append(errors, ValidationError{"password", "must be less than 128 characters"})
		}
	}

	return errors
}

// Validate validates a LoginRequest
func (l *LoginRequest) Validate() ValidationErrors {
	var errors ValidationErrors

	if l.User.Email == "" {
		errors = append(errors, ValidationError{"email", "is required"})
	} else {
		if !isValidEmail(l.User.Email) {
			errors = append(errors, ValidationError{"email", "is invalid"})
		}
	}

	if l.User.Password == "" {
		errors = append(errors, ValidationError{"password", "is required"})
	}

	return errors
}

// Validate validates an UpdateUserRequest
func (u *UpdateUserRequest) Validate() ValidationErrors {
	var errors ValidationErrors

	// Username validation (optional)
	if u.User.Username != "" {
		if len(u.User.Username) < 3 {
			errors = append(errors, ValidationError{"username", "must be at least 3 characters long"})
		}
		if len(u.User.Username) > 50 {
			errors = append(errors, ValidationError{"username", "must be less than 50 characters"})
		}
		if matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, u.User.Username); !matched {
			errors = append(errors, ValidationError{"username", "can only contain letters, numbers, underscores, and hyphens"})
		}
	}

	// Email validation (optional)
	if u.User.Email != "" {
		if !isValidEmail(u.User.Email) {
			errors = append(errors, ValidationError{"email", "is invalid"})
		}
	}

	// Password validation (optional)
	if u.User.Password != "" {
		if len(u.User.Password) < 6 {
			errors = append(errors, ValidationError{"password", "must be at least 6 characters long"})
		}
		if len(u.User.Password) > 128 {
			errors = append(errors, ValidationError{"password", "must be less than 128 characters"})
		}
	}

	// Bio validation (optional)
	if len(u.User.Bio) > 1000 {
		errors = append(errors, ValidationError{"bio", "must be less than 1000 characters"})
	}

	// Image URL validation (optional)
	if u.User.Image != "" {
		if len(u.User.Image) > 500 {
			errors = append(errors, ValidationError{"image", "URL must be less than 500 characters"})
		}
		if !isValidURL(u.User.Image) {
			errors = append(errors, ValidationError{"image", "must be a valid URL"})
		}
	}

	return errors
}

// ToUserData converts a User model to UserData for API responses
func (u *User) ToUserData(token string) UserData {
	return UserData{
		Username: u.Username,
		Email:    u.Email,
		Bio:      u.Bio,
		Image:    u.Image,
		Token:    token,
	}
}

// ToProfile converts a User model to Profile for API responses
func (u *User) ToProfile(following bool) Profile {
	return Profile{
		Username:  u.Username,
		Bio:       u.Bio,
		Image:     u.Image,
		Following: following,
	}
}

// Helper function to validate email format
func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email) && len(email) <= 254
}

// Helper function to validate URL format
func isValidURL(url string) bool {
	urlRegex := regexp.MustCompile(`^https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_\+.~#?&//=]*)$`)
	return urlRegex.MatchString(url)
}

// Common errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
)