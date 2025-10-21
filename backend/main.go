package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"todo-app/backend/handlers"
	"todo-app/backend/middleware"
	"todo-app/backend/utils"
)

func main() {
	_ = godotenv.Load()

	mongoURI := getenv("MONGODB_URI", "mongodb://localhost:27017")
	dbName := getenv("DB_NAME", "todoapp")
	frontendDir := filepath.Join("..", "frontend")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := utils.Connect(ctx, mongoURI, dbName)
	if err != nil {
		log.Fatal(err)
	}
	if err := utils.EnsureIndexes(ctx, db); err != nil {
		log.Fatal(err)
	}

	handlers.Init(db, frontendDir)

	r := chi.NewRouter()

	// serve static (CSS, images)
	fs := http.StripPrefix("/static/", http.FileServer(http.Dir(frontendDir)))
	r.Handle("/static/*", fs)

	// auth-less pages
	r.Get("/", handlers.Home)
	r.Get("/login", handlers.ServeLoginPage)
	r.Post("/login", handlers.LoginPost)
	r.Get("/register", handlers.ServeRegisterPage)
	r.Post("/register", handlers.RegisterPost)

	// authed routes
	r.With(middleware.AuthRequired).Get("/app", handlers.AppPage)
	r.With(middleware.AuthRequired).Post("/logout", handlers.LogoutPost)
	r.With(middleware.AuthRequired).Post("/todo", handlers.AddTodo)
	r.With(middleware.AuthRequired).Post("/todo/{id}/toggle", handlers.ToggleTodo)
	r.With(middleware.AuthRequired).Post("/todo/{id}/delete", handlers.DeleteTodo)

	addr := getenv("ADDR", ":8080")
	log.Printf("Server listening on %s", addr)
	log.Printf("Open: http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
