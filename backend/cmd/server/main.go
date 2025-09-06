package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/realworld/backend/internal/database"
	"github.com/realworld/backend/internal/handlers"
	"github.com/realworld/backend/internal/middleware"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
	// Environment configuration
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "./data/realworld.db")
	jwtSecret := getEnv("JWT_SECRET", "your-development-secret-change-in-production")

	// Initialize logger
	logger := log.New(os.Stdout, "realworld-api: ", log.LstdFlags)

	// Initialize database
	db, err := database.New(dbPath)
	if err != nil {
		logger.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	logger.Println("Database initialized successfully")

	// Initialize handlers
	h := &handlers.Handler{
		DB:        db.DB,
		JWTSecret: jwtSecret,
		Logger:    logger,
	}

	// Setup routes
	mux := setupRoutes(h)

	// Setup middleware chain
	handler := middleware.Chain(mux,
		middleware.CORS(),
		middleware.Logging(logger),
		middleware.Recovery(logger),
		middleware.RateLimit(),
	)

	// HTTP server configuration
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Printf("Server starting on port %s", port)
		logger.Printf("API available at: http://localhost:%s/api", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Println("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown:", err)
	}

	logger.Println("Server exited")
}

func setupRoutes(h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", h.Health)

	// Authentication routes - public
	mux.HandleFunc("POST /api/users/login", h.Login)
	mux.HandleFunc("POST /api/users", h.Register)

	// User routes - protected
	mux.Handle("GET /api/user", middleware.Auth(h.JWTSecret)(http.HandlerFunc(h.GetCurrentUser)))
	mux.Handle("PUT /api/user", middleware.Auth(h.JWTSecret)(http.HandlerFunc(h.UpdateUser)))

	// Profile routes
	mux.HandleFunc("GET /api/profiles/{username}", h.GetProfile)
	mux.Handle("POST /api/profiles/{username}/follow", middleware.Auth(h.JWTSecret)(http.HandlerFunc(h.FollowUser)))
	mux.Handle("DELETE /api/profiles/{username}/follow", middleware.Auth(h.JWTSecret)(http.HandlerFunc(h.UnfollowUser)))

	// Article routes
	mux.HandleFunc("GET /api/articles", h.ListArticles)
	mux.HandleFunc("GET /api/articles/{slug}", h.GetArticle)
	mux.Handle("GET /api/articles/feed", middleware.Auth(h.JWTSecret)(http.HandlerFunc(h.GetFeed)))
	mux.Handle("POST /api/articles", middleware.Auth(h.JWTSecret)(http.HandlerFunc(h.CreateArticle)))
	mux.Handle("PUT /api/articles/{slug}", middleware.Auth(h.JWTSecret)(http.HandlerFunc(h.UpdateArticle)))
	mux.Handle("DELETE /api/articles/{slug}", middleware.Auth(h.JWTSecret)(http.HandlerFunc(h.DeleteArticle)))

	// Favorite routes
	mux.Handle("POST /api/articles/{slug}/favorite", middleware.Auth(h.JWTSecret)(http.HandlerFunc(h.FavoriteArticle)))
	mux.Handle("DELETE /api/articles/{slug}/favorite", middleware.Auth(h.JWTSecret)(http.HandlerFunc(h.UnfavoriteArticle)))

	// Comment routes
	mux.HandleFunc("GET /api/articles/{slug}/comments", h.GetComments)
	mux.Handle("POST /api/articles/{slug}/comments", middleware.Auth(h.JWTSecret)(http.HandlerFunc(h.CreateComment)))
	mux.Handle("DELETE /api/articles/{slug}/comments/{id}", middleware.Auth(h.JWTSecret)(http.HandlerFunc(h.DeleteComment)))

	// Tag routes
	mux.HandleFunc("GET /api/tags", h.GetTags)

	return mux
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}