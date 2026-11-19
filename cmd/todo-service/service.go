package main

import (
	"context"
	"strings"
	"time"

	"example.com/devops-lab6/internal/api"
	"github.com/google/uuid"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
	CreatedOn   int64  `json:"created_on"`
	Priority    int    `json:"priority"`
}

type Service interface {
	Create(context.Context, Task) (string, error)
	GetByID(context.Context, string) (Task, error)
	MarkDone(context.Context, string) error
}

type Repository interface {
	Put(string, Task)
	Get(string) (Task, error)
	Change(string, func(*Task)) error
}

type service struct{ repo Repository }

func (s service) Create(_ context.Context, task Task) (string, error) {
	if strings.TrimSpace(task.Title) == "" {
		return "", api.ErrInvalid
	}
	task.ID = uuid.NewString()
	task.CreatedOn = time.Now().Unix()
	task.Done = false
	s.repo.Put(task.ID, task)
	return task.ID, nil
}

func (s service) GetByID(_ context.Context, id string) (Task, error) {
	return s.repo.Get(id)
}

func (s service) MarkDone(_ context.Context, id string) error {
	return s.repo.Change(id, func(task *Task) { task.Done = true })
}
