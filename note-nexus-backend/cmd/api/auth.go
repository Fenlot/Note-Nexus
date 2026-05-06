package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"unicode"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (app *application) signupHandler(w http.ResponseWriter, r *http.Request) {
	var input AuthRequest

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if !isValidPassword(input.Password) {
		http.Error(w, "password must be at least 8 characters and include uppercase, number, and special character", http.StatusBadRequest)
		return
	}

	user, err := app.models.Users.Insert(input.Email, input.Password)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	// Create a default workspace for the user
	workspaceName := input.Email + "'s Workspace"
	_, err = app.models.Workspaces.Insert(workspaceName, user.ID)
	if err != nil {
		http.Error(w, "failed to create workspace", http.StatusInternalServerError)
		return
	}

	app.issueTokenResponse(w, user.ID)
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input AuthRequest

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user, err := app.models.Users.Authenticate(input.Email, input.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	app.issueTokenResponse(w, user.ID)
}

func (app *application) issueTokenResponse(w http.ResponseWriter, userID int64) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret_for_development_only"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{Token: tokenString})
}

func isValidPassword(s string) bool {
	if len(s) < 8 {
		return false
	}
	var hasUpper, hasNumber, hasSpecial bool
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	return hasUpper && hasNumber && hasSpecial
}
