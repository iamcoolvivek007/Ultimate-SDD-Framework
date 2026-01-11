package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the SQLite database connection
type DB struct {
	conn *sql.DB
	path string
}

// Config holds database configuration
type Config struct {
	Path string // Path to SQLite database file
}

// DefaultConfig returns default database configuration
func DefaultConfig() Config {
	homeDir, _ := os.UserHomeDir()
	return Config{
		Path: filepath.Join(homeDir, ".viki", "viki.db"),
	}
}

// ProjectConfig returns project-specific database configuration
func ProjectConfig(projectDir string) Config {
	return Config{
		Path: filepath.Join(projectDir, ".sdd", "viki.db"),
	}
}

// New creates a new database connection
func New(cfg Config) (*DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	conn, err := sql.Open("sqlite3", cfg.Path+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{
		conn: conn,
		path: cfg.Path,
	}

	// Run migrations
	if err := db.migrate(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// Conn returns the underlying sql.DB connection
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// migrate runs database migrations
func (db *DB) migrate() error {
	migrations := []string{
		// Sessions table
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL DEFAULT 'New Session',
			model TEXT,
			provider TEXT,
			project_path TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			summary TEXT,
			is_active INTEGER DEFAULT 0,
			token_count INTEGER DEFAULT 0,
			message_count INTEGER DEFAULT 0
		)`,

		// Messages table
		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL,
			role TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			token_count INTEGER DEFAULT 0,
			model TEXT,
			tool_calls TEXT,
			tool_results TEXT,
			FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
		)`,

		// File changes table
		`CREATE TABLE IF NOT EXISTS file_changes (
			id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL,
			message_id TEXT,
			file_path TEXT NOT NULL,
			operation TEXT NOT NULL,
			before_content TEXT,
			after_content TEXT,
			diff TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			reverted INTEGER DEFAULT 0,
			FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
			FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE SET NULL
		)`,

		// Tool executions table
		`CREATE TABLE IF NOT EXISTS tool_executions (
			id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL,
			message_id TEXT,
			tool_name TEXT NOT NULL,
			tool_input TEXT,
			tool_output TEXT,
			status TEXT NOT NULL,
			duration_ms INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE,
			FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE SET NULL
		)`,

		// Custom commands table
		`CREATE TABLE IF NOT EXISTS custom_commands (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			description TEXT,
			content TEXT NOT NULL,
			scope TEXT NOT NULL DEFAULT 'user',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Learning records table
		`CREATE TABLE IF NOT EXISTS learning_records (
			id TEXT PRIMARY KEY,
			type TEXT NOT NULL,
			context TEXT,
			action TEXT NOT NULL,
			outcome TEXT,
			success INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Indexes
		`CREATE INDEX IF NOT EXISTS idx_messages_session ON messages(session_id)`,
		`CREATE INDEX IF NOT EXISTS idx_file_changes_session ON file_changes(session_id)`,
		`CREATE INDEX IF NOT EXISTS idx_tool_executions_session ON tool_executions(session_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_project ON sessions(project_path)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_active ON sessions(is_active)`,
	}

	for _, migration := range migrations {
		if _, err := db.conn.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

// GenerateID generates a unique ID with prefix
func GenerateID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}
