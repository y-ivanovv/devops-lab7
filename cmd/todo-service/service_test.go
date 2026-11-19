package main

import (
	"context"
	"errors"
	"testing"

	"example.com/devops-lab6/internal/api"
)

func TestTaskLifecycle(t *testing.T) {
	ctx  := context.Background()
	s := service{api.NewMemory[Task]()}
	id, err := s.Create(ctx, Task{Title: "Prepare lab", Priority: 2, Done: true})
	if err != nil || id == "" {
		t.Fatalf("Create: id=%q err=%v", id, err)
	}
	task, err := s.GetByID(ctx, id)
	if err != nil || task.ID != id || task.Title != "Prepare lab" || task.Priority != 2 || task.Done || task.CreatedOn <= 0 {
		t.Fatalf("GetByID after create: task=%+v err=%v", task, err)
	}
	if err := s.MarkDone(ctx, id); err != nil {
		t.Fatal(err)
	}
	task, err = s.GetByID(ctx, id)
	if err != nil || !task.Done {
		t.Fatalf("GetByID after MarkDone: task=%+v err=%v", task, err)
	}
}

func TestInvalidRequests(t *testing.T) {
	ctx  := context.Background()
	s := service{api.NewMemory[Task]()}
	if _, err := s.Create(ctx, Task{Title: "  "}); !errors.Is(err, api.ErrInvalid) {
		t.Fatalf("blank title: got %v", err)
	}
	if _, err := s.GetByID(ctx, "missing"); !errors.Is(err, api.ErrNotFound) {
		t.Fatalf("missing read: got %v", err)
	}
	if err := s.MarkDone(ctx, "missing"); !errors.Is(err, api.ErrNotFound) {
		t.Fatalf("missing update: got %v", err)
	}
}
