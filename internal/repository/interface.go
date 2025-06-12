package repository

import "context"

// Хранилище задач.
type Repository interface {
	List(ctx context.Context) ([]Task, error)
	Find(ctx context.Context, taskID uint64) (*Task, error)
	Create(ctx context.Context, task Task) (*Task, error)
	Delete(ctx context.Context, taskID uint64) error
	Update(ctx context.Context, taskID uint64, update func(t Task) (Task, error)) (*Task, error)
}
