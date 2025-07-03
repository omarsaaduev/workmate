package main

import (
	"log"
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

	log.Printf("Listening on port 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}

}
