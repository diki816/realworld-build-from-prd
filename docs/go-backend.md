# Go Backend Documentation

## Go Language

### Version Information
- **Current Stable**: Go 1.25.1 (Latest Patch)
- **Recommended**: Go 1.25.x or Go 1.24.x for production
- **Support Model**: Rolling support (each version supported until 2 newer majors released)
- **Migration Required**: From Go 1.21+ specified in PRD

### Official Documentation
- **Main Documentation**: https://golang.org/doc/
- **Language Specification**: https://golang.org/ref/spec
- **Standard Library**: https://pkg.go.dev/std
- **Effective Go**: https://golang.org/doc/effective_go
- **Go Blog**: https://blog.golang.org/
- **Tour of Go**: https://tour.golang.org/
- **GitHub Repository**: https://github.com/golang/go

### Key Features for RealWorld
- **Simplicity**: Minimal syntax, easy to read and maintain
- **Concurrency**: Built-in goroutines and channels for concurrent operations
- **Standard Library**: Rich standard library including `net/http`
- **Performance**: Compiled language with excellent runtime performance
- **Cross-Platform**: Build for multiple operating systems and architectures
- **Memory Safety**: Garbage collected with efficient memory management

### Installation & Setup
```bash
# Install Go 1.25.1
# Download from https://golang.org/dl/

# Verify installation
go version
# go version go1.25.1 linux/amd64

# Initialize Go module
cd backend
go mod init github.com/youruser/realworld-backend

# Set Go environment variables
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1  # Required for SQLite
```

### Go Module Configuration
```go
// go.mod
module github.com/youruser/realworld-backend

go 1.25

require (
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/mattn/go-sqlite3 v1.14.17
    golang.org/x/crypto v0.17.0
)

require (
    github.com/davecgh/go-spew v1.1.1 // indirect
    github.com/pmezard/go-difflib v1.0.0 // indirect
    github.com/stretchr/testify v1.8.4 // test
    gopkg.in/yaml.v3 v3.0.1 // indirect
)
```

## HTTP Server with Standard Library

### Core HTTP Features
- **Native HTTP Server**: `net/http` package for HTTP server functionality
- **Middleware Support**: Request/response middleware pattern
- **Context Propagation**: Built-in context package for request scoping
- **TLS Support**: HTTPS with automatic certificate management
- **HTTP/2**: Built-in HTTP/2 support

### Basic Server Setup
```go
// cmd/server/main.go
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

    "github.com/youruser/realworld-backend/internal/database"
    "github.com/youruser/realworld-backend/internal/handlers"
    "github.com/youruser/realworld-backend/internal/middleware"
    
    _ "github.com/mattn/go-sqlite3" // SQLite driver
)

func main() {
    // Environment configuration
    port := getEnv("PORT", "8080")
    dbPath := getEnv("DB_PATH", "./data/realworld.db")
    jwtSecret := getEnv("JWT_SECRET", "your-development-secret")
    
    // Initialize database
    db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
    if err != nil {
        log.Fatal("Failed to open database:", err)
    }
    defer db.Close()
    
    // Run database migrations
    if err := database.Migrate(db); err != nil {
        log.Fatal("Failed to migrate database:", err)
    }
    
    // Initialize handlers
    h := &handlers.Handler{
        DB:        db,
        JWTSecret: jwtSecret,
        Logger:    log.New(os.Stdout, "realworld-api: ", log.LstdFlags),
    }
    
    // Setup routes
    mux := setupRoutes(h)
    
    // Setup middleware
    handler := middleware.Chain(mux,
        middleware.CORS(),
        middleware.Logging(h.Logger),
        middleware.Recovery(h.Logger),
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
    
    // Graceful shutdown
    go func() {
        log.Printf("Server starting on port %s", port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatal("Failed to start server:", err)
        }
    }()
    
    // Wait for interrupt signal
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := server.Shutdown(ctx); err != nil {
        log.Fatal("Server forced to shutdown:", err)
    }
    
    log.Println("Server exited")
}

func setupRoutes(h *handlers.Handler) *http.ServeMux {
    mux := http.NewServeMux()
    
    // Health check
    mux.HandleFunc("GET /health", h.Health)
    
    // Authentication routes
    mux.HandleFunc("POST /api/users/login", h.Login)
    mux.HandleFunc("POST /api/users", h.Register)
    
    // Protected user routes
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
```

### Handler Architecture
```go
// internal/handlers/handler.go
package handlers

import (
    "database/sql"
    "log"
)

type Handler struct {
    DB        *sql.DB
    JWTSecret string
    Logger    *log.Logger
}

// Response helpers
type Response struct {
    User    *User    `json:"user,omitempty"`
    Profile *Profile `json:"profile,omitempty"`
    Article *Article `json:"article,omitempty"`
    Articles struct {
        Articles      []Article `json:"articles"`
        ArticlesCount int       `json:"articlesCount"`
    } `json:",omitempty"`
    Comments []Comment `json:"comments,omitempty"`
    Tags     []string  `json:"tags,omitempty"`
    Errors   map[string][]string `json:"errors,omitempty"`
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(data); err != nil {
        h.Logger.Printf("Error encoding JSON: %v", err)
    }
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
    response := Response{
        Errors: map[string][]string{
            "body": {message},
        },
    }
    h.writeJSON(w, status, response)
}
```

## JWT Authentication

### JWT Implementation
```go
// internal/utils/jwt.go
package utils

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID   int    `json:"user_id"`
    Username string `json:"username"`
    jwt.RegisteredClaims
}

func GenerateToken(userID int, username, secret string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Username: username,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "realworld-api",
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString, secret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("invalid signing method")
        }
        return []byte(secret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}
```

### Authentication Middleware
```go
// internal/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/youruser/realworld-backend/internal/utils"
)

type contextKey string

const UserContextKey = contextKey("user")

type User struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

func Auth(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization header required", http.StatusUnauthorized)
                return
            }

            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
                return
            }

            token := parts[1]
            claims, err := utils.ValidateToken(token, secret)
            if err != nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            user := &User{
                ID:       claims.UserID,
                Username: claims.Username,
            }

            ctx := context.WithValue(r.Context(), UserContextKey, user)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetUserFromContext(ctx context.Context) (*User, bool) {
    user, ok := ctx.Value(UserContextKey).(*User)
    return user, ok
}
```

## Error Handling Patterns

### Structured Error Handling
```go
// internal/errors/errors.go
package errors

import (
    "fmt"
    "net/http"
)

type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

func (e APIError) Error() string {
    return e.Message
}

var (
    ErrUserNotFound     = APIError{http.StatusNotFound, "User not found", ""}
    ErrArticleNotFound  = APIError{http.StatusNotFound, "Article not found", ""}
    ErrUnauthorized     = APIError{http.StatusUnauthorized, "Unauthorized", ""}
    ErrForbidden        = APIError{http.StatusForbidden, "Forbidden", ""}
    ErrInvalidInput     = APIError{http.StatusBadRequest, "Invalid input", ""}
    ErrInternalServer   = APIError{http.StatusInternalServerError, "Internal server error", ""}
)

func NewValidationError(field, message string) APIError {
    return APIError{
        Code:    http.StatusBadRequest,
        Message: fmt.Sprintf("Validation error: %s", message),
        Details: field,
    }
}
```

### Middleware Chain
```go
// internal/middleware/middleware.go
package middleware

import (
    "log"
    "net/http"
    "time"
)

func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

func CORS() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}

func Logging(logger *log.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            
            next.ServeHTTP(w, r)
            
            logger.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
        })
    }
}

func Recovery(logger *log.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            defer func() {
                if err := recover(); err != nil {
                    logger.Printf("Panic recovered: %v", err)
                    http.Error(w, "Internal server error", http.StatusInternalServerError)
                }
            }()
            
            next.ServeHTTP(w, r)
        })
    }
}
```

## Performance Optimization

### Database Connection Pooling
```go
// Configure database connection pool
func setupDB(dbPath string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=on&_journal_mode=WAL&_cache_size=1000")
    if err != nil {
        return nil, err
    }

    // Connection pool settings
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(25)
    db.SetConnMaxLifetime(5 * time.Minute)

    // Test connection
    if err := db.Ping(); err != nil {
        return nil, err
    }

    return db, nil
}
```

### Production Configuration
```go
// Production environment configuration
func setupProduction() {
    // Set GOMAXPROCS to match container limits
    runtime.GOMAXPROCS(runtime.NumCPU())
    
    // Configure garbage collector
    debug.SetGCPercent(100)
    
    // Enable CPU profiling in development
    if os.Getenv("ENABLE_PPROF") == "true" {
        go func() {
            log.Println(http.ListenAndServe("localhost:6060", nil))
        }()
    }
}
```

This Go backend implementation provides a robust, performant foundation for the RealWorld API following Go best practices and the "dumbest possible thing that works" philosophy.