package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Fenlot/note-nexus-backend/internal/data"
	"github.com/go-chi/chi/v5"
)

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	workspaceIDStr := chi.URLParam(r, "workspaceId")
	workspaceID, _ := strconv.ParseInt(workspaceIDStr, 10, 64)
	userID := r.Context().Value(userContextKey).(int64)

	var input struct {
		Title string `json:"title"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	id, err := app.models.Tasks.Insert(input.Title, workspaceID, userID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"id":           id,
		"title":        input.Title,
		"is_completed": false,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (app *application) listTasksHandler(w http.ResponseWriter, r *http.Request) {
	workspaceIDStr := chi.URLParam(r, "workspaceId")
	workspaceID, _ := strconv.ParseInt(workspaceIDStr, 10, 64)

	tasks, err := app.models.Tasks.GetAll(workspaceID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []data.Task{} 
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (app *application) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	workspaceIDStr := chi.URLParam(r, "workspaceId")
	workspaceID, _ := strconv.ParseInt(workspaceIDStr, 10, 64)

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = app.models.Tasks.Delete(id, workspaceID)
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "deleted successfully"}`))
}

func (app *application) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	workspaceIDStr := chi.URLParam(r, "workspaceId")
	workspaceID, _ := strconv.ParseInt(workspaceIDStr, 10, 64)

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	task, err := app.models.Tasks.Get(id, workspaceID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	task.IsCompleted = !task.IsCompleted

	err = app.models.Tasks.Update(task)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (app *application) updateTaskContentHandler(w http.ResponseWriter, r *http.Request) {
	workspaceIDStr := chi.URLParam(r, "workspaceId")
	workspaceID, _ := strconv.ParseInt(workspaceIDStr, 10, 64)

	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var input struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	task, err := app.models.Tasks.Get(id, workspaceID)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	task.Title = input.Title

	err = app.models.Tasks.Update(task)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "updated"}`))
}
