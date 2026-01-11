package db

import (
	"database/sql"
	"time"
)

// FileChange represents a file modification tracked by the system
type FileChange struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id"`
	MessageID     string    `json:"message_id,omitempty"`
	FilePath      string    `json:"file_path"`
	Operation     string    `json:"operation"` // "create", "modify", "delete"
	BeforeContent string    `json:"before_content,omitempty"`
	AfterContent  string    `json:"after_content,omitempty"`
	Diff          string    `json:"diff,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	Reverted      bool      `json:"reverted"`
}

// FileChangeStore handles file change tracking
type FileChangeStore struct {
	db *DB
}

// NewFileChangeStore creates a new file change store
func NewFileChangeStore(db *DB) *FileChangeStore {
	return &FileChangeStore{db: db}
}

// Create records a new file change
func (f *FileChangeStore) Create(change *FileChange) error {
	if change.ID == "" {
		change.ID = GenerateID("fc")
	}
	change.CreatedAt = time.Now()

	_, err := f.db.conn.Exec(`
		INSERT INTO file_changes (id, session_id, message_id, file_path, operation, before_content, after_content, diff, created_at, reverted)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, change.ID, change.SessionID, change.MessageID, change.FilePath, change.Operation,
		change.BeforeContent, change.AfterContent, change.Diff, change.CreatedAt, change.Reverted)

	return err
}

// GetByID retrieves a file change by ID
func (f *FileChangeStore) GetByID(id string) (*FileChange, error) {
	change := &FileChange{}
	err := f.db.conn.QueryRow(`
		SELECT id, session_id, message_id, file_path, operation, before_content, after_content, diff, created_at, reverted
		FROM file_changes WHERE id = ?
	`, id).Scan(&change.ID, &change.SessionID, &change.MessageID, &change.FilePath, &change.Operation,
		&change.BeforeContent, &change.AfterContent, &change.Diff, &change.CreatedAt, &change.Reverted)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return change, err
}

// ListBySession retrieves all file changes for a session
func (f *FileChangeStore) ListBySession(sessionID string) ([]*FileChange, error) {
	rows, err := f.db.conn.Query(`
		SELECT id, session_id, message_id, file_path, operation, before_content, after_content, diff, created_at, reverted
		FROM file_changes WHERE session_id = ? ORDER BY created_at DESC
	`, sessionID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changes []*FileChange
	for rows.Next() {
		change := &FileChange{}
		if err := rows.Scan(&change.ID, &change.SessionID, &change.MessageID, &change.FilePath, &change.Operation,
			&change.BeforeContent, &change.AfterContent, &change.Diff, &change.CreatedAt, &change.Reverted); err != nil {
			return nil, err
		}
		changes = append(changes, change)
	}

	return changes, rows.Err()
}

// ListByFile retrieves all changes for a specific file
func (f *FileChangeStore) ListByFile(filePath string, limit int) ([]*FileChange, error) {
	rows, err := f.db.conn.Query(`
		SELECT id, session_id, message_id, file_path, operation, before_content, after_content, diff, created_at, reverted
		FROM file_changes WHERE file_path = ? ORDER BY created_at DESC LIMIT ?
	`, filePath, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changes []*FileChange
	for rows.Next() {
		change := &FileChange{}
		if err := rows.Scan(&change.ID, &change.SessionID, &change.MessageID, &change.FilePath, &change.Operation,
			&change.BeforeContent, &change.AfterContent, &change.Diff, &change.CreatedAt, &change.Reverted); err != nil {
			return nil, err
		}
		changes = append(changes, change)
	}

	return changes, rows.Err()
}

// GetLastForFile retrieves the most recent change for a file that hasn't been reverted
func (f *FileChangeStore) GetLastForFile(filePath string) (*FileChange, error) {
	change := &FileChange{}
	err := f.db.conn.QueryRow(`
		SELECT id, session_id, message_id, file_path, operation, before_content, after_content, diff, created_at, reverted
		FROM file_changes WHERE file_path = ? AND reverted = 0 ORDER BY created_at DESC LIMIT 1
	`, filePath).Scan(&change.ID, &change.SessionID, &change.MessageID, &change.FilePath, &change.Operation,
		&change.BeforeContent, &change.AfterContent, &change.Diff, &change.CreatedAt, &change.Reverted)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return change, err
}

// MarkReverted marks a file change as reverted
func (f *FileChangeStore) MarkReverted(id string) error {
	_, err := f.db.conn.Exec(`UPDATE file_changes SET reverted = 1 WHERE id = ?`, id)
	return err
}

// GetPendingReverts returns changes that can be reverted (not already reverted)
func (f *FileChangeStore) GetPendingReverts(sessionID string, limit int) ([]*FileChange, error) {
	rows, err := f.db.conn.Query(`
		SELECT id, session_id, message_id, file_path, operation, before_content, after_content, diff, created_at, reverted
		FROM file_changes WHERE session_id = ? AND reverted = 0 ORDER BY created_at DESC LIMIT ?
	`, sessionID, limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changes []*FileChange
	for rows.Next() {
		change := &FileChange{}
		if err := rows.Scan(&change.ID, &change.SessionID, &change.MessageID, &change.FilePath, &change.Operation,
			&change.BeforeContent, &change.AfterContent, &change.Diff, &change.CreatedAt, &change.Reverted); err != nil {
			return nil, err
		}
		changes = append(changes, change)
	}

	return changes, rows.Err()
}

// CountBySession counts file changes in a session
func (f *FileChangeStore) CountBySession(sessionID string) (int, error) {
	var count int
	err := f.db.conn.QueryRow(`SELECT COUNT(*) FROM file_changes WHERE session_id = ?`, sessionID).Scan(&count)
	return count, err
}

// GetAffectedFiles returns unique file paths affected in a session
func (f *FileChangeStore) GetAffectedFiles(sessionID string) ([]string, error) {
	rows, err := f.db.conn.Query(`
		SELECT DISTINCT file_path FROM file_changes WHERE session_id = ? ORDER BY file_path
	`, sessionID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []string
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			return nil, err
		}
		files = append(files, path)
	}

	return files, rows.Err()
}
