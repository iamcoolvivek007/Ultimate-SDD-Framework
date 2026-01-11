package db

import (
	"database/sql"
	"time"
)

// Message represents a chat message
type Message struct {
	ID          string    `json:"id"`
	SessionID   string    `json:"session_id"`
	Role        string    `json:"role"` // "user", "assistant", "system", "tool"
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	TokenCount  int       `json:"token_count"`
	Model       string    `json:"model"`
	ToolCalls   string    `json:"tool_calls,omitempty"`   // JSON array of tool calls
	ToolResults string    `json:"tool_results,omitempty"` // JSON array of tool results
}

// MessageStore handles message CRUD operations
type MessageStore struct {
	db *DB
}

// NewMessageStore creates a new message store
func NewMessageStore(db *DB) *MessageStore {
	return &MessageStore{db: db}
}

// Create creates a new message
func (m *MessageStore) Create(msg *Message) error {
	if msg.ID == "" {
		msg.ID = GenerateID("msg")
	}
	msg.CreatedAt = time.Now()

	_, err := m.db.conn.Exec(`
		INSERT INTO messages (id, session_id, role, content, created_at, token_count, model, tool_calls, tool_results)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, msg.ID, msg.SessionID, msg.Role, msg.Content, msg.CreatedAt, msg.TokenCount, msg.Model, msg.ToolCalls, msg.ToolResults)

	return err
}

// GetByID retrieves a message by ID
func (m *MessageStore) GetByID(id string) (*Message, error) {
	msg := &Message{}
	err := m.db.conn.QueryRow(`
		SELECT id, session_id, role, content, created_at, token_count, model, tool_calls, tool_results
		FROM messages WHERE id = ?
	`, id).Scan(&msg.ID, &msg.SessionID, &msg.Role, &msg.Content, &msg.CreatedAt, &msg.TokenCount, &msg.Model, &msg.ToolCalls, &msg.ToolResults)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return msg, err
}

// ListBySession retrieves all messages for a session
func (m *MessageStore) ListBySession(sessionID string, limit int) ([]*Message, error) {
	var rows *sql.Rows
	var err error

	if limit > 0 {
		rows, err = m.db.conn.Query(`
			SELECT id, session_id, role, content, created_at, token_count, model, tool_calls, tool_results
			FROM messages WHERE session_id = ? ORDER BY created_at ASC LIMIT ?
		`, sessionID, limit)
	} else {
		rows, err = m.db.conn.Query(`
			SELECT id, session_id, role, content, created_at, token_count, model, tool_calls, tool_results
			FROM messages WHERE session_id = ? ORDER BY created_at ASC
		`, sessionID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		if err := rows.Scan(&msg.ID, &msg.SessionID, &msg.Role, &msg.Content, &msg.CreatedAt, &msg.TokenCount, &msg.Model, &msg.ToolCalls, &msg.ToolResults); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

// GetLastN retrieves the last N messages from a session
func (m *MessageStore) GetLastN(sessionID string, n int) ([]*Message, error) {
	rows, err := m.db.conn.Query(`
		SELECT id, session_id, role, content, created_at, token_count, model, tool_calls, tool_results
		FROM messages WHERE session_id = ? ORDER BY created_at DESC LIMIT ?
	`, sessionID, n)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		if err := rows.Scan(&msg.ID, &msg.SessionID, &msg.Role, &msg.Content, &msg.CreatedAt, &msg.TokenCount, &msg.Model, &msg.ToolCalls, &msg.ToolResults); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	// Reverse to get chronological order
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, rows.Err()
}

// CountBySession returns the number of messages in a session
func (m *MessageStore) CountBySession(sessionID string) (int, error) {
	var count int
	err := m.db.conn.QueryRow(`SELECT COUNT(*) FROM messages WHERE session_id = ?`, sessionID).Scan(&count)
	return count, err
}

// TotalTokensBySession returns the total token count for a session
func (m *MessageStore) TotalTokensBySession(sessionID string) (int, error) {
	var total int
	err := m.db.conn.QueryRow(`SELECT COALESCE(SUM(token_count), 0) FROM messages WHERE session_id = ?`, sessionID).Scan(&total)
	return total, err
}

// DeleteBySession deletes all messages for a session
func (m *MessageStore) DeleteBySession(sessionID string) error {
	_, err := m.db.conn.Exec(`DELETE FROM messages WHERE session_id = ?`, sessionID)
	return err
}

// Search searches messages by content
func (m *MessageStore) Search(query string, limit int) ([]*Message, error) {
	rows, err := m.db.conn.Query(`
		SELECT id, session_id, role, content, created_at, token_count, model, tool_calls, tool_results
		FROM messages WHERE content LIKE ? ORDER BY created_at DESC LIMIT ?
	`, "%"+query+"%", limit)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		if err := rows.Scan(&msg.ID, &msg.SessionID, &msg.Role, &msg.Content, &msg.CreatedAt, &msg.TokenCount, &msg.Model, &msg.ToolCalls, &msg.ToolResults); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}
