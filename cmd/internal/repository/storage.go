package repository

import (
	"context"
	"errors"
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
	return nil
}

func (ts *TaskStorageMap) GetTask(ctx context.Context, id string) (*model.Task, error) {
	val, ok := ts.tasks.Load(id)
	if !ok {
		return nil, ErrTaskNotFound
	}
	task, ok := val.(*model.Task)
	if !ok {
		return nil, errors.New("invalid task type")
	}
	return task, nil
}

func (ts *TaskStorageMap) DeleteTask(ctx context.Context, id string) error {
	_, ok := ts.tasks.Load(id)
	if !ok {
		return ErrTaskNotFound
	}
	ts.tasks.Delete(id)
	return nil
}
