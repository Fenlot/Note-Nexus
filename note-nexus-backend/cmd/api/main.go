package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Fenlot/note-nexus-backend/internal/data"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type application struct {
	models data.Models
}

func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	log.Println("Initializing database connection...")
	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := &application{
		models: data.NewModels(db),
	}

	if err := applyMigrations(db); err != nil {
		log.Fatal("failed to apply migrations: ", err)
	}
	log.Println("Database connection established and migrations applied.")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://note-nexusgh.onrender.com",
			"https://notenexusgh.netlify.app",
			"http://localhost:5173",
			"http://localhost:5432",
			"http://localhost:3000",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("NoteNexus Backend is ready with DB!"))
	})

	// Auth routes
	r.Post("/v1/signup", app.signupHandler)
	r.Post("/v1/login", app.loginHandler)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(app.requireAuthenticatedUser)

		r.Get("/v1/workspaces", app.listWorkspacesHandler)

		// Workspace scoped routes
		r.Route("/v1/workspaces/{workspaceId}", func(r chi.Router) {
			r.Use(app.requireWorkspaceAccess)

			// Todo routes
			r.Post("/todos", app.createTaskHandler)
			r.Get("/todos", app.listTasksHandler)
			r.Delete("/todos/{id}", app.deleteTaskHandler)
			r.Patch("/todos/{id}", app.updateTaskHandler)
			r.Put("/todos/{id}", app.updateTaskContentHandler)

			// Note routes
			r.Post("/notes", app.createNoteHandler)
			r.Get("/notes", app.listNotesHandler)
			r.Put("/notes/{id}", app.updateNoteHandler)
			r.Delete("/notes/{id}", app.deleteNoteHandler)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "10000"
	}

	addr := "0.0.0.0:" + port
	log.Printf("Server starting on %s...\n", addr)

	err = http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func applyMigrations(db *sql.DB) error {
	migration, err := os.ReadFile("./migrations/001_create_tables.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(migration))
	return err
}
