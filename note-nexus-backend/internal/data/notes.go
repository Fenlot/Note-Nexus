package data

import (
	"database/sql"
	"strings"
	"time"
)

type Note struct {
	ID          int64     `json:"id"`
	WorkspaceID int64     `json:"workspaceId"`
	UserID      int64     `json:"userId"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type NoteModel struct {
	DB *sql.DB
}

// Insert creates a new note.
func (m NoteModel) Insert(content string, workspaceID, userID int64) (int64, error) {
	title := "Untitled Note"
	lines := strings.Split(content, "\n")
	if len(lines) > 0 && strings.TrimSpace(lines[0]) != "" {
		title = strings.TrimSpace(lines[0])
		if len(title) > 50 {
			title = title[:50] + "..."
		}
	}

	stmt := `INSERT INTO notes (workspace_id, user_id, title, content, updated_at) VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) RETURNING id`

	var id int64
	err := m.DB.QueryRow(stmt, workspaceID, userID, title, content).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetAll retrieves all notes for a workspace.
func (m NoteModel) GetAll(workspaceID int64) ([]Note, error) {
	stmt := `SELECT id, workspace_id, user_id, title, content, updated_at FROM notes WHERE workspace_id = $1 ORDER BY updated_at DESC`
	rows, err := m.DB.Query(stmt, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []Note{}
	for rows.Next() {
		var n Note
		err := rows.Scan(&n.ID, &n.WorkspaceID, &n.UserID, &n.Title, &n.Content, &n.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}

	return notes, nil
}

// Update modifies an existing note.
func (m NoteModel) Update(note *Note) error {
	lines := strings.Split(note.Content, "\n")
	if len(lines) > 0 && strings.TrimSpace(lines[0]) != "" {
		note.Title = strings.TrimSpace(lines[0])
		if len(note.Title) > 50 {
			note.Title = note.Title[:50] + "..."
		}
	}

	stmt := `UPDATE notes SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 AND workspace_id = $4`
	_, err := m.DB.Exec(stmt, note.Title, note.Content, note.ID, note.WorkspaceID)
	return err
}

// Delete removes a note by ID.
func (m NoteModel) Delete(id, workspaceID int64) error {
	stmt := `DELETE FROM notes WHERE id = $1 AND workspace_id = $2`
	_, err := m.DB.Exec(stmt, id, workspaceID)
	return err
}

// Get retrieves a single note.
func (m NoteModel) Get(id, workspaceID int64) (*Note, error) {
	stmt := `SELECT id, workspace_id, user_id, title, content, updated_at FROM notes WHERE id = $1 AND workspace_id = $2`

	var n Note
	err := m.DB.QueryRow(stmt, id, workspaceID).Scan(&n.ID, &n.WorkspaceID, &n.UserID, &n.Title, &n.Content, &n.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &n, nil
}