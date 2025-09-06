package database

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

type DB struct {
	*sql.DB
}

func New(dbPath string) (*DB, error) {
	// Connection string with optimizations as per documentation
	connStr := fmt.Sprintf(
		"%s?_foreign_keys=on&_journal_mode=WAL&_synchronous=NORMAL&_cache_size=1000&_temp_store=memory&_timeout=5000",
		dbPath,
	)

	sqlDB, err := sql.Open("sqlite3", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool for production use
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

	// Configure production optimizations
	if err := db.configureProduction(); err != nil {
		return nil, fmt.Errorf("failed to configure database: %w", err)
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

			// Execute migration in transaction
			tx, err := db.Begin()
			if err != nil {
				return fmt.Errorf("failed to begin transaction for migration %s: %w", name, err)
			}

			_, err = tx.Exec(string(content))
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to execute migration %s: %w", name, err)
			}

			// Record migration
			_, err = tx.Exec("INSERT INTO migrations (name) VALUES (?)", name)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("failed to record migration %s: %w", name, err)
			}

			if err = tx.Commit(); err != nil {
				return fmt.Errorf("failed to commit migration %s: %w", name, err)
			}

			fmt.Printf("Executed migration: %s\n", name)
		}
	}

	return nil
}

func (db *DB) configureProduction() error {
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

// Backup creates a backup of the database
func (db *DB) Backup(backupPath string) error {
	query := fmt.Sprintf("VACUUM INTO '%s'", backupPath)
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to backup database: %w", err)
	}
	return nil
}

// Maintenance performs database maintenance tasks
func (db *DB) Maintenance() error {
	queries := []string{
		"PRAGMA optimize",
		"PRAGMA wal_checkpoint(TRUNCATE)",
		"ANALYZE",
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("maintenance query failed: %s: %w", query, err)
		}
	}

	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}