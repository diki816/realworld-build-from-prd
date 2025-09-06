package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/realworld/backend/internal/middleware"
	"github.com/realworld/backend/internal/models"
	"github.com/realworld/backend/internal/utils"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	DB        *sql.DB
	JWTSecret string
	Logger    *log.Logger
}

// Health handler for health checks
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	models.WriteJSONResponse(w, http.StatusOK, map[string]string{
		"status": "ok",
		"message": "RealWorld API is running",
	})
}

// Authentication handlers - implemented in Phase 1.2
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate request
	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		models.WriteErrorResponse(w, http.StatusUnprocessableEntity, validationErrors)
		return
	}

	// Check if user already exists
	var existingCount int
	err := h.DB.QueryRow(`
		SELECT COUNT(*) FROM users 
		WHERE email = ? OR username = ?
	`, req.User.Email, req.User.Username).Scan(&existingCount)
	
	if err != nil {
		h.Logger.Printf("Database error checking existing user: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if existingCount > 0 {
		// Check which field conflicts
		var emailCount, usernameCount int
		h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", req.User.Email).Scan(&emailCount)
		h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", req.User.Username).Scan(&usernameCount)
		
		var errors models.ValidationErrors
		if emailCount > 0 {
			errors = append(errors, models.ValidationError{"email", "already exists"})
		}
		if usernameCount > 0 {
			errors = append(errors, models.ValidationError{"username", "already exists"})
		}
		
		models.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.User.Password)
	if err != nil {
		h.Logger.Printf("Password hashing error: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Insert user into database
	result, err := h.DB.Exec(`
		INSERT INTO users (username, email, password_hash, bio, image) 
		VALUES (?, ?, ?, '', '')
	`, req.User.Username, req.User.Email, hashedPassword)
	
	if err != nil {
		h.Logger.Printf("Database error creating user: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Get the newly created user ID
	userID, err := result.LastInsertId()
	if err != nil {
		h.Logger.Printf("Error getting user ID: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(int(userID), req.User.Username, h.JWTSecret)
	if err != nil {
		h.Logger.Printf("Token generation error: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Create user response
	user := models.User{
		ID:       int(userID),
		Username: req.User.Username,
		Email:    req.User.Email,
		Bio:      "",
		Image:    "",
	}

	response := models.UserResponse{
		User: user.ToUserData(token),
	}

	models.WriteJSONResponse(w, http.StatusCreated, response)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate request
	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		models.WriteErrorResponse(w, http.StatusUnprocessableEntity, validationErrors)
		return
	}

	// Find user by email
	var user models.User
	var passwordHash string
	err := h.DB.QueryRow(`
		SELECT id, username, email, password_hash, bio, image, created_at, updated_at 
		FROM users WHERE email = ?
	`, req.User.Email).Scan(
		&user.ID, &user.Username, &user.Email, &passwordHash, 
		&user.Bio, &user.Image, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		models.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if err != nil {
		h.Logger.Printf("Database error during login: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Check password
	if err := utils.CheckPassword(req.User.Password, passwordHash); err != nil {
		models.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Username, h.JWTSecret)
	if err != nil {
		h.Logger.Printf("Token generation error: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Create user response
	response := models.UserResponse{
		User: user.ToUserData(token),
	}

	models.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		models.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get full user details from database
	var user models.User
	err := h.DB.QueryRow(`
		SELECT id, username, email, bio, image, created_at, updated_at 
		FROM users WHERE id = ?
	`, authUser.ID).Scan(
		&user.ID, &user.Username, &user.Email, 
		&user.Bio, &user.Image, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		models.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	if err != nil {
		h.Logger.Printf("Database error getting current user: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Generate new token to refresh expiration
	token, err := utils.GenerateToken(user.ID, user.Username, h.JWTSecret)
	if err != nil {
		h.Logger.Printf("Token generation error: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Create user response
	response := models.UserResponse{
		User: user.ToUserData(token),
	}

	models.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		models.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate request
	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		models.WriteErrorResponse(w, http.StatusUnprocessableEntity, validationErrors)
		return
	}

	// Get current user data
	var currentUser models.User
	err := h.DB.QueryRow(`
		SELECT id, username, email, bio, image, created_at, updated_at 
		FROM users WHERE id = ?
	`, authUser.ID).Scan(
		&currentUser.ID, &currentUser.Username, &currentUser.Email,
		&currentUser.Bio, &currentUser.Image, &currentUser.CreatedAt, &currentUser.UpdatedAt,
	)

	if err != nil {
		h.Logger.Printf("Database error getting current user: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Check for conflicts with existing users
	if req.User.Email != "" && req.User.Email != currentUser.Email {
		var emailCount int
		h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND id != ?", req.User.Email, authUser.ID).Scan(&emailCount)
		if emailCount > 0 {
			var errors models.ValidationErrors
			errors = append(errors, models.ValidationError{"email", "already exists"})
			models.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors)
			return
		}
	}

	if req.User.Username != "" && req.User.Username != currentUser.Username {
		var usernameCount int
		h.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND id != ?", req.User.Username, authUser.ID).Scan(&usernameCount)
		if usernameCount > 0 {
			var errors models.ValidationErrors
			errors = append(errors, models.ValidationError{"username", "already exists"})
			models.WriteErrorResponse(w, http.StatusUnprocessableEntity, errors)
			return
		}
	}

	// Prepare update values
	updateValues := make(map[string]interface{})
	if req.User.Username != "" {
		updateValues["username"] = req.User.Username
	}
	if req.User.Email != "" {
		updateValues["email"] = req.User.Email
	}
	if req.User.Bio != "" || req.User.Bio == "" { // Allow empty bio
		updateValues["bio"] = req.User.Bio
	}
	if req.User.Image != "" || req.User.Image == "" { // Allow empty image
		updateValues["image"] = req.User.Image
	}

	// Handle password update
	if req.User.Password != "" {
		hashedPassword, err := utils.HashPassword(req.User.Password)
		if err != nil {
			h.Logger.Printf("Password hashing error: %v", err)
			models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}
		updateValues["password_hash"] = hashedPassword
	}

	// Build dynamic update query
	if len(updateValues) == 0 {
		models.WriteErrorResponse(w, http.StatusBadRequest, "No fields to update")
		return
	}

	query := "UPDATE users SET "
	args := make([]interface{}, 0, len(updateValues)+1)
	setParts := make([]string, 0, len(updateValues))

	for field, value := range updateValues {
		setParts = append(setParts, field+" = ?")
		args = append(args, value)
	}

	query += strings.Join(setParts, ", ")
	query += " WHERE id = ?"
	args = append(args, authUser.ID)

	// Execute update
	_, err = h.DB.Exec(query, args...)
	if err != nil {
		h.Logger.Printf("Database error updating user: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Get updated user data
	var updatedUser models.User
	err = h.DB.QueryRow(`
		SELECT id, username, email, bio, image, created_at, updated_at 
		FROM users WHERE id = ?
	`, authUser.ID).Scan(
		&updatedUser.ID, &updatedUser.Username, &updatedUser.Email,
		&updatedUser.Bio, &updatedUser.Image, &updatedUser.CreatedAt, &updatedUser.UpdatedAt,
	)

	if err != nil {
		h.Logger.Printf("Database error getting updated user: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Generate new token with updated username if needed
	username := updatedUser.Username
	token, err := utils.GenerateToken(updatedUser.ID, username, h.JWTSecret)
	if err != nil {
		h.Logger.Printf("Token generation error: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Create user response
	response := models.UserResponse{
		User: updatedUser.ToUserData(token),
	}

	models.WriteJSONResponse(w, http.StatusOK, response)
}

// Profile handlers - implemented in Phase 1.2
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Extract username from URL path
	username := r.PathValue("username")
	if username == "" {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Username is required")
		return
	}

	// Get user profile from database
	var user models.User
	err := h.DB.QueryRow(`
		SELECT id, username, email, bio, image, created_at, updated_at 
		FROM users WHERE username = ?
	`, username).Scan(
		&user.ID, &user.Username, &user.Email,
		&user.Bio, &user.Image, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		models.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	if err != nil {
		h.Logger.Printf("Database error getting profile: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Check if current user is following this profile (if authenticated)
	following := false
	if authUser, ok := middleware.GetUserFromContext(r.Context()); ok {
		var followCount int
		h.DB.QueryRow(`
			SELECT COUNT(*) FROM follows 
			WHERE follower_id = ? AND following_id = ?
		`, authUser.ID, user.ID).Scan(&followCount)
		following = followCount > 0
	}

	// Create profile response
	response := models.ProfileResponse{
		Profile: user.ToProfile(following),
	}

	models.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		models.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Extract username from URL path
	username := r.PathValue("username")
	if username == "" {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Username is required")
		return
	}

	// Get target user
	var targetUser models.User
	err := h.DB.QueryRow(`
		SELECT id, username, email, bio, image, created_at, updated_at 
		FROM users WHERE username = ?
	`, username).Scan(
		&targetUser.ID, &targetUser.Username, &targetUser.Email,
		&targetUser.Bio, &targetUser.Image, &targetUser.CreatedAt, &targetUser.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		models.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	if err != nil {
		h.Logger.Printf("Database error getting target user: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Prevent self-following
	if authUser.ID == targetUser.ID {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Cannot follow yourself")
		return
	}

	// Check if already following
	var followCount int
	h.DB.QueryRow(`
		SELECT COUNT(*) FROM follows 
		WHERE follower_id = ? AND following_id = ?
	`, authUser.ID, targetUser.ID).Scan(&followCount)

	if followCount == 0 {
		// Create follow relationship
		_, err = h.DB.Exec(`
			INSERT INTO follows (follower_id, following_id) 
			VALUES (?, ?)
		`, authUser.ID, targetUser.ID)

		if err != nil {
			h.Logger.Printf("Database error creating follow: %v", err)
			models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	// Create profile response (always following = true after successful follow)
	response := models.ProfileResponse{
		Profile: targetUser.ToProfile(true),
	}

	models.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		models.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Extract username from URL path
	username := r.PathValue("username")
	if username == "" {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Username is required")
		return
	}

	// Get target user
	var targetUser models.User
	err := h.DB.QueryRow(`
		SELECT id, username, email, bio, image, created_at, updated_at 
		FROM users WHERE username = ?
	`, username).Scan(
		&targetUser.ID, &targetUser.Username, &targetUser.Email,
		&targetUser.Bio, &targetUser.Image, &targetUser.CreatedAt, &targetUser.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		models.WriteErrorResponse(w, http.StatusNotFound, "User not found")
		return
	}

	if err != nil {
		h.Logger.Printf("Database error getting target user: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Delete follow relationship (ignore if not following)
	_, err = h.DB.Exec(`
		DELETE FROM follows 
		WHERE follower_id = ? AND following_id = ?
	`, authUser.ID, targetUser.ID)

	if err != nil {
		h.Logger.Printf("Database error removing follow: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Create profile response (always following = false after successful unfollow)
	response := models.ProfileResponse{
		Profile: targetUser.ToProfile(false),
	}

	models.WriteJSONResponse(w, http.StatusOK, response)
}

// Article handlers - to be implemented in Phase 1.3
func (h *Handler) ListArticles(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "ListArticles endpoint not implemented yet")
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "GetFeed endpoint not implemented yet")
}

func (h *Handler) GetArticle(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "GetArticle endpoint not implemented yet")
}

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "CreateArticle endpoint not implemented yet")
}

func (h *Handler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "UpdateArticle endpoint not implemented yet")
}

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "DeleteArticle endpoint not implemented yet")
}

func (h *Handler) FavoriteArticle(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "FavoriteArticle endpoint not implemented yet")
}

func (h *Handler) UnfavoriteArticle(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "UnfavoriteArticle endpoint not implemented yet")
}

// Comment handlers - to be implemented in Phase 1.4
func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "GetComments endpoint not implemented yet")
}

func (h *Handler) CreateComment(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "CreateComment endpoint not implemented yet")
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "DeleteComment endpoint not implemented yet")
}

// Tag handlers - to be implemented in Phase 1.4
func (h *Handler) GetTags(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "GetTags endpoint not implemented yet")
}