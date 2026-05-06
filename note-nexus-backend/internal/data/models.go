package data

import (
	"database/sql"
	"errors"
)

var (
	// ErrRecordNotFound is a custom error we can check against in our handlers
	ErrRecordNotFound = errors.New("record not found")
)

// Models wraps all our individual database models.
// This allows us to pass a single struct to our handlers containing everything we need.
type Models struct {
	Tasks         TaskModel
	Notes         NoteModel
	Users         UserModel
	Workspaces    WorkspaceModel
	Reminders     ReminderModel
	Notifications NotificationModel
}

// NewModels returns a Models struct containing the initialized models.
func NewModels(db *sql.DB) Models {
	return Models{
		Tasks:         TaskModel{DB: db},
		Notes:         NoteModel{DB: db},
		Users:         UserModel{DB: db},
		Workspaces:    WorkspaceModel{DB: db},
		Reminders:     ReminderModel{DB: db},
		Notifications: NotificationModel{DB: db},
	}
}
