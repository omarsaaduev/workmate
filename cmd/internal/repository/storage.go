package repository

import (
	"context"
	"errors"
	"log"
	"sync"
	"workmate/cmd/internal/model"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskStorage interface {
	SaveTask(ctx context.Context, task *model.Task) error
	GetTask(ctx context.Context, id string) (*model.Task, error)
	DeleteTask(ctx context.Context, id string) error
}

type TaskStorageMap struct {
	tasks *sync.Map
}

func NewTaskStorageMap(tasks *sync.Map) *TaskStorageMap {
	return &TaskStorageMap{tasks: tasks}
}

func (ts *TaskStorageMap) SaveTask(ctx context.Context, task *model.Task) error {
	ts.tasks.Store(task.ID, task)
	log.Printf("Repo: saved task id=%s", task.ID)
	return nil
}

func (ts *TaskStorageMap) GetTask(ctx context.Context, id string) (*model.Task, error) {
	val, ok := ts.tasks.Load(id)
	if !ok {
		log.Printf("Repo: task not found id=%s", id)
		return nil, ErrTaskNotFound
	}
	task, ok := val.(*model.Task)
	if !ok {
		log.Printf("Repo: task type assertion failed id=%s", id)
		return nil, errors.New("invalid task type")
	}
	log.Printf("Repo: task loaded id=%s", id)
	return task, nil
}

func (ts *TaskStorageMap) DeleteTask(ctx context.Context, id string) error {
	_, ok := ts.tasks.Load(id)
	if !ok {
		log.Printf("Repo: task not found for delete id=%s", id)
		return ErrTaskNotFound
	}
	ts.tasks.Delete(id)
	log.Printf("Repo: task deleted id=%s", id)
	return nil
}
