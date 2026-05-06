package data

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          int64     `json:"id"`
	WorkspaceID int64     `json:"workspaceId"`
	UserID      int64     `json:"userId"`
	Title       string    `json:"title"`
	IsCompleted bool      `json:"is_completed"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskModel struct {
	DB *sql.DB
}

func (m TaskModel) Insert(title string, workspaceID, userID int64) (int64, error) {
	stmt := `INSERT INTO tasks (workspace_id, user_id, title, is_completed, created_at) 
	VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) RETURNING id`

	var id int64
	err := m.DB.QueryRow(stmt, workspaceID, userID, title, false).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m TaskModel) GetAll(workspaceID int64) ([]Task, error) {
	stmt := `SELECT id, workspace_id, user_id, title, is_completed, created_at FROM tasks WHERE workspace_id = $1 ORDER BY id DESC`

	rows, err := m.DB.Query(stmt, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []Task{}

	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.WorkspaceID, &t.UserID, &t.Title, &t.IsCompleted, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (m TaskModel) Update(task *Task) error {
	stmt := `UPDATE tasks SET title = $1, is_completed = $2 WHERE id = $3 AND workspace_id = $4`
	_, err := m.DB.Exec(stmt, task.Title, task.IsCompleted, task.ID, task.WorkspaceID)
	return err
}

func (m TaskModel) Get(id, workspaceID int64) (*Task, error) {
	stmt := `SELECT id, workspace_id, user_id, title, is_completed, created_at FROM tasks WHERE id = $1 AND workspace_id = $2`

	var t Task
	err := m.DB.QueryRow(stmt, id, workspaceID).Scan(&t.ID, &t.WorkspaceID, &t.UserID, &t.Title, &t.IsCompleted, &t.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (m TaskModel) Delete(id, workspaceID int64) error {
	stmt := `DELETE FROM tasks WHERE id = $1 AND workspace_id = $2`
	_, err := m.DB.Exec(stmt, id, workspaceID)
	return err
}