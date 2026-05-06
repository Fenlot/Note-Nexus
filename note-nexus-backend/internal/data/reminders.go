package data

import (
	"database/sql"
	"time"
)

// Reminder represents a reminder for a task, note, or custom event
type Reminder struct {
	ID           int64     `json:"id"`
	WorkspaceID  int64     `json:"workspace_id"`
	UserID       int64     `json:"user_id"`
	TargetID     *int64    `json:"target_id"`     // ID of task or note, nullable for custom reminders
	TargetType   string    `json:"target_type"`   // 'task', 'note', 'custom'
	ReminderType string    `json:"reminder_type"` // 'task_due', 'note_review', 'custom'
	Title        string    `json:"title"`
	Description  *string   `json:"description"`
	DueDate      time.Time `json:"due_date"`
	ScheduleType string    `json:"schedule_type"` // 'once', 'daily', 'weekly', 'monthly'
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ReminderModel handles database operations for reminders
type ReminderModel struct {
	DB *sql.DB
}

// Create inserts a new reminder
func (m ReminderModel) Create(reminder *Reminder) error {
	query := `
		INSERT INTO reminders (workspace_id, user_id, target_id, target_type, reminder_type, title, description, due_date, schedule_type, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at
	`

	err := m.DB.QueryRow(query, reminder.WorkspaceID, reminder.UserID, reminder.TargetID, reminder.TargetType,
		reminder.ReminderType, reminder.Title, reminder.Description, reminder.DueDate, reminder.ScheduleType, reminder.IsActive).
		Scan(&reminder.ID, &reminder.CreatedAt, &reminder.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a reminder by ID
func (m ReminderModel) GetByID(reminderID int64) (*Reminder, error) {
	query := `
		SELECT id, workspace_id, user_id, target_id, target_type, reminder_type, title, description, due_date, schedule_type, is_active, created_at, updated_at
		FROM reminders
		WHERE id = $1
	`

	reminder := &Reminder{}
	err := m.DB.QueryRow(query, reminderID).Scan(
		&reminder.ID, &reminder.WorkspaceID, &reminder.UserID, &reminder.TargetID, &reminder.TargetType,
		&reminder.ReminderType, &reminder.Title, &reminder.Description, &reminder.DueDate, &reminder.ScheduleType,
		&reminder.IsActive, &reminder.CreatedAt, &reminder.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return reminder, nil
}

// GetByWorkspace retrieves all reminders for a workspace
func (m ReminderModel) GetByWorkspace(workspaceID int64) ([]*Reminder, error) {
	query := `
		SELECT id, workspace_id, user_id, target_id, target_type, reminder_type, title, description, due_date, schedule_type, is_active, created_at, updated_at
		FROM reminders
		WHERE workspace_id = $1 AND is_active = true
		ORDER BY due_date ASC
	`

	rows, err := m.DB.Query(query, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders := []*Reminder{}
	for rows.Next() {
		reminder := &Reminder{}
		err := rows.Scan(
			&reminder.ID, &reminder.WorkspaceID, &reminder.UserID, &reminder.TargetID, &reminder.TargetType,
			&reminder.ReminderType, &reminder.Title, &reminder.Description, &reminder.DueDate, &reminder.ScheduleType,
			&reminder.IsActive, &reminder.CreatedAt, &reminder.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}

	return reminders, rows.Err()
}

// GetDue retrieves reminders that are due (past due_date and active)
func (m ReminderModel) GetDue() ([]*Reminder, error) {
	query := `
		SELECT id, workspace_id, user_id, target_id, target_type, reminder_type, title, description, due_date, schedule_type, is_active, created_at, updated_at
		FROM reminders
		WHERE is_active = true AND due_date <= NOW()
		ORDER BY due_date ASC
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminders := []*Reminder{}
	for rows.Next() {
		reminder := &Reminder{}
		err := rows.Scan(
			&reminder.ID, &reminder.WorkspaceID, &reminder.UserID, &reminder.TargetID, &reminder.TargetType,
			&reminder.ReminderType, &reminder.Title, &reminder.Description, &reminder.DueDate, &reminder.ScheduleType,
			&reminder.IsActive, &reminder.CreatedAt, &reminder.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		reminders = append(reminders, reminder)
	}

	return reminders, rows.Err()
}

// Update modifies an existing reminder
func (m ReminderModel) Update(reminder *Reminder) error {
	query := `
		UPDATE reminders
		SET title = $1, description = $2, due_date = $3, schedule_type = $4, is_active = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`

	err := m.DB.QueryRow(query, reminder.Title, reminder.Description, reminder.DueDate, reminder.ScheduleType, reminder.IsActive, reminder.ID).
		Scan(&reminder.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrRecordNotFound
		}
		return err
	}
	return nil
}

// Delete removes a reminder by ID
func (m ReminderModel) Delete(reminderID int64) error {
	query := `DELETE FROM reminders WHERE id = $1`

	result, err := m.DB.Exec(query, reminderID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// Notification represents a user notification
type Notification struct {
	ID               int64      `json:"id"`
	UserID           int64      `json:"user_id"`
	ReminderID       *int64     `json:"reminder_id"`
	Title            string     `json:"title"`
	Message          *string    `json:"message"`
	NotificationType string     `json:"notification_type"` // 'reminder', 'task_update', 'mention'
	IsRead           bool       `json:"is_read"`
	Channels         string     `json:"channels"` // JSON: {"in_app": true, "email": true}
	CreatedAt        time.Time  `json:"created_at"`
	ReadAt           *time.Time `json:"read_at"`
}

// NotificationModel handles database operations for notifications
type NotificationModel struct {
	DB *sql.DB
}

// Create inserts a new notification
func (m NotificationModel) Create(notification *Notification) error {
	query := `
		INSERT INTO notifications (user_id, reminder_id, title, message, notification_type, channels)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := m.DB.QueryRow(query, notification.UserID, notification.ReminderID, notification.Title, notification.Message,
		notification.NotificationType, notification.Channels).
		Scan(&notification.ID, &notification.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}

// GetByUser retrieves notifications for a user
func (m NotificationModel) GetByUser(userID int64, limit int, offset int) ([]*Notification, error) {
	query := `
		SELECT id, user_id, reminder_id, title, message, notification_type, is_read, channels, created_at, read_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := m.DB.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notifications := []*Notification{}
	for rows.Next() {
		notification := &Notification{}
		err := rows.Scan(
			&notification.ID, &notification.UserID, &notification.ReminderID, &notification.Title, &notification.Message,
			&notification.NotificationType, &notification.IsRead, &notification.Channels, &notification.CreatedAt, &notification.ReadAt,
		)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, rows.Err()
}

// GetUnreadCount returns the count of unread notifications for a user
func (m NotificationModel) GetUnreadCount(userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`

	var count int
	err := m.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// MarkAsRead marks a notification as read
func (m NotificationModel) MarkAsRead(notificationID int64) error {
	query := `UPDATE notifications SET is_read = true, read_at = NOW() WHERE id = $1`

	result, err := m.DB.Exec(query, notificationID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// MarkAllAsRead marks all unread notifications for a user as read
func (m NotificationModel) MarkAllAsRead(userID int64) error {
	query := `UPDATE notifications SET is_read = true, read_at = NOW() WHERE user_id = $1 AND is_read = false`

	_, err := m.DB.Exec(query, userID)
	return err
}

// Delete removes a notification
func (m NotificationModel) Delete(notificationID int64) error {
	query := `DELETE FROM notifications WHERE id = $1`

	result, err := m.DB.Exec(query, notificationID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
