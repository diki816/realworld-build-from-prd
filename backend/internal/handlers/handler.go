package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/realworld/backend/internal/models"
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

// Authentication handlers - to be implemented in Phase 1.2
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "Register endpoint not implemented yet")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "Login endpoint not implemented yet")
}

func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "GetCurrentUser endpoint not implemented yet")
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "UpdateUser endpoint not implemented yet")
}

// Profile handlers - to be implemented in Phase 1.2
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "GetProfile endpoint not implemented yet")
}

func (h *Handler) FollowUser(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "FollowUser endpoint not implemented yet")
}

func (h *Handler) UnfollowUser(w http.ResponseWriter, r *http.Request) {
	models.WriteErrorResponse(w, http.StatusNotImplemented, "UnfollowUser endpoint not implemented yet")
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