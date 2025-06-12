package repository

import (
	"context"
	"fmt"
	"slices"
	"sync"
)

type Task struct {
	ID         uint64
	CreatedAt  int64
	FinishedAt int64
	Type       string
	Options    map[string]any
	Aborted    bool
	Error      string
	Result     any
}

var _ Repository = (*repository)(nil)

type repository struct {
	mu            sync.RWMutex
	currentTaskID uint64
	store         map[uint64]Task
}

func New() *repository {
	return &repository{
		currentTaskID: 1,
		store:         make(map[uint64]Task),
	}
}

// Create implements Repository.
func (r *repository) Create(ctx context.Context, task Task) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	task.ID = r.currentTaskID
	r.store[task.ID] = task
	r.currentTaskID++
	return &task, nil
}

// Find implements Repository.
func (r *repository) Find(ctx context.Context, taskID uint64) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.store[taskID]
	if !ok {
		msg := fmt.Sprintf("задача с id %d не найдена", taskID)
		return nil, NewError(ErrCodeNotFound, msg)
	}
	return &task, nil
}

// List implements Repository.
func (r *repository) List(context.Context) ([]Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	tasks := make([]Task, 0, len(r.store))
	for _, task := range r.store {
		tasks = append(tasks, task)
	}
	slices.SortFunc(tasks, func(a, b Task) int { return int(a.ID - b.ID) })
	return tasks, nil
}

// Update implements Repository.
func (r *repository) Update(ctx context.Context, taskID uint64, update func(t Task) (Task, error)) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, ok := r.store[taskID]
	if !ok {
		msg := fmt.Sprintf("задача с id %d не найдена", taskID)
		return nil, NewError(ErrCodeNotFound, msg)
	}
	updated, err := update(task)
	if err != nil {
		return nil, err
	}
	r.store[taskID] = updated
	return &updated, nil
}

// Delete implements Repository.
func (r *repository) Delete(ctx context.Context, taskID uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.store, taskID)
	return nil
}
