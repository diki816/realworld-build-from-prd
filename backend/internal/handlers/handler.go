package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

// Article handlers - implemented in Phase 1.3
func (h *Handler) ListArticles(w http.ResponseWriter, r *http.Request) {
	// Get user ID for favorite/follow status (0 if not authenticated)
	var userID int
	if authUser, ok := middleware.GetUserFromContext(r.Context()); ok {
		userID = authUser.ID
	}

	// Parse query parameters
	query := r.URL.Query()
	filters := models.ArticleFilters{
		Tag:       query.Get("tag"),
		Author:    query.Get("author"),
		Favorited: query.Get("favorited"),
		Limit:     20, // default
		Offset:    0,  // default
	}

	// Parse limit and offset
	if limitStr := query.Get("limit"); limitStr != "" {
		if limit := parseIntDefault(limitStr, 20); limit > 0 && limit <= 100 {
			filters.Limit = limit
		}
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		if offset := parseIntDefault(offsetStr, 0); offset >= 0 {
			filters.Offset = offset
		}
	}

	// Build the base query
	baseQuery := `
		SELECT DISTINCT
			a.id, a.slug, a.title, a.description, a.body, a.author_id,
			a.created_at, a.updated_at,
			u.username, u.bio, u.image,
			COALESCE(
				(SELECT COUNT(*) FROM favorites f WHERE f.article_id = a.id AND f.user_id = ?), 
				0
			) > 0 as favorited,
			(SELECT COUNT(*) FROM favorites f WHERE f.article_id = a.id) as favorites_count
		FROM articles a
		JOIN users u ON a.author_id = u.id
	`

	countQuery := `
		SELECT COUNT(DISTINCT a.id)
		FROM articles a
		JOIN users u ON a.author_id = u.id
	`

	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	var countArgs []interface{}

	args = append(args, userID)
	
	// Filter by tag
	if filters.Tag != "" {
		baseQuery += " JOIN article_tags at ON a.id = at.article_id JOIN tags t ON at.tag_id = t.id"
		countQuery += " JOIN article_tags at ON a.id = at.article_id JOIN tags t ON at.tag_id = t.id"
		conditions = append(conditions, "t.name = ?")
		args = append(args, filters.Tag)
		countArgs = append(countArgs, filters.Tag)
	}

	// Filter by author
	if filters.Author != "" {
		conditions = append(conditions, "u.username = ?")
		args = append(args, filters.Author)
		countArgs = append(countArgs, filters.Author)
	}

	// Filter by favorited user
	if filters.Favorited != "" {
		baseQuery += " JOIN favorites fav ON a.id = fav.article_id JOIN users fav_user ON fav.user_id = fav_user.id"
		countQuery += " JOIN favorites fav ON a.id = fav.article_id JOIN users fav_user ON fav.user_id = fav_user.id"
		conditions = append(conditions, "fav_user.username = ?")
		args = append(args, filters.Favorited)
		countArgs = append(countArgs, filters.Favorited)
	}

	// Add WHERE clause if conditions exist
	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Add ordering and pagination
	baseQuery += " ORDER BY a.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, filters.Limit, filters.Offset)

	// Get total count
	var totalCount int
	err := h.DB.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		h.Logger.Printf("Database error getting article count: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Get articles
	rows, err := h.DB.Query(baseQuery, args...)
	if err != nil {
		h.Logger.Printf("Database error getting articles: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		var authorUsername, authorBio, authorImage string
		var favorited bool
		var favoritesCount int

		err := rows.Scan(
			&article.ID, &article.Slug, &article.Title, &article.Description, 
			&article.Body, &article.AuthorID, &article.CreatedAt, &article.UpdatedAt,
			&authorUsername, &authorBio, &authorImage,
			&favorited, &favoritesCount,
		)
		if err != nil {
			h.Logger.Printf("Error scanning article row: %v", err)
			models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		// Check if current user follows the author
		var following bool
		if userID > 0 {
			var followCount int
			h.DB.QueryRow(`
				SELECT COUNT(*) FROM follows 
				WHERE follower_id = ? AND following_id = ?
			`, userID, article.AuthorID).Scan(&followCount)
			following = followCount > 0
		}

		// Set article fields
		article.Favorited = favorited
		article.FavoritesCount = favoritesCount
		article.Author = models.Profile{
			Username:  authorUsername,
			Bio:       authorBio,
			Image:     authorImage,
			Following: following,
		}

		// Get article tags
		tagRows, err := h.DB.Query(`
			SELECT t.name 
			FROM tags t 
			JOIN article_tags at ON t.id = at.tag_id 
			WHERE at.article_id = ?
			ORDER BY t.name
		`, article.ID)
		
		if err != nil {
			h.Logger.Printf("Error getting article tags: %v", err)
			models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		var tags []string
		for tagRows.Next() {
			var tagName string
			if err := tagRows.Scan(&tagName); err != nil {
				tagRows.Close()
				h.Logger.Printf("Error scanning tag: %v", err)
				models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			tags = append(tags, tagName)
		}
		tagRows.Close()
		
		article.TagList = tags
		if article.TagList == nil {
			article.TagList = make([]string, 0)
		}

		articles = append(articles, article)
	}

	if articles == nil {
		articles = make([]models.Article, 0)
	}

	response := models.ArticlesResponse{
		Articles:      articles,
		ArticlesCount: totalCount,
	}

	models.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) GetFeed(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		models.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse query parameters for pagination
	query := r.URL.Query()
	limit := 20 // default
	offset := 0 // default

	if limitStr := query.Get("limit"); limitStr != "" {
		if l := parseIntDefault(limitStr, 20); l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		if o := parseIntDefault(offsetStr, 0); o >= 0 {
			offset = o
		}
	}

	// Query articles from followed users
	baseQuery := `
		SELECT DISTINCT
			a.id, a.slug, a.title, a.description, a.body, a.author_id,
			a.created_at, a.updated_at,
			u.username, u.bio, u.image,
			COALESCE(
				(SELECT COUNT(*) FROM favorites f WHERE f.article_id = a.id AND f.user_id = ?), 
				0
			) > 0 as favorited,
			(SELECT COUNT(*) FROM favorites f WHERE f.article_id = a.id) as favorites_count
		FROM articles a
		JOIN users u ON a.author_id = u.id
		JOIN follows f ON a.author_id = f.following_id
		WHERE f.follower_id = ?
		ORDER BY a.created_at DESC
		LIMIT ? OFFSET ?
	`

	countQuery := `
		SELECT COUNT(DISTINCT a.id)
		FROM articles a
		JOIN follows f ON a.author_id = f.following_id
		WHERE f.follower_id = ?
	`

	// Get total count
	var totalCount int
	err := h.DB.QueryRow(countQuery, authUser.ID).Scan(&totalCount)
	if err != nil {
		h.Logger.Printf("Database error getting feed count: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Get articles
	rows, err := h.DB.Query(baseQuery, authUser.ID, authUser.ID, limit, offset)
	if err != nil {
		h.Logger.Printf("Database error getting feed: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		var authorUsername, authorBio, authorImage string
		var favorited bool
		var favoritesCount int

		err := rows.Scan(
			&article.ID, &article.Slug, &article.Title, &article.Description, 
			&article.Body, &article.AuthorID, &article.CreatedAt, &article.UpdatedAt,
			&authorUsername, &authorBio, &authorImage,
			&favorited, &favoritesCount,
		)
		if err != nil {
			h.Logger.Printf("Error scanning feed article row: %v", err)
			models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		// User is always following authors in their feed
		article.Favorited = favorited
		article.FavoritesCount = favoritesCount
		article.Author = models.Profile{
			Username:  authorUsername,
			Bio:       authorBio,
			Image:     authorImage,
			Following: true, // Always true in feed
		}

		// Get article tags
		tagRows, err := h.DB.Query(`
			SELECT t.name 
			FROM tags t 
			JOIN article_tags at ON t.id = at.tag_id 
			WHERE at.article_id = ?
			ORDER BY t.name
		`, article.ID)
		
		if err != nil {
			h.Logger.Printf("Error getting feed article tags: %v", err)
			models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		var tags []string
		for tagRows.Next() {
			var tagName string
			if err := tagRows.Scan(&tagName); err != nil {
				tagRows.Close()
				h.Logger.Printf("Error scanning feed tag: %v", err)
				models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			tags = append(tags, tagName)
		}
		tagRows.Close()
		
		article.TagList = tags
		if article.TagList == nil {
			article.TagList = make([]string, 0)
		}

		articles = append(articles, article)
	}

	if articles == nil {
		articles = make([]models.Article, 0)
	}

	response := models.ArticlesResponse{
		Articles:      articles,
		ArticlesCount: totalCount,
	}

	models.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) GetArticle(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL path
	slug := r.PathValue("slug")
	if slug == "" {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Article slug is required")
		return
	}

	// Get user ID for favorite/follow status (0 if not authenticated)
	var userID int
	if authUser, ok := middleware.GetUserFromContext(r.Context()); ok {
		userID = authUser.ID
	}

	// Get article by slug
	article, err := h.getArticleBySlug(slug, userID)
	if err == sql.ErrNoRows {
		models.WriteErrorResponse(w, http.StatusNotFound, "Article not found")
		return
	}

	if err != nil {
		h.Logger.Printf("Database error getting article: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := models.ArticleResponse{
		Article: *article,
	}

	models.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		models.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate request
	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		models.WriteErrorResponse(w, http.StatusUnprocessableEntity, validationErrors)
		return
	}

	// Generate unique slug
	checkSlugExists := func(slug string) bool {
		var count int
		h.DB.QueryRow("SELECT COUNT(*) FROM articles WHERE slug = ?", slug).Scan(&count)
		return count > 0
	}
	slug := utils.GenerateUniqueSlug(req.Article.Title, checkSlugExists)

	// Begin transaction
	tx, err := h.DB.Begin()
	if err != nil {
		h.Logger.Printf("Database error starting transaction: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer tx.Rollback()

	// Insert article
	result, err := tx.Exec(`
		INSERT INTO articles (slug, title, description, body, author_id) 
		VALUES (?, ?, ?, ?, ?)
	`, slug, req.Article.Title, req.Article.Description, req.Article.Body, authUser.ID)
	
	if err != nil {
		h.Logger.Printf("Database error creating article: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	articleID, err := result.LastInsertId()
	if err != nil {
		h.Logger.Printf("Error getting article ID: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Handle tags
	for _, tagName := range req.Article.TagList {
		if tagName == "" {
			continue
		}
		
		// Insert or get tag
		var tagID int64
		err = tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
		if err == sql.ErrNoRows {
			// Create new tag
			tagResult, err := tx.Exec("INSERT INTO tags (name) VALUES (?)", tagName)
			if err != nil {
				h.Logger.Printf("Error creating tag: %v", err)
				models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
				return
			}
			tagID, _ = tagResult.LastInsertId()
		} else if err != nil {
			h.Logger.Printf("Error querying tag: %v", err)
			models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		// Link article to tag
		_, err = tx.Exec("INSERT OR IGNORE INTO article_tags (article_id, tag_id) VALUES (?, ?)", articleID, tagID)
		if err != nil {
			h.Logger.Printf("Error linking article to tag: %v", err)
			models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		h.Logger.Printf("Error committing transaction: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Get the created article with all details
	article, err := h.getArticleBySlug(slug, authUser.ID)
	if err != nil {
		h.Logger.Printf("Error retrieving created article: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := models.ArticleResponse{
		Article: *article,
	}

	models.WriteJSONResponse(w, http.StatusCreated, response)
}

func (h *Handler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		models.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Extract slug from URL path
	slug := r.PathValue("slug")
	if slug == "" {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Article slug is required")
		return
	}

	var req models.UpdateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Validate request
	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		models.WriteErrorResponse(w, http.StatusUnprocessableEntity, validationErrors)
		return
	}

	// Get current article to verify ownership
	var currentArticle models.Article
	err := h.DB.QueryRow(`
		SELECT id, slug, title, description, body, author_id, created_at, updated_at
		FROM articles WHERE slug = ?
	`, slug).Scan(
		&currentArticle.ID, &currentArticle.Slug, &currentArticle.Title, 
		&currentArticle.Description, &currentArticle.Body, &currentArticle.AuthorID,
		&currentArticle.CreatedAt, &currentArticle.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		models.WriteErrorResponse(w, http.StatusNotFound, "Article not found")
		return
	}

	if err != nil {
		h.Logger.Printf("Database error getting article: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Check if user is the author
	if currentArticle.AuthorID != authUser.ID {
		models.WriteErrorResponse(w, http.StatusForbidden, "You can only update your own articles")
		return
	}

	// Begin transaction
	tx, err := h.DB.Begin()
	if err != nil {
		h.Logger.Printf("Database error starting transaction: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	defer tx.Rollback()

	// Prepare update values
	updateValues := make(map[string]interface{})
	newSlug := slug

	if req.Article.Title != "" && req.Article.Title != currentArticle.Title {
		updateValues["title"] = req.Article.Title
		
		// Generate new slug if title changed
		checkSlugExists := func(s string) bool {
			if s == slug {
				return false // Current slug is allowed
			}
			var count int
			h.DB.QueryRow("SELECT COUNT(*) FROM articles WHERE slug = ?", s).Scan(&count)
			return count > 0
		}
		newSlug = utils.GenerateUniqueSlug(req.Article.Title, checkSlugExists)
		updateValues["slug"] = newSlug
	}

	if req.Article.Description != "" {
		updateValues["description"] = req.Article.Description
	}

	if req.Article.Body != "" {
		updateValues["body"] = req.Article.Body
	}

	// Update article if there are changes
	if len(updateValues) > 0 {
		query := "UPDATE articles SET "
		args := make([]interface{}, 0, len(updateValues)+1)
		setParts := make([]string, 0, len(updateValues))

		for field, value := range updateValues {
			setParts = append(setParts, field+" = ?")
			args = append(args, value)
		}

		query += strings.Join(setParts, ", ")
		query += " WHERE id = ?"
		args = append(args, currentArticle.ID)

		_, err = tx.Exec(query, args...)
		if err != nil {
			h.Logger.Printf("Database error updating article: %v", err)
			models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}
	}

	// Handle tags if provided
	if req.Article.TagList != nil {
		// Remove existing tags
		_, err = tx.Exec("DELETE FROM article_tags WHERE article_id = ?", currentArticle.ID)
		if err != nil {
			h.Logger.Printf("Error removing existing tags: %v", err)
			models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		// Add new tags
		for _, tagName := range req.Article.TagList {
			if tagName == "" {
				continue
			}
			
			// Insert or get tag
			var tagID int64
			err = tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
			if err == sql.ErrNoRows {
				// Create new tag
				tagResult, err := tx.Exec("INSERT INTO tags (name) VALUES (?)", tagName)
				if err != nil {
					h.Logger.Printf("Error creating tag: %v", err)
					models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
					return
				}
				tagID, _ = tagResult.LastInsertId()
			} else if err != nil {
				h.Logger.Printf("Error querying tag: %v", err)
				models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
				return
			}

			// Link article to tag
			_, err = tx.Exec("INSERT OR IGNORE INTO article_tags (article_id, tag_id) VALUES (?, ?)", currentArticle.ID, tagID)
			if err != nil {
				h.Logger.Printf("Error linking article to tag: %v", err)
				models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
				return
			}
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		h.Logger.Printf("Error committing transaction: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Get updated article
	article, err := h.getArticleBySlug(newSlug, authUser.ID)
	if err != nil {
		h.Logger.Printf("Error retrieving updated article: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := models.ArticleResponse{
		Article: *article,
	}

	models.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		models.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Extract slug from URL path
	slug := r.PathValue("slug")
	if slug == "" {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Article slug is required")
		return
	}

	// Get article to verify ownership
	var authorID int
	err := h.DB.QueryRow(`
		SELECT author_id FROM articles WHERE slug = ?
	`, slug).Scan(&authorID)

	if err == sql.ErrNoRows {
		models.WriteErrorResponse(w, http.StatusNotFound, "Article not found")
		return
	}

	if err != nil {
		h.Logger.Printf("Database error getting article: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Check if user is the author
	if authorID != authUser.ID {
		models.WriteErrorResponse(w, http.StatusForbidden, "You can only delete your own articles")
		return
	}

	// Delete article (CASCADE will handle related records)
	_, err = h.DB.Exec("DELETE FROM articles WHERE slug = ?", slug)
	if err != nil {
		h.Logger.Printf("Database error deleting article: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Return 200 OK with empty response
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (h *Handler) FavoriteArticle(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		models.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Extract slug from URL path
	slug := r.PathValue("slug")
	if slug == "" {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Article slug is required")
		return
	}

	// Check if article exists and get its ID
	var articleID int
	err := h.DB.QueryRow("SELECT id FROM articles WHERE slug = ?", slug).Scan(&articleID)
	if err == sql.ErrNoRows {
		models.WriteErrorResponse(w, http.StatusNotFound, "Article not found")
		return
	}

	if err != nil {
		h.Logger.Printf("Database error getting article ID: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Add to favorites (ignore if already favorited)
	_, err = h.DB.Exec(`
		INSERT OR IGNORE INTO favorites (user_id, article_id) 
		VALUES (?, ?)
	`, authUser.ID, articleID)

	if err != nil {
		h.Logger.Printf("Database error favoriting article: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Get updated article
	article, err := h.getArticleBySlug(slug, authUser.ID)
	if err != nil {
		h.Logger.Printf("Error retrieving favorited article: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := models.ArticleResponse{
		Article: *article,
	}

	models.WriteJSONResponse(w, http.StatusOK, response)
}

func (h *Handler) UnfavoriteArticle(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	authUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		models.WriteErrorResponse(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Extract slug from URL path
	slug := r.PathValue("slug")
	if slug == "" {
		models.WriteErrorResponse(w, http.StatusBadRequest, "Article slug is required")
		return
	}

	// Check if article exists and get its ID
	var articleID int
	err := h.DB.QueryRow("SELECT id FROM articles WHERE slug = ?", slug).Scan(&articleID)
	if err == sql.ErrNoRows {
		models.WriteErrorResponse(w, http.StatusNotFound, "Article not found")
		return
	}

	if err != nil {
		h.Logger.Printf("Database error getting article ID: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Remove from favorites (ignore if not favorited)
	_, err = h.DB.Exec(`
		DELETE FROM favorites 
		WHERE user_id = ? AND article_id = ?
	`, authUser.ID, articleID)

	if err != nil {
		h.Logger.Printf("Database error unfavoriting article: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Get updated article
	article, err := h.getArticleBySlug(slug, authUser.ID)
	if err != nil {
		h.Logger.Printf("Error retrieving unfavorited article: %v", err)
		models.WriteErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	response := models.ArticleResponse{
		Article: *article,
	}

	models.WriteJSONResponse(w, http.StatusOK, response)
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

// Helper functions

// parseIntDefault parses a string to int with a default value
func parseIntDefault(s string, defaultValue int) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return defaultValue
}

// getArticleBySlug retrieves a complete article by slug with author profile, tags, and favorite status
func (h *Handler) getArticleBySlug(slug string, userID int) (*models.Article, error) {
	var article models.Article
	var authorUsername, authorBio, authorImage string
	var favorited bool
	var favoritesCount int
	
	// Query article with author details
	err := h.DB.QueryRow(`
		SELECT 
			a.id, a.slug, a.title, a.description, a.body, a.author_id,
			a.created_at, a.updated_at,
			u.username, u.bio, u.image,
			COALESCE(
				(SELECT COUNT(*) FROM favorites f WHERE f.article_id = a.id AND f.user_id = ?), 
				0
			) > 0 as favorited,
			(SELECT COUNT(*) FROM favorites f WHERE f.article_id = a.id) as favorites_count
		FROM articles a
		JOIN users u ON a.author_id = u.id
		WHERE a.slug = ?
	`, userID, slug).Scan(
		&article.ID, &article.Slug, &article.Title, &article.Description, 
		&article.Body, &article.AuthorID, &article.CreatedAt, &article.UpdatedAt,
		&authorUsername, &authorBio, &authorImage,
		&favorited, &favoritesCount,
	)
	
	if err != nil {
		return nil, err
	}

	// Check if current user follows the author
	var following bool
	if userID > 0 {
		var followCount int
		h.DB.QueryRow(`
			SELECT COUNT(*) FROM follows 
			WHERE follower_id = ? AND following_id = ?
		`, userID, article.AuthorID).Scan(&followCount)
		following = followCount > 0
	}

	// Set article fields
	article.Favorited = favorited
	article.FavoritesCount = favoritesCount
	article.Author = models.Profile{
		Username:  authorUsername,
		Bio:       authorBio,
		Image:     authorImage,
		Following: following,
	}

	// Get article tags
	rows, err := h.DB.Query(`
		SELECT t.name 
		FROM tags t 
		JOIN article_tags at ON t.id = at.tag_id 
		WHERE at.article_id = ?
		ORDER BY t.name
	`, article.ID)
	
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tagName string
		if err := rows.Scan(&tagName); err != nil {
			return nil, err
		}
		tags = append(tags, tagName)
	}
	
	article.TagList = tags
	if article.TagList == nil {
		article.TagList = make([]string, 0)
	}

	return &article, nil
}