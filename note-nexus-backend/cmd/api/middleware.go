package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := headerParts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			secret = "default_secret_for_development_only"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}

		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, "invalid sub claim", http.StatusUnauthorized)
			return
		}

		userID := int64(userIDFloat)

		ctx := context.WithValue(r.Context(), userContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) requireWorkspaceAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(userContextKey).(int64)

		workspaceIDStr := chi.URLParam(r, "workspaceId")
		workspaceID, err := strconv.ParseInt(workspaceIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid workspace id", http.StatusBadRequest)
			return
		}

		hasAccess, err := app.models.Workspaces.HasAccess(workspaceID, userID)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		if !hasAccess {
			http.Error(w, "forbidden: you do not have access to this workspace", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
