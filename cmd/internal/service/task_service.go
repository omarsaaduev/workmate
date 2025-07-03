package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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
		return "", err
	}
	go processTask(task)
	return id, nil
}

func (service *TaskServiceImpl) GetTask(ctx context.Context, id string) (*model.Task, error) {
	return service.repo.GetTask(ctx, id)
}

func (service *TaskServiceImpl) DeleteTask(ctx context.Context, id string) error {
	return service.repo.DeleteTask(ctx, id)
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
