package main

import (
	"encoding/json"
	"net/http"
)

func (app *application) listWorkspacesHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userContextKey).(int64)

	workspaces, err := app.models.Workspaces.GetUserWorkspaces(userID)
	if err != nil {
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workspaces)
}
