package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Fenlot/note-nexus-backend/internal/data"
	"github.com/go-chi/chi/v5"
)

// CreateReminderRequest represents the request payload for creating a reminder
type CreateReminderRequest struct {
	TargetID     *int64  `json:"target_id"`
	TargetType   string  `json:"target_type"` // 'task', 'note', 'custom'
	ReminderType string  `json:"reminder_type"`
	Title        string  `json:"title"`
	Description  *string `json:"description"`
	DueDate      string  `json:"due_date"` // RFC3339 format
	ScheduleType string  `json:"schedule_type"`
}

// UpdateReminderRequest represents the request payload for updating a reminder
type UpdateReminderRequest struct {
	Title        string  `json:"title"`
	Description  *string `json:"description"`
	DueDate      string  `json:"due_date"`
	ScheduleType string  `json:"schedule_type"`
	IsActive     bool    `json:"is_active"`
}

// createReminderHandler creates a new reminder
func (app *application) createReminderHandler(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.Context().Value("workspaceId").(int64)
	userID := r.Context().Value("userId").(int64)

	var req CreateReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse due date
	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		http.Error(w, "Invalid due_date format. Use RFC3339", http.StatusBadRequest)
		return
	}

	reminder := &data.Reminder{
		WorkspaceID:  workspaceID,
		UserID:       userID,
		TargetID:     req.TargetID,
		TargetType:   req.TargetType,
		ReminderType: req.ReminderType,
		Title:        req.Title,
		Description:  req.Description,
		DueDate:      dueDate,
		ScheduleType: req.ScheduleType,
		IsActive:     true,
	}

	if err := app.models.Reminders.Create(reminder); err != nil {
		http.Error(w, "Failed to create reminder", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(reminder)
}

// listRemindersHandler lists all reminders for a workspace
func (app *application) listRemindersHandler(w http.ResponseWriter, r *http.Request) {
	workspaceID := r.Context().Value("workspaceId").(int64)

	reminders, err := app.models.Reminders.GetByWorkspace(workspaceID)
	if err != nil {
		http.Error(w, "Failed to retrieve reminders", http.StatusInternalServerError)
		return
	}

	if reminders == nil {
		reminders = []*data.Reminder{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminders)
}

// getReminderHandler retrieves a single reminder
func (app *application) getReminderHandler(w http.ResponseWriter, r *http.Request) {
	reminderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid reminder ID", http.StatusBadRequest)
		return
	}

	reminder, err := app.models.Reminders.GetByID(reminderID)
	if err != nil {
		if err == data.ErrRecordNotFound {
			http.Error(w, "Reminder not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve reminder", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminder)
}

// updateReminderHandler updates an existing reminder
func (app *application) updateReminderHandler(w http.ResponseWriter, r *http.Request) {
	reminderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid reminder ID", http.StatusBadRequest)
		return
	}

	reminder, err := app.models.Reminders.GetByID(reminderID)
	if err != nil {
		if err == data.ErrRecordNotFound {
			http.Error(w, "Reminder not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve reminder", http.StatusInternalServerError)
		return
	}

	var req UpdateReminderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse due date
	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		http.Error(w, "Invalid due_date format. Use RFC3339", http.StatusBadRequest)
		return
	}

	reminder.Title = req.Title
	reminder.Description = req.Description
	reminder.DueDate = dueDate
	reminder.ScheduleType = req.ScheduleType
	reminder.IsActive = req.IsActive

	if err := app.models.Reminders.Update(reminder); err != nil {
		http.Error(w, "Failed to update reminder", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reminder)
}

// deleteReminderHandler deletes a reminder
func (app *application) deleteReminderHandler(w http.ResponseWriter, r *http.Request) {
	reminderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid reminder ID", http.StatusBadRequest)
		return
	}

	if err := app.models.Reminders.Delete(reminderID); err != nil {
		if err == data.ErrRecordNotFound {
			http.Error(w, "Reminder not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete reminder", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
