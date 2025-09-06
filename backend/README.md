# RealWorld Backend (Go)

A Go implementation of the RealWorld API specification.

## Features

- Clean Architecture with Go standard library
- SQLite database with migrations
- JWT authentication
- Comprehensive validation
- CORS middleware
- Rate limiting
- Structured logging

## Project Structure

```
backend/
├── cmd/server/          # Application entry point
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models and validation
│   ├── database/        # Database connection and migrations
│   └── utils/           # Utility functions
└── tests/               # Test files
```

## Quick Start

1. Install Go 1.25+
2. Copy `.env.example` to `.env` and configure
3. Run the server:

```bash
go run cmd/server/main.go
```

## Environment Variables

- `PORT`: Server port (default: 8080)
- `DB_PATH`: SQLite database file path
- `JWT_SECRET`: Secret key for JWT tokens
- `LOG_LEVEL`: Logging level (debug, info, warn, error)

## API Endpoints

The API follows the [RealWorld specification](https://realworld-docs.netlify.app/docs/specs/backend-specs/introduction).

### Authentication
- `POST /api/users/login` - User login
- `POST /api/users` - User registration
- `GET /api/user` - Get current user
- `PUT /api/user` - Update user

### Profiles
- `GET /api/profiles/:username` - Get user profile
- `POST /api/profiles/:username/follow` - Follow user
- `DELETE /api/profiles/:username/follow` - Unfollow user

### Articles
- `GET /api/articles` - List articles
- `GET /api/articles/feed` - Get user feed
- `GET /api/articles/:slug` - Get single article
- `POST /api/articles` - Create article
- `PUT /api/articles/:slug` - Update article
- `DELETE /api/articles/:slug` - Delete article
- `POST /api/articles/:slug/favorite` - Favorite article
- `DELETE /api/articles/:slug/favorite` - Unfavorite article

### Comments
- `GET /api/articles/:slug/comments` - Get article comments
- `POST /api/articles/:slug/comments` - Add comment
- `DELETE /api/articles/:slug/comments/:id` - Delete comment

### Tags
- `GET /api/tags` - Get all tags

## Development

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o realworld cmd/server/main.go
```

### Database Migrations

Migrations run automatically on server start. Migration files are in `internal/database/migrations/`.

## Implementation Status

- ✅ Phase 1.1: Project setup and foundation
- 🔄 Phase 1.2: Authentication and profiles (next)
- ⏳ Phase 1.3: Articles system
- ⏳ Phase 1.4: Comments and tags