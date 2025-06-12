package operator

import (
	"context"
	"task-api/internal/executor"
	"task-api/internal/repository"
)

type Task interface {
	executor.Task
	Options() map[string]any
	Type() string
}

// Оператор отдает задачи на исполнение и взаимодействует с хранилищем.
type Operator interface {
	Create(ctx context.Context, task Task) (*repository.Task, error)
	Cancel(ctx context.Context, taskID uint64) (*repository.Task, error)
	Delete(ctx context.Context, taskID uint64) error
}
