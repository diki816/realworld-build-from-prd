# RealWorld App - Product Requirements Document (PRD)

## 1. Project Overview

### 1.1 Product Name
RealWorld (Conduit) - Medium.com Clone

### 1.2 Project Description
RealWorld is a full-stack web application that serves as an exemplary implementation of a social blogging platform, modeled after Medium.com. The project demonstrates how to build production-ready applications using various frontend and backend technology stacks while maintaining a consistent API specification.

### 1.3 Purpose & Goals
- Provide a comprehensive, real-world example beyond basic "todo" applications
- Demonstrate technology interoperability across different frameworks
- Offer developers a practical learning resource for full-stack development
- Enable mix-and-match of frontend and backend implementations

## 2. Product Vision

### 2.1 Vision Statement
Create an open-source platform that showcases how the exact same Medium.com clone can be built using different frontends (React, Angular, etc.) and backends (Node.js, Django, etc.) while maintaining consistent functionality and user experience.

### 2.2 Target Audience
- **Primary**: Developers learning full-stack web development
- **Secondary**: Technology enthusiasts exploring framework interoperability
- **Tertiary**: Software engineers seeking practical implementation examples

## 3. Core Features & Functionality

### 3.1 User Authentication & Management
- **User Registration**: Email, username, and password-based signup
- **User Login**: JWT-based authentication system
- **User Profile Management**: Update email, username, password, bio, and avatar
- **User Profiles**: Public profile pages with user information and articles

### 3.2 Article Management
- **Article Creation**: Rich text editor with markdown support
- **Article Editing**: Update existing articles (author-only)
- **Article Deletion**: Remove articles (author-only)
- **Article Viewing**: Public article pages with markdown rendering
- **Article Feed**: Personalized feed for authenticated users
- **Global Article List**: Public article listing with pagination
- **Article Favoriting**: Like/unlike articles functionality

### 3.3 Social Features
- **User Following**: Follow/unfollow other users
- **Comments System**: Add, view, and delete comments on articles
- **Tag System**: Categorize articles with tags
- **Popular Tags**: Display trending tags

### 3.4 Content Discovery
- **Article Filtering**: Filter by tag, author, or favorited status
- **Pagination**: Navigate through article lists
- **Search**: Tag-based content discovery
- **Feed Types**: Global feed, personal feed, and tag-specific feeds

## 4. Technical Specifications

### 4.1 Architecture
- **Type**: Single Page Application (SPA)
- **Architecture Pattern**: Client-server with RESTful API
- **Authentication**: JWT-based token authentication
- **Routing**: Hash-based routing (`#/path`)

### 4.2 Frontend Requirements

#### 4.2.1 Technology Stack
- **Framework**: React 18+ with TypeScript
- **Build Tool**: Vite for fast development and optimized builds
- **State Management**: TanStack Query for server state, React Context for client state
- **Routing**: TanStack Router for type-safe routing
- **Development**: Hot Module Replacement (HMR) for fast iteration

#### 4.2.2 Styling & UI System
- **CSS Framework**: Tailwind CSS for utility-first styling
- **Component Library**: Shadcn/UI for pre-built, customizable React components
- **Icons**: Lucide React (integrated with Shadcn/UI)
- **Theme System**: Built-in dark/light mode support
- **Typography**: System fonts with fallbacks for optimal performance
- **Design Tokens**: Consistent spacing, colors, and typography scales
- **Responsive Design**: Mobile-first approach with Tailwind breakpoints

#### 4.2.3 Component Architecture
- **Design System**: Atomic design principles (atoms, molecules, organisms)
- **Reusable Components**: Button, Input, Card, Modal, etc. from Shadcn/UI
- **Custom Components**: Article editor, comment system, user feed
- **Accessibility**: WCAG 2.1 AA compliance with proper ARIA attributes
- **Performance**: Code splitting and lazy loading for optimal bundle size

#### 4.2.3 Routing Structure
```
/#/                           # Home page
/#/login                      # Login page
/#/register                   # Registration page
/#/settings                   # User settings
/#/editor                     # Create new article
/#/editor/:slug               # Edit existing article
/#/article/:slug              # Article detail page
/#/profile/:username          # User profile
/#/profile/:username/favorites # User's favorited articles
```

### 4.3 Backend Technical Specification

#### 4.3.1 Technology Stack
- **Language**: Go 1.21+ for simplicity and performance
- **HTTP Framework**: Standard library `net/http` with lightweight routing
- **Database**: SQLite for lightweight, file-based storage
- **Database Access**: Plain SQL with prepared statements (no ORM)
- **Authentication**: JWT with secure signing and validation
- **Logging**: Structured logging with configurable levels
- **Configuration**: Environment variables with sensible defaults

#### 4.3.2 Architecture Principles
- **Functional Design**: Prefer functions over complex classes
- **Explicit Context**: Use Go's context package for request scoping
- **Simple Error Handling**: Clear error messages and proper HTTP status codes
- **Local Permissions**: Permission checks close to business logic
- **Minimal Dependencies**: "Dumbest possible thing that works" approach
- **Agent-Friendly**: Clear interfaces and predictable behavior

#### 4.3.3 API Configuration
- **Base URL**: `http://localhost:8080/api` (development)
- **Content-Type**: `application/json; charset=utf-8`
- **CORS**: Configured for frontend development
- **Rate Limiting**: Basic rate limiting for API protection
- **Request/Response**: JSON format with consistent structure

#### 4.3.4 Authentication System
- **Method**: JWT (JSON Web Token) with HS256 signing
- **Header Format**: `Authorization: Bearer <jwt-token>`
- **Token Storage**: localStorage (frontend) with secure flags
- **Token Expiration**: 7 days with refresh mechanism
- **Password Security**: bcrypt hashing with appropriate cost
- **Session Management**: Stateless JWT-based authentication

#### 4.3.3 API Endpoints

**Authentication**
- `POST /api/users/login` - User login
- `POST /api/users` - User registration
- `GET /api/user` - Get current user (auth required)
- `PUT /api/user` - Update user (auth required)

**Profiles**
- `GET /api/profiles/:username` - Get user profile
- `POST /api/profiles/:username/follow` - Follow user (auth required)
- `DELETE /api/profiles/:username/follow` - Unfollow user (auth required)

**Articles**
- `GET /api/articles` - List articles (with filtering & pagination)
- `GET /api/articles/feed` - Get user feed (auth required)
- `GET /api/articles/:slug` - Get single article
- `POST /api/articles` - Create article (auth required)
- `PUT /api/articles/:slug` - Update article (auth required)
- `DELETE /api/articles/:slug` - Delete article (auth required)

**Favorites**
- `POST /api/articles/:slug/favorite` - Favorite article (auth required)
- `DELETE /api/articles/:slug/favorite` - Unfavorite article (auth required)

**Comments**
- `GET /api/articles/:slug/comments` - Get article comments
- `POST /api/articles/:slug/comments` - Add comment (auth required)
- `DELETE /api/articles/:slug/comments/:id` - Delete comment (auth required)

**Tags**
- `GET /api/tags` - Get all tags

### 4.4 Database Design (SQLite)

#### 4.4.1 Database Schema
```sql
-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    bio TEXT,
    image VARCHAR(500),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Articles table
CREATE TABLE articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug VARCHAR(255) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    body TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES users(id)
);

-- Tags table
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL
);

-- Article tags junction table
CREATE TABLE article_tags (
    article_id INTEGER,
    tag_id INTEGER,
    PRIMARY KEY (article_id, tag_id),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Comments table
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    body TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    article_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES users(id),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
);

-- Favorites table
CREATE TABLE favorites (
    user_id INTEGER,
    article_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, article_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
);

-- Follows table
CREATE TABLE follows (
    follower_id INTEGER,
    following_id INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (follower_id, following_id),
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE
);
```

#### 4.4.2 Indexes for Performance
```sql
CREATE INDEX idx_articles_author_id ON articles(author_id);
CREATE INDEX idx_articles_created_at ON articles(created_at DESC);
CREATE INDEX idx_comments_article_id ON comments(article_id);
CREATE INDEX idx_favorites_article_id ON favorites(article_id);
CREATE INDEX idx_follows_following_id ON follows(following_id);
```

### 4.5 Data Models (JSON API)

#### 4.4.1 User
```json
{
  "email": "string",
  "username": "string", 
  "password": "string",
  "bio": "string",
  "image": "string"
}
```

#### 4.4.2 Article
```json
{
  "title": "string",
  "description": "string",
  "body": "string",
  "tagList": ["string"],
  "slug": "string",
  "createdAt": "datetime",
  "updatedAt": "datetime",
  "favorited": "boolean",
  "favoritesCount": "number",
  "author": "Profile"
}
```

#### 4.4.3 Comment
```json
{
  "id": "number",
  "body": "string",
  "createdAt": "datetime",
  "updatedAt": "datetime",
  "author": "Profile"
}
```

## 5. User Stories

### 5.1 Guest User Stories
- As a guest, I can view the global article feed
- As a guest, I can view individual articles
- As a guest, I can view user profiles
- As a guest, I can register for a new account
- As a guest, I can login to my account

### 5.2 Authenticated User Stories
- As a user, I can view my personalized article feed
- As a user, I can create new articles with markdown
- As a user, I can edit my own articles
- As a user, I can delete my own articles
- As a user, I can favorite/unfavorite articles
- As a user, I can follow/unfollow other users
- As a user, I can comment on articles
- As a user, I can delete my own comments
- As a user, I can update my profile settings
- As a user, I can view my profile and favorited articles

## 6. User Interface Requirements

### 6.1 Page Components

#### 6.1.1 Header
- Logo/Brand name
- Navigation menu (conditional based on auth status)
- User menu (when authenticated)

#### 6.1.2 Home Page
- Article feed tabs (Your Feed, Global Feed, Tag Feed)
- Article previews with metadata
- Popular tags sidebar
- Pagination controls

#### 6.1.3 Authentication Pages
- Login form (email, password)
- Registration form (username, email, password)
- Form validation and error handling

#### 6.1.4 Article Pages
- Article content with markdown rendering
- Article metadata (author, date, tags)
- Follow/favorite buttons
- Comments section
- Author-only edit/delete controls

#### 6.1.5 Profile Pages
- User information display
- Article tabs (My Articles, Favorited Articles)
- Follow button (for other users)

#### 6.1.6 Editor
- Article title input
- Article description input
- Markdown editor for body
- Tags input
- Publish button

### 6.2 Responsive Design
- Mobile-first approach
- Breakpoints for tablet and desktop
- Touch-friendly interface elements

## 7. Performance Requirements

### 7.1 Frontend Performance
- Fast initial page load
- Smooth client-side navigation
- Efficient API data caching
- Optimized asset loading

### 7.2 Backend Performance
- API response time < 200ms for most endpoints
- Support for pagination to handle large datasets
- Efficient database queries
- CORS support for cross-origin requests

## 8. Security Requirements

### 8.1 Authentication Security
- Secure password storage (hashing)
- JWT token expiration
- Protected routes requiring authentication
- Input validation and sanitization

### 8.2 Content Security
- XSS protection for user-generated content
- CSRF protection
- Safe markdown rendering
- Content moderation capabilities

## 9. Implementation Guidelines

### 9.1 Development Standards
- Follow existing codebase conventions
- Implement consistent error handling
- Use proper HTTP status codes
- Maintain API specification compliance

### 9.2 Testing Strategy

#### 9.2.1 Frontend Testing
- **Unit Tests**: Vitest for component and utility testing
- **Component Tests**: React Testing Library for user interaction testing
- **Integration Tests**: API integration with MSW (Mock Service Worker)
- **E2E Tests**: Playwright for critical user journeys
- **Type Safety**: TypeScript for compile-time error catching

#### 9.2.2 Backend Testing
- **Unit Tests**: Go's built-in testing package
- **Integration Tests**: Database integration with test SQLite instances
- **API Tests**: HTTP endpoint testing with Go's httptest package
- **Test Data**: Fixtures and factories for consistent test data
- **Coverage**: Aim for >80% code coverage on critical paths

#### 9.2.3 Testing Configuration
```go
// Example Go test structure
func TestArticleHandler(t *testing.T) {
    // Setup test database
    db := setupTestDB()
    defer db.Close()
    
    // Test cases
    tests := []struct {
        name     string
        method   string
        path     string
        body     string
        expected int
    }{
        {"Create Article", "POST", "/api/articles", validArticleJSON, 201},
        {"Invalid Article", "POST", "/api/articles", invalidArticleJSON, 400},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 10. Development Environment & Deployment

### 10.1 Docker Development Setup

#### 10.1.1 Docker Compose Configuration
```yaml
version: '3.8'
services:
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile.dev
    ports:
      - "3000:3000"
    volumes:
      - ./frontend:/app
      - /app/node_modules
    environment:
      - VITE_API_URL=http://localhost:8080
    depends_on:
      - backend

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    volumes:
      - ./backend:/app
      - ./data:/app/data
    environment:
      - DB_PATH=/app/data/realworld.db
      - JWT_SECRET=your-development-secret
      - PORT=8080
    restart: unless-stopped
```

#### 10.1.2 Frontend Dockerfile
```dockerfile
# Development
FROM node:18-alpine as development
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
EXPOSE 3000
CMD ["npm", "run", "dev", "--", "--host"]

# Production build
FROM node:18-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build

# Production serve
FROM nginx:alpine as production
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

#### 10.1.3 Backend Dockerfile
```dockerfile
# Development
FROM golang:1.21-alpine as development
RUN go install github.com/cosmtrek/air@latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8080
CMD ["air"]

# Production build
FROM golang:1.21-alpine as build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Production
FROM alpine:latest as production
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /app/main .
EXPOSE 8080
CMD ["./main"]
```

### 10.2 Development Workflow
- **Quick Start**: `docker-compose up` for full environment
- **Hot Reload**: Both frontend and backend support live reloading
- **Database**: SQLite file persisted in `./data` directory
- **Environment Variables**: Managed through `.env` files
- **Port Configuration**: Frontend (3000), Backend (8080)

### 10.3 Production Deployment
- **Frontend**: Static build deployed to CDN (Vercel, Netlify)
- **Backend**: Containerized deployment (Docker, Railway, Fly.io)
- **Database**: SQLite file with backup strategies
- **Environment**: Production environment variables
- **Monitoring**: Health checks and logging integration

## 11. Success Metrics

### 11.1 Technical Metrics
- Application load time < 3 seconds
- 99.9% API uptime
- Zero critical security vulnerabilities
- Cross-browser compatibility

### 11.2 Learning Objectives
- Demonstrate full-stack development concepts
- Showcase framework interoperability
- Provide production-ready code examples
- Enable developer skill progression

## 12. Future Enhancements

### 12.1 Potential Features
- Rich text editor with WYSIWYG
- Image upload and management
- Advanced search functionality
- Email notifications
- Social media integration
- Article drafts and scheduling

### 12.2 Technical Improvements
- Real-time notifications
- Progressive Web App (PWA) features
- Advanced caching strategies
- Microservices architecture
- GraphQL API option

## 13. Project Structure & Development Guidelines

### 13.1 Recommended Project Structure
```
realworld-app/
├── docker-compose.yml
├── .env.example
├── README.md
│
├── frontend/
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   ├── components.json          # Shadcn/UI config
│   ├── Dockerfile.dev
│   ├── Dockerfile
│   │
│   ├── src/
│   │   ├── components/
│   │   │   ├── ui/              # Shadcn/UI components
│   │   │   ├── articles/
│   │   │   ├── auth/
│   │   │   └── layout/
│   │   │
│   │   ├── pages/
│   │   ├── hooks/
│   │   ├── lib/
│   │   │   ├── api.ts
│   │   │   ├── auth.ts
│   │   │   └── utils.ts
│   │   │
│   │   ├── types/
│   │   └── styles/
│   │       └── globals.css
│   │
│   └── public/
│
├── backend/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile.dev
│   ├── Dockerfile
│   ├── .air.toml              # Hot reload config
│   │
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── handlers/
│   │   │   ├── articles.go
│   │   │   ├── auth.go
│   │   │   ├── comments.go
│   │   │   ├── profiles.go
│   │   │   └── tags.go
│   │   │
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── logging.go
│   │   │
│   │   ├── models/
│   │   │   ├── article.go
│   │   │   ├── comment.go
│   │   │   └── user.go
│   │   │
│   │   ├── database/
│   │   │   ├── db.go
│   │   │   ├── migrations/
│   │   │   └── queries/
│   │   │
│   │   └── utils/
│   │       ├── jwt.go
│   │       ├── password.go
│   │       └── slugify.go
│   │
│   └── tests/
│       ├── integration/
│       └── fixtures/
│
└── data/                      # SQLite database files
    └── .gitkeep
```

### 13.2 Development Best Practices

#### 13.2.1 Code Organization
- **Separation of Concerns**: Clear boundaries between presentation, business logic, and data
- **Function-First**: Prefer pure functions over stateful classes
- **Explicit Dependencies**: Use dependency injection for testability
- **Error Handling**: Consistent error patterns across the application
- **Context Usage**: Leverage Go's context for request scoping and cancellation

#### 13.2.2 API Design Patterns
```go
// Example handler signature
type Handler struct {
    db     *sql.DB
    logger *log.Logger
}

func (h *Handler) CreateArticle(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Extract user from JWT context
    user, ok := ctx.Value(\"user\").(User)
    if !ok {
        http.Error(w, \"Unauthorized\", http.StatusUnauthorized)
        return
    }
    
    // Validate request
    var req CreateArticleRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, \"Invalid request\", http.StatusBadRequest)
        return
    }
    
    // Business logic
    article, err := h.createArticle(ctx, user.ID, req)
    if err != nil {
        h.logger.Printf(\"Error creating article: %v\", err)
        http.Error(w, \"Internal server error\", http.StatusInternalServerError)
        return
    }
    
    // Response
    w.Header().Set(\"Content-Type\", \"application/json\")
    json.NewEncoder(w).Encode(article)
}
```

#### 13.2.3 Database Patterns
```go
// Example database query pattern
func (h *Handler) GetArticleBySlug(ctx context.Context, slug string) (*Article, error) {
    query := `
        SELECT a.id, a.slug, a.title, a.description, a.body, 
               a.created_at, a.updated_at,
               u.username, u.bio, u.image,
               COUNT(f.user_id) as favorites_count
        FROM articles a
        JOIN users u ON a.author_id = u.id
        LEFT JOIN favorites f ON a.id = f.article_id
        WHERE a.slug = ?
        GROUP BY a.id
    `
    
    var article Article
    err := h.db.QueryRowContext(ctx, query, slug).Scan(
        &article.ID, &article.Slug, &article.Title,
        &article.Description, &article.Body,
        &article.CreatedAt, &article.UpdatedAt,
        &article.Author.Username, &article.Author.Bio, &article.Author.Image,
        &article.FavoritesCount,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrArticleNotFound
        }
        return nil, err
    }
    
    return &article, nil
}
```

#### 13.2.4 Frontend State Management
```typescript
// Example TanStack Query usage
export function useArticles(filters: ArticleFilters) {
  return useQuery({
    queryKey: ['articles', filters],
    queryFn: () => api.getArticles(filters),
    staleTime: 5 * 60 * 1000, // 5 minutes
    gcTime: 10 * 60 * 1000,   // 10 minutes
  });
}

export function useCreateArticle() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: api.createArticle,
    onSuccess: () => {
      // Invalidate and refetch articles
      queryClient.invalidateQueries({ queryKey: ['articles'] });
    },
  });
}
```

### 13.3 Environment Configuration

#### 13.3.1 Frontend Environment Variables
```bash
# .env.development
VITE_API_URL=http://localhost:8080
VITE_APP_TITLE=RealWorld
VITE_DEBUG=true

# .env.production  
VITE_API_URL=https://api.yourapp.com
VITE_APP_TITLE=RealWorld
VITE_DEBUG=false
```

#### 13.3.2 Backend Environment Variables
```bash
# .env
DB_PATH=./data/realworld.db
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRY=168h  # 7 days
PORT=8080
LOG_LEVEL=info
CORS_ORIGINS=http://localhost:3000
```

### 13.4 Performance Optimization

#### 13.4.1 Frontend Optimization
- **Code Splitting**: Route-based and component-based splitting
- **Tree Shaking**: Eliminate unused code with Vite
- **Asset Optimization**: Image compression and WebP support
- **Caching Strategy**: Aggressive caching with proper cache busting
- **Bundle Analysis**: Regular bundle size monitoring

#### 13.4.2 Backend Optimization
- **Connection Pooling**: SQLite with proper connection management
- **Query Optimization**: Indexed queries and efficient joins
- **Response Caching**: Cache frequently accessed data
- **Compression**: gzip compression for API responses
- **Rate Limiting**: Protect against abuse

## 14. Appendices

### 14.1 Reference Links
- Official Documentation: https://docs.realworld.show/
- GitHub Repository: https://github.com/gothinkster/realworld
- Demo API: https://api.realworld.io/api
- Community Discussions: GitHub Issues and Discussions

### 14.2 Implementation Examples
- 100+ implementations across different technology stacks
- Frontend: React, Angular, Vue.js, Svelte, etc.
- Backend: Node.js, Django, Rails, .NET, Go, etc.
- Database: PostgreSQL, MySQL, MongoDB, SQLite, etc.

### 14.3 Modern Technology Stack
Based on the recommendations from modern agentic coding practices:

#### 14.3.1 Recommended Stack
- **Frontend**: React + TypeScript + Vite + Tailwind CSS + Shadcn/UI + TanStack Query/Router
- **Backend**: Go + Standard Library HTTP + SQLite + JWT
- **Development**: Docker + Docker Compose + Hot Reload
- **Testing**: Vitest (Frontend) + Go Testing (Backend) + Playwright (E2E)
- **Deployment**: Containerized deployment with environment-specific configurations

#### 14.3.2 Key Architectural Decisions
- **Simplicity Over Complexity**: Choose the simplest solution that works
- **Function-First Design**: Prefer functions over complex object hierarchies
- **Plain SQL**: Direct database queries without ORM complexity
- **Explicit Context**: Clear request scoping and error handling
- **Agent-Friendly**: Code patterns optimized for AI assistance and maintenance

---

*This PRD serves as a comprehensive specification for building the RealWorld application, ensuring consistency across different technology implementations while maintaining the core Medium.com clone functionality.*