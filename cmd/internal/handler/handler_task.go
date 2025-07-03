package handler

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"workmate/cmd/internal/pkg/utils"
	"workmate/cmd/internal/repository"
	"workmate/cmd/internal/service"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	id, err := h.service.CreateTask(r.Context())
	if err != nil {
		log.Printf("Error creating task: %v", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}
	log.Printf("Task created successfully: id=%s", id)
	utils.RespondJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := extractIDFromRequest(r)
	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			log.Printf("Task not found: id=%s", id)
			http.NotFound(w, r)
			return
		}
		log.Printf("Error getting task (id=%s): %v", id, err)
		http.Error(w, "Failed to get task", http.StatusInternalServerError)
		return
	}
	log.Printf("Task retrieved: id=%s", id)
	utils.RespondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := extractIDFromRequest(r)
	err := h.service.DeleteTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			log.Printf("Task not found for delete: id=%s", id)
			http.NotFound(w, r)
			return
		}
		log.Printf("Error deleting task (id=%s): %v", id, err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	log.Printf("Task deleted successfully: id=%s", id)
	w.WriteHeader(http.StatusNoContent)
}

// функция для извлечения ID из URL
func extractIDFromRequest(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Path, "/tasks/")
}
