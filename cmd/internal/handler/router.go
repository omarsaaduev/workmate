package handler

import (
	"net/http"
	"strings"
)

// NewRouter — создает http.Handler с маршрутизацией только на /tasks/
func NewRouter(taskHandler *TaskHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем ID из пути
		id := strings.TrimPrefix(r.URL.Path, "/tasks/")
		id = strings.TrimSuffix(id, "/") // на всякий случай убираем финальный слэш

		switch {
		case id == "" && r.Method == http.MethodPost:
			// POST /tasks/ — создать задачу
			taskHandler.CreateTask(w, r)
			return

		case id != "" && r.Method == http.MethodGet:
			// GET /tasks/{id} — получить задачу
			taskHandler.GetTask(w, r)
			return

		case id != "" && r.Method == http.MethodDelete:
			// DELETE /tasks/{id} — удалить задачу
			taskHandler.DeleteTask(w, r)
			return

		default:
			if r.Method != http.MethodPost && r.Method != http.MethodGet && r.Method != http.MethodDelete {
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			} else {
				http.NotFound(w, r)
			}
			return
		}
	})

	return mux
}
