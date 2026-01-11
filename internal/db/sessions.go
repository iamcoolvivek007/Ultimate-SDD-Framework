package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Session represents a chat session
type Session struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Model        string    `json:"model"`
	Provider     string    `json:"provider"`
	ProjectPath  string    `json:"project_path"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Summary      string    `json:"summary"`
	IsActive     bool      `json:"is_active"`
	TokenCount   int       `json:"token_count"`
	MessageCount int       `json:"message_count"`
}

// SessionStore handles session CRUD operations
type SessionStore struct {
	db *DB
}

// NewSessionStore creates a new session store
func NewSessionStore(db *DB) *SessionStore {
	return &SessionStore{db: db}
}

// Create creates a new session
func (s *SessionStore) Create(session *Session) error {
	if session.ID == "" {
		session.ID = GenerateID("sess")
	}
	if session.Title == "" {
		session.Title = "New Session"
	}
	session.CreatedAt = time.Now()
	session.UpdatedAt = time.Now()

	_, err := s.db.conn.Exec(`
		INSERT INTO sessions (id, title, model, provider, project_path, created_at, updated_at, summary, is_active, token_count, message_count)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, session.ID, session.Title, session.Model, session.Provider, session.ProjectPath,
		session.CreatedAt, session.UpdatedAt, session.Summary, session.IsActive, session.TokenCount, session.MessageCount)

	return err
}

// GetByID retrieves a session by ID
func (s *SessionStore) GetByID(id string) (*Session, error) {
	session := &Session{}
	err := s.db.conn.QueryRow(`
		SELECT id, title, model, provider, project_path, created_at, updated_at, summary, is_active, token_count, message_count
		FROM sessions WHERE id = ?
	`, id).Scan(&session.ID, &session.Title, &session.Model, &session.Provider, &session.ProjectPath,
		&session.CreatedAt, &session.UpdatedAt, &session.Summary, &session.IsActive, &session.TokenCount, &session.MessageCount)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

// GetActive retrieves the currently active session
func (s *SessionStore) GetActive() (*Session, error) {
	session := &Session{}
	err := s.db.conn.QueryRow(`
		SELECT id, title, model, provider, project_path, created_at, updated_at, summary, is_active, token_count, message_count
		FROM sessions WHERE is_active = 1 ORDER BY updated_at DESC LIMIT 1
	`).Scan(&session.ID, &session.Title, &session.Model, &session.Provider, &session.ProjectPath,
		&session.CreatedAt, &session.UpdatedAt, &session.Summary, &session.IsActive, &session.TokenCount, &session.MessageCount)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return session, err
}

// SetActive sets a session as active (and deactivates others)
func (s *SessionStore) SetActive(id string) error {
	tx, err := s.db.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Deactivate all sessions
	if _, err := tx.Exec(`UPDATE sessions SET is_active = 0`); err != nil {
		return err
	}

	// Activate the specified session
	if _, err := tx.Exec(`UPDATE sessions SET is_active = 1, updated_at = ? WHERE id = ?`, time.Now(), id); err != nil {
		return err
	}

	return tx.Commit()
}

// List retrieves all sessions, optionally filtered by project
func (s *SessionStore) List(projectPath string, limit int) ([]*Session, error) {
	var rows *sql.Rows
	var err error

	if projectPath != "" {
		rows, err = s.db.conn.Query(`
			SELECT id, title, model, provider, project_path, created_at, updated_at, summary, is_active, token_count, message_count
			FROM sessions WHERE project_path = ? ORDER BY updated_at DESC LIMIT ?
		`, projectPath, limit)
	} else {
		rows, err = s.db.conn.Query(`
			SELECT id, title, model, provider, project_path, created_at, updated_at, summary, is_active, token_count, message_count
			FROM sessions ORDER BY updated_at DESC LIMIT ?
		`, limit)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*Session
	for rows.Next() {
		session := &Session{}
		if err := rows.Scan(&session.ID, &session.Title, &session.Model, &session.Provider, &session.ProjectPath,
			&session.CreatedAt, &session.UpdatedAt, &session.Summary, &session.IsActive, &session.TokenCount, &session.MessageCount); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}

	return sessions, rows.Err()
}

// Update updates a session
func (s *SessionStore) Update(session *Session) error {
	session.UpdatedAt = time.Now()
	_, err := s.db.conn.Exec(`
		UPDATE sessions SET title = ?, model = ?, provider = ?, summary = ?, token_count = ?, message_count = ?, updated_at = ?
		WHERE id = ?
	`, session.Title, session.Model, session.Provider, session.Summary, session.TokenCount, session.MessageCount, session.UpdatedAt, session.ID)
	return err
}

// Delete deletes a session and its messages
func (s *SessionStore) Delete(id string) error {
	_, err := s.db.conn.Exec(`DELETE FROM sessions WHERE id = ?`, id)
	return err
}

// UpdateTitle updates the session title
func (s *SessionStore) UpdateTitle(id, title string) error {
	_, err := s.db.conn.Exec(`UPDATE sessions SET title = ?, updated_at = ? WHERE id = ?`, title, time.Now(), id)
	return err
}

// IncrementTokenCount increments the token count for a session
func (s *SessionStore) IncrementTokenCount(id string, tokens int) error {
	_, err := s.db.conn.Exec(`UPDATE sessions SET token_count = token_count + ?, message_count = message_count + 1, updated_at = ? WHERE id = ?`,
		tokens, time.Now(), id)
	return err
}

// Compact summarizes and creates a new session from an existing one
func (s *SessionStore) Compact(oldSessionID string, summary string) (*Session, error) {
	oldSession, err := s.GetByID(oldSessionID)
	if err != nil || oldSession == nil {
		return nil, fmt.Errorf("session not found: %s", oldSessionID)
	}

	// Create new session with summary as context
	newSession := &Session{
		Title:       oldSession.Title + " (continued)",
		Model:       oldSession.Model,
		Provider:    oldSession.Provider,
		ProjectPath: oldSession.ProjectPath,
		Summary:     summary,
		IsActive:    true,
	}

	if err := s.Create(newSession); err != nil {
		return nil, err
	}

	// Deactivate old session
	if _, err := s.db.conn.Exec(`UPDATE sessions SET is_active = 0 WHERE id = ?`, oldSessionID); err != nil {
		return nil, err
	}

	return newSession, nil
}

// ExportToJSON exports a session and its messages to JSON
func (s *SessionStore) ExportToJSON(id string) ([]byte, error) {
	session, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	msgStore := NewMessageStore(s.db)
	messages, err := msgStore.ListBySession(id, 0)
	if err != nil {
		return nil, err
	}

	export := struct {
		Session  *Session   `json:"session"`
		Messages []*Message `json:"messages"`
	}{
		Session:  session,
		Messages: messages,
	}

	return json.MarshalIndent(export, "", "  ")
}
