package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Fenlot/note-nexus-backend/internal/data"
	"github.com/go-chi/chi/v5"
)

// listNotificationsHandler lists notifications for the authenticated user
func (app *application) listNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userId").(int64)

	// Get limit and offset from query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	notifications, err := app.models.Notifications.GetByUser(userID, limit, offset)
	if err != nil {
		http.Error(w, "Failed to retrieve notifications", http.StatusInternalServerError)
		return
	}

	if notifications == nil {
		notifications = []*data.Notification{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

// getUnreadCountHandler returns the count of unread notifications for the user
func (app *application) getUnreadCountHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userId").(int64)

	count, err := app.models.Notifications.GetUnreadCount(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve unread count", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"unread_count": count})
}

// MarkNotificationAsReadRequest represents the payload for marking a notification as read
type MarkNotificationAsReadRequest struct {
	IsRead bool `json:"is_read"`
}

// markNotificationAsReadHandler marks a notification as read
func (app *application) markNotificationAsReadHandler(w http.ResponseWriter, r *http.Request) {
	notificationID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	var req MarkNotificationAsReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.IsRead {
		if err := app.models.Notifications.MarkAsRead(notificationID); err != nil {
			if err == data.ErrRecordNotFound {
				http.Error(w, "Notification not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Failed to update notification", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

// markAllNotificationsAsReadHandler marks all unread notifications as read
func (app *application) markAllNotificationsAsReadHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userId").(int64)

	if err := app.models.Notifications.MarkAllAsRead(userID); err != nil {
		http.Error(w, "Failed to update notifications", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// deleteNotificationHandler deletes a notification
func (app *application) deleteNotificationHandler(w http.ResponseWriter, r *http.Request) {
	notificationID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid notification ID", http.StatusBadRequest)
		return
	}

	if err := app.models.Notifications.Delete(notificationID); err != nil {
		if err == data.ErrRecordNotFound {
			http.Error(w, "Notification not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete notification", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
