package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"strings"
	"task-management-api/internal/handlers"
	"task-management-api/internal/store"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("cannot connect to database: ", err)
	}
	defer db.Close()
	repo := store.NewPostgresTaskStore(db)
	taskHandler := handlers.NewTaskHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			taskHandler.CreateTask(w, r)
		case http.MethodGet:
			taskHandler.GetTasks(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/tasks/")

		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			taskHandler.GetTask(w, r, id)
		case http.MethodPut:
			taskHandler.UpdateTask(w, r, id)
		case http.MethodDelete:
			taskHandler.DeleteTask(w, r, id)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
