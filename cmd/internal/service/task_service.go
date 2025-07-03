package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log"
	_mathRand "math/rand"
	"time"
	"workmate/cmd/internal/model"
	"workmate/cmd/internal/repository"
)

type TaskService interface {
	CreateTask(ctx context.Context) (string, error)
	GetTask(ctx context.Context, id string) (*model.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type TaskServiceImpl struct {
	repo repository.TaskStorage
}

func NewTaskService(repo repository.TaskStorage) *TaskServiceImpl {
	return &TaskServiceImpl{repo: repo}
}

func (service *TaskServiceImpl) CreateTask(ctx context.Context) (string, error) {
	id := generateID()
	task := &model.Task{
		ID:        id,
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
		Duration:  _mathRand.Intn(3) + 3, // 3–5 minutes
	}
	if err := service.repo.SaveTask(ctx, task); err != nil {
		log.Printf("Service: failed to save task id=%s: %v", id, err)
		return "", err
	}
	log.Printf("Service: task created id=%s", id)
	go processTask(task)
	return id, nil
}

func (service *TaskServiceImpl) GetTask(ctx context.Context, id string) (*model.Task, error) {
	task, err := service.repo.GetTask(ctx, id)
	if err != nil {
		log.Printf("Service: failed to get task id=%s: %v", id, err)
		return nil, err
	}
	log.Printf("Service: task retrieved id=%s", id)
	return task, nil
}

func (service *TaskServiceImpl) DeleteTask(ctx context.Context, id string) error {
	err := service.repo.DeleteTask(ctx, id)
	if err != nil {
		log.Printf("Service: failed to delete task id=%s: %v", id, err)
		return err
	}
	log.Printf("Service: task deleted id=%s", id)
	return nil
}

func processTask(task *model.Task) {
	task.Status = model.StatusInProgress
	time.Sleep(time.Duration(task.Duration) * time.Minute)
	task.Status = model.StatusDone
}

func generateID() string {
	buf := make([]byte, 6)
	rand.Read(buf)
	return base64.URLEncoding.EncodeToString(buf)
}
