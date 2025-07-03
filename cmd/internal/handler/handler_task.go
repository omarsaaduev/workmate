package handler

import (
	"errors"
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
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := extractIDFromRequest(r)

	task, err := h.service.GetTask(r.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to get task", http.StatusInternalServerError)
		return
	}
	utils.RespondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := extractIDFromRequest(r)

	err := h.service.DeleteTask(r.Context(), id)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// функция для извлечения ID из URL
func extractIDFromRequest(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Path, "/tasks/")
}
