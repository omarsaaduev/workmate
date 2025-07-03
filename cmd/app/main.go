package main

import (
	"net/http"
	"sync"
	"workmate/cmd/internal/handler"
	"workmate/cmd/internal/repository"
	"workmate/cmd/internal/service"
)

func main() {
	repo := repository.NewTaskStorageMap(new(sync.Map))
	taskService := service.NewTaskService(repo)
	delivery := handler.NewTaskHandler(taskService)

	router := handler.NewRouter(delivery)
	http.ListenAndServe(":8080", router)

}
