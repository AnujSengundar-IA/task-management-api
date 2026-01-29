package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"task-management-api/internal/models"
	"task-management-api/internal/store"
	"time"

	"github.com/google/uuid"
)

type TaskHandler struct {
	repo store.TaskRepository
}

func NewTaskHandler(repo store.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{
		"error": msg,
	})
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if input.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is Required")
		return
	}

	task := models.Task{
		ID:        uuid.NewString(),
		Title:     input.Title,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if err := h.repo.Create(r.Context(), task); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create task")
		log.Fatal(err)
		return
	}
	writeJSON(w, http.StatusCreated, task)

}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.repo.GetAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get tasks")
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request, id string) {
	task, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Task Not Found")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request, id string) {
	var input struct {
		Title  string `json:"title"`
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}
	task, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusNotFound, "task not found")
		return
	}

	if input.Title != "" {
		task.Title = input.Title
	}
	if input.Status != "" {
		task.Status = input.Status
	}
	if err := h.repo.Update(r.Context(), task); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update task")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request, id string) {
	if err := h.repo.Delete(r.Context(), id); err != nil {
		if err == store.ErrNotFound {
			writeError(w, http.StatusNotFound, "task not found")
		}
		writeError(w, http.StatusInternalServerError, "failed to delete task")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
