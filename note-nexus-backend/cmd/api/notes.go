package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Fenlot/note-nexus-backend/internal/data"
	"github.com/go-chi/chi/v5"
)

func (app *application) createNoteHandler(w http.ResponseWriter, r *http.Request) {
	workspaceIDStr := chi.URLParam(r, "workspaceId")
	workspaceID, _ := strconv.ParseInt(workspaceIDStr, 10, 64)
	userID := r.Context().Value(userContextKey).(int64)

	var input struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	id, err := app.models.Notes.Insert(input.Content, workspaceID, userID)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"status": "created", "id": id})
}

func (app *application) listNotesHandler(w http.ResponseWriter, r *http.Request) {
	workspaceIDStr := chi.URLParam(r, "workspaceId")
	workspaceID, _ := strconv.ParseInt(workspaceIDStr, 10, 64)

	notes, err := app.models.Notes.GetAll(workspaceID)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}
	
	if notes == nil {
		notes = []data.Note{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (app *application) updateNoteHandler(w http.ResponseWriter, r *http.Request) {
	workspaceIDStr := chi.URLParam(r, "workspaceId")
	workspaceID, _ := strconv.ParseInt(workspaceIDStr, 10, 64)
	idParam := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idParam, 10, 64)

	var input struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	note := &data.Note{
		ID:          id,
		WorkspaceID: workspaceID,
		Content:     input.Content,
	}

	err := app.models.Notes.Update(note)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "updated"}`))
}

func (app *application) deleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	workspaceIDStr := chi.URLParam(r, "workspaceId")
	workspaceID, _ := strconv.ParseInt(workspaceIDStr, 10, 64)
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = app.models.Notes.Delete(id, workspaceID)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "deleted"}`))
}