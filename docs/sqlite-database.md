# SQLite Database Documentation

## SQLite Database System

### Version Information
- **Current Stable**: SQLite 3.50.4 (2025-07-30)
- **Production Status**: Stable and production-ready
- **Features**: Query planner optimizations, jsonb_set() improvements
- **Testing**: 100% branch test coverage, rigorous testing standards
- **Compatibility**: Backward compatible with all 3.x versions

### Official Documentation
- **Main Documentation**: https://www.sqlite.org/docs.html
- **SQL Syntax**: https://www.sqlite.org/lang.html
- **C API Reference**: https://www.sqlite.org/c3ref/intro.html
- **SQL Features**: https://www.sqlite.org/features.html
- **Performance**: https://www.sqlite.org/speed.html
- **When to Use SQLite**: https://www.sqlite.org/whentouse.html
- **GitHub Mirror**: https://github.com/sqlite/sqlite

### Key Features for RealWorld
- **Serverless**: Zero-configuration database, no separate server process
- **Self-Contained**: Single file database, easy backup and deployment
- **Cross-Platform**: Runs on all major operating systems
- **ACID Compliant**: Full ACID (Atomicity, Consistency, Isolation, Durability) support
- **SQL Standard**: Supports most of SQL-92 standard
- **Lightweight**: Small memory footprint, suitable for embedded applications
- **Concurrent Access**: Multiple readers, single writer with WAL mode

### Database Schema for RealWorld

#### Complete Schema Implementation
```sql
-- RealWorld Database Schema
-- SQLite 3.50+ compatible

PRAGMA foreign_keys = ON;
PRAGMA journal_mode = WAL;
PRAGMA synchronous = NORMAL;
PRAGMA cache_size = 1000;
PRAGMA temp_store = memory;

-- Users table - Core user authentication and profile data
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(255) UNIQUE NOT NULL COLLATE NOCASE,
    email VARCHAR(255) UNIQUE NOT NULL COLLATE NOCASE,
    password_hash VARCHAR(255) NOT NULL,
    bio TEXT DEFAULT '',
    image VARCHAR(500) DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT username_length CHECK (length(username) >= 3 AND length(username) <= 50),
    CONSTRAINT email_format CHECK (email LIKE '%_@_%._%')
);

-- Articles table - Main content storage
CREATE TABLE articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    slug VARCHAR(255) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    body TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key relationships
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Constraints
    CONSTRAINT title_length CHECK (length(title) >= 1 AND length(title) <= 255),
    CONSTRAINT slug_format CHECK (slug GLOB '*[a-z0-9-]*')
);

-- Tags table - Article categorization
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(255) UNIQUE NOT NULL COLLATE NOCASE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT tag_name_length CHECK (length(name) >= 1 AND length(name) <= 50)
);

-- Article tags junction table - Many-to-many relationship
CREATE TABLE article_tags (
    article_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (article_id, tag_id),
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Comments table - Article comments
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    body TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    article_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    -- Foreign key relationships
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    
    -- Constraints
    CONSTRAINT comment_body_length CHECK (length(body) >= 1 AND length(body) <= 2000)
);

-- Favorites table - User article favorites (many-to-many)
CREATE TABLE favorites (
    user_id INTEGER NOT NULL,
    article_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (user_id, article_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE
);

-- Follows table - User following relationships (many-to-many)
CREATE TABLE follows (
    follower_id INTEGER NOT NULL,
    following_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY (follower_id, following_id),
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE,
    
    -- Prevent self-following
    CONSTRAINT no_self_follow CHECK (follower_id != following_id)
);

-- Performance indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_created_at ON users(created_at DESC);

CREATE INDEX idx_articles_slug ON articles(slug);
CREATE INDEX idx_articles_author_id ON articles(author_id);
CREATE INDEX idx_articles_created_at ON articles(created_at DESC);
CREATE INDEX idx_articles_updated_at ON articles(updated_at DESC);

CREATE INDEX idx_tags_name ON tags(name);

CREATE INDEX idx_article_tags_article_id ON article_tags(article_id);
CREATE INDEX idx_article_tags_tag_id ON article_tags(tag_id);

CREATE INDEX idx_comments_article_id ON comments(article_id);
CREATE INDEX idx_comments_author_id ON comments(author_id);
CREATE INDEX idx_comments_created_at ON comments(created_at DESC);

CREATE INDEX idx_favorites_user_id ON favorites(user_id);
CREATE INDEX idx_favorites_article_id ON favorites(article_id);

CREATE INDEX idx_follows_follower_id ON follows(follower_id);
CREATE INDEX idx_follows_following_id ON follows(following_id);

-- Triggers for updated_at timestamps
CREATE TRIGGER users_updated_at 
    AFTER UPDATE ON users
    FOR EACH ROW
    WHEN OLD.updated_at = NEW.updated_at
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER articles_updated_at 
    AFTER UPDATE ON articles
    FOR EACH ROW
    WHEN OLD.updated_at = NEW.updated_at
BEGIN
    UPDATE articles SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER comments_updated_at 
    AFTER UPDATE ON comments
    FOR EACH ROW
    WHEN OLD.updated_at = NEW.updated_at
BEGIN
    UPDATE comments SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Populate initial data
INSERT INTO users (username, email, password_hash, bio, image) VALUES
('demo', 'demo@realworld.io', '$2a$10$example_hash_here', 'Demo user account', ''),
('admin', 'admin@realworld.io', '$2a$10$example_hash_here', 'Administrator account', '');

INSERT INTO tags (name) VALUES
('welcome'),
('demo'),
('introduction'),
('getting-started');
```

### Go Database Integration

#### Database Connection and Setup
```go
// internal/database/db.go
package database

import (
    "database/sql"
    "embed"
    "fmt"
    "io/fs"
    "path/filepath"
    "sort"
    "strings"
    
    _ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type DB struct {
    *sql.DB
}

func New(dbPath string) (*DB, error) {
    // Connection string with optimizations
    connStr := fmt.Sprintf(
        "%s?_foreign_keys=on&_journal_mode=WAL&_synchronous=NORMAL&_cache_size=1000&_temp_store=memory&_timeout=5000",
        dbPath,
    )
    
    sqlDB, err := sql.Open("sqlite3", connStr)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    
    // Configure connection pool
    sqlDB.SetMaxOpenConns(25)
    sqlDB.SetMaxIdleConns(25)
    sqlDB.SetConnMaxLifetime(0) // SQLite doesn't need connection rotation
    
    // Test connection
    if err := sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    
    db := &DB{sqlDB}
    
    // Run migrations
    if err := db.migrate(); err != nil {
        return nil, fmt.Errorf("failed to migrate database: %w", err)
    }
    
    return db, nil
}

func (db *DB) migrate() error {
    // Create migrations table if it doesn't exist
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS migrations (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name VARCHAR(255) UNIQUE NOT NULL,
            executed_at DATETIME DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        return fmt.Errorf("failed to create migrations table: %w", err)
    }
    
    // Get list of migration files
    entries, err := fs.ReadDir(migrationFiles, "migrations")
    if err != nil {
        return fmt.Errorf("failed to read migration directory: %w", err)
    }
    
    var migrationNames []string
    for _, entry := range entries {
        if strings.HasSuffix(entry.Name(), ".sql") {
            migrationNames = append(migrationNames, entry.Name())
        }
    }
    sort.Strings(migrationNames)
    
    // Execute pending migrations
    for _, name := range migrationNames {
        var count int
        err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", name).Scan(&count)
        if err != nil {
            return fmt.Errorf("failed to check migration status: %w", err)
        }
        
        if count == 0 {
            // Read migration file
            content, err := fs.ReadFile(migrationFiles, filepath.Join("migrations", name))
            if err != nil {
                return fmt.Errorf("failed to read migration %s: %w", name, err)
            }
            
            // Execute migration
            _, err = db.Exec(string(content))
            if err != nil {
                return fmt.Errorf("failed to execute migration %s: %w", name, err)
            }
            
            // Record migration
            _, err = db.Exec("INSERT INTO migrations (name) VALUES (?)", name)
            if err != nil {
                return fmt.Errorf("failed to record migration %s: %w", name, err)
            }
            
            fmt.Printf("Executed migration: %s\n", name)
        }
    }
    
    return nil
}
```

#### Database Models and Queries
```go
// internal/database/queries.go
package database

import (
    "context"
    "database/sql"
    "fmt"
    "strings"
    "time"
)

// User operations
func (db *DB) CreateUser(ctx context.Context, username, email, passwordHash string) (int, error) {
    query := `
        INSERT INTO users (username, email, password_hash)
        VALUES (?, ?, ?)
    `
    
    result, err := db.ExecContext(ctx, query, username, email, passwordHash)
    if err != nil {
        return 0, fmt.Errorf("failed to create user: %w", err)
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to get user ID: %w", err)
    }
    
    return int(id), nil
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*User, error) {
    query := `
        SELECT id, username, email, password_hash, bio, image, created_at, updated_at
        FROM users
        WHERE email = ? COLLATE NOCASE
    `
    
    var user User
    err := db.QueryRowContext(ctx, query, email).Scan(
        &user.ID, &user.Username, &user.Email, &user.PasswordHash,
        &user.Bio, &user.Image, &user.CreatedAt, &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, ErrUserNotFound
        }
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }
    
    return &user, nil
}

// Article operations with complex queries
func (db *DB) GetArticlesWithFilters(ctx context.Context, filters ArticleFilters) ([]Article, int, error) {
    conditions := []string{"1=1"}
    args := []interface{}{}
    
    if filters.Author != "" {
        conditions = append(conditions, "u.username = ?")
        args = append(args, filters.Author)
    }
    
    if filters.Tag != "" {
        conditions = append(conditions, "a.id IN (SELECT at.article_id FROM article_tags at JOIN tags t ON at.tag_id = t.id WHERE t.name = ?)")
        args = append(args, filters.Tag)
    }
    
    if filters.FavoritedBy != "" {
        conditions = append(conditions, "a.id IN (SELECT f.article_id FROM favorites f JOIN users u2 ON f.user_id = u2.id WHERE u2.username = ?)")
        args = append(args, filters.FavoritedBy)
    }
    
    whereClause := strings.Join(conditions, " AND ")
    
    // Count query
    countQuery := fmt.Sprintf(`
        SELECT COUNT(DISTINCT a.id)
        FROM articles a
        JOIN users u ON a.author_id = u.id
        WHERE %s
    `, whereClause)
    
    var totalCount int
    err := db.QueryRowContext(ctx, countQuery, args...).Scan(&totalCount)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to count articles: %w", err)
    }
    
    // Main query with pagination
    query := fmt.Sprintf(`
        SELECT DISTINCT
            a.id, a.slug, a.title, a.description, a.body,
            a.created_at, a.updated_at,
            u.username, u.bio, u.image,
            (SELECT COUNT(*) FROM favorites f WHERE f.article_id = a.id) as favorites_count,
            (SELECT GROUP_CONCAT(t.name) FROM article_tags at JOIN tags t ON at.tag_id = t.id WHERE at.article_id = a.id) as tags
        FROM articles a
        JOIN users u ON a.author_id = u.id
        WHERE %s
        ORDER BY a.created_at DESC
        LIMIT ? OFFSET ?
    `, whereClause)
    
    args = append(args, filters.Limit, filters.Offset)
    
    rows, err := db.QueryContext(ctx, query, args...)
    if err != nil {
        return nil, 0, fmt.Errorf("failed to query articles: %w", err)
    }
    defer rows.Close()
    
    var articles []Article
    for rows.Next() {
        var article Article
        var tagString sql.NullString
        
        err := rows.Scan(
            &article.ID, &article.Slug, &article.Title, &article.Description, &article.Body,
            &article.CreatedAt, &article.UpdatedAt,
            &article.Author.Username, &article.Author.Bio, &article.Author.Image,
            &article.FavoritesCount, &tagString,
        )
        
        if err != nil {
            return nil, 0, fmt.Errorf("failed to scan article: %w", err)
        }
        
        // Parse tags
        if tagString.Valid && tagString.String != "" {
            article.TagList = strings.Split(tagString.String, ",")
        } else {
            article.TagList = []string{}
        }
        
        articles = append(articles, article)
    }
    
    if err = rows.Err(); err != nil {
        return nil, 0, fmt.Errorf("error iterating articles: %w", err)
    }
    
    return articles, totalCount, nil
}

// Transaction example for complex operations
func (db *DB) CreateArticleWithTags(ctx context.Context, article CreateArticleRequest, authorID int) (*Article, error) {
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to begin transaction: %w", err)
    }
    defer tx.Rollback()
    
    // Insert article
    articleQuery := `
        INSERT INTO articles (slug, title, description, body, author_id)
        VALUES (?, ?, ?, ?, ?)
    `
    
    result, err := tx.ExecContext(ctx, articleQuery, 
        article.Slug, article.Title, article.Description, article.Body, authorID)
    if err != nil {
        return nil, fmt.Errorf("failed to create article: %w", err)
    }
    
    articleID, err := result.LastInsertId()
    if err != nil {
        return nil, fmt.Errorf("failed to get article ID: %w", err)
    }
    
    // Insert tags and create associations
    for _, tagName := range article.TagList {
        // Insert or get tag ID
        var tagID int
        err := tx.QueryRowContext(ctx, 
            "INSERT OR IGNORE INTO tags (name) VALUES (?) RETURNING id", 
            tagName).Scan(&tagID)
        
        if err != nil {
            // Fallback for older SQLite versions
            err = tx.QueryRowContext(ctx, 
                "SELECT id FROM tags WHERE name = ?", 
                tagName).Scan(&tagID)
            if err != nil {
                return nil, fmt.Errorf("failed to get tag ID: %w", err)
            }
        }
        
        // Create article-tag association
        _, err = tx.ExecContext(ctx, 
            "INSERT INTO article_tags (article_id, tag_id) VALUES (?, ?)",
            articleID, tagID)
        if err != nil {
            return nil, fmt.Errorf("failed to create article-tag association: %w", err)
        }
    }
    
    // Commit transaction
    if err = tx.Commit(); err != nil {
        return nil, fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    // Return created article
    return db.GetArticleBySlug(ctx, article.Slug)
}
```

### Performance Optimization

#### Query Optimization Techniques
```sql
-- Efficient article feed query with user following
SELECT DISTINCT a.id, a.slug, a.title, a.description, 
       a.created_at, u.username, u.image,
       COUNT(f.user_id) as favorites_count
FROM articles a
JOIN users u ON a.author_id = u.id
JOIN follows fo ON a.author_id = fo.following_id
LEFT JOIN favorites f ON a.id = f.article_id
WHERE fo.follower_id = ?
GROUP BY a.id
ORDER BY a.created_at DESC
LIMIT ? OFFSET ?;

-- Optimized tag popularity query
SELECT t.name, COUNT(at.article_id) as article_count
FROM tags t
JOIN article_tags at ON t.id = at.tag_id
JOIN articles a ON at.article_id = a.id
WHERE a.created_at > datetime('now', '-30 days')
GROUP BY t.id, t.name
ORDER BY article_count DESC
LIMIT 10;

-- User activity summary with single query
SELECT 
    u.username,
    COUNT(DISTINCT a.id) as article_count,
    COUNT(DISTINCT c.id) as comment_count,
    COUNT(DISTINCT f.article_id) as favorites_count
FROM users u
LEFT JOIN articles a ON u.id = a.author_id
LEFT JOIN comments c ON u.id = c.author_id
LEFT JOIN favorites f ON u.id = f.user_id
WHERE u.id = ?
GROUP BY u.id, u.username;
```

#### Database Configuration for Production
```go
// Production database configuration
func configureProduction(db *sql.DB) error {
    pragmas := []string{
        "PRAGMA foreign_keys = ON",
        "PRAGMA journal_mode = WAL",
        "PRAGMA synchronous = NORMAL",
        "PRAGMA cache_size = -64000", // 64MB cache
        "PRAGMA temp_store = memory",
        "PRAGMA mmap_size = 268435456", // 256MB mmap
        "PRAGMA optimize",
    }
    
    for _, pragma := range pragmas {
        if _, err := db.Exec(pragma); err != nil {
            return fmt.Errorf("failed to execute %s: %w", pragma, err)
        }
    }
    
    return nil
}
```

### Backup and Maintenance

#### Backup Strategy
```go
// Database backup functionality
func (db *DB) Backup(ctx context.Context, backupPath string) error {
    query := fmt.Sprintf("VACUUM INTO '%s'", backupPath)
    _, err := db.ExecContext(ctx, query)
    if err != nil {
        return fmt.Errorf("failed to backup database: %w", err)
    }
    
    return nil
}

// Database maintenance
func (db *DB) Maintenance(ctx context.Context) error {
    queries := []string{
        "PRAGMA optimize",
        "PRAGMA wal_checkpoint(TRUNCATE)",
        "ANALYZE",
    }
    
    for _, query := range queries {
        if _, err := db.ExecContext(ctx, query); err != nil {
            return fmt.Errorf("maintenance query failed: %s: %w", query, err)
        }
    }
    
    return nil
}
```

This SQLite implementation provides a robust, performant database foundation for the RealWorld application with proper indexing, constraints, and Go integration patterns.