package operator

import (
	"context"
	"sync"
	"task-api/internal/executor"
	"task-api/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOperatorCreate(t *testing.T) {
	repo := &mockRepo{}
	exec := &mockExec{}
	exectask := &mockTask{}
	oper := New(repo, exec)
	ctx := context.Background()

	task, err := oper.Create(ctx, exectask)
	assert.Nil(t, err)
	assert.Equal(t, task.ID, uint64(1))
	assert.Equal(t, exec.taskID, uint64(1))
	assert.Equal(t, exec.task, exectask)
	assert.Equal(t, task.Options["test"], 42)
}

func TestOperatorCancel(t *testing.T) {
	repo := &mockRepo{}
	exec := &mockExec{}
	exectask := &mockTask{}
	oper := New(repo, exec)
	ctx := context.Background()

	task, _ := oper.Create(ctx, exectask)
	task, err := oper.Cancel(ctx, task.ID)
	assert.Nil(t, err)
	assert.Equal(t, task.Aborted, true)
	assert.Equal(t, exec.canceledTaskID, task.ID)
}

func TestOperatorDelete(t *testing.T) {
	repo := &mockRepo{}
	exec := &mockExec{}
	exectask := &mockTask{}
	oper := New(repo, exec)
	ctx := context.Background()

	task, _ := oper.Create(ctx, exectask)
	err := oper.Delete(ctx, task.ID)
	assert.Nil(t, err)
	assert.Equal(t, exec.canceledTaskID, task.ID)
	assert.Nil(t, repo.task)
}

type mockTask struct{}

// Execute implements Task.
func (t *mockTask) Execute(context.Context) (any, error) {
	return nil, nil
}

// Options implements Task.
func (t *mockTask) Options() map[string]any {
	return map[string]any{
		"test": 42,
	}
}

// Type implements Task.
func (t *mockTask) Type() string {
	return "test-type"
}

type mockRepo struct {
	mu   sync.RWMutex
	task *repository.Task
}

// Create implements repository.Repository.
func (r *mockRepo) Create(_ context.Context, task repository.Task) (*repository.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	task.ID = 1
	r.task = &task
	return r.task, nil
}

// Delete implements repository.Repository.
func (r *mockRepo) Delete(ctx context.Context, taskID uint64) error {
	r.task = nil
	return nil
}

// Find implements repository.Repository.
func (r *mockRepo) Find(ctx context.Context, taskID uint64) (*repository.Task, error) {
	if r.task == nil || r.task.ID != taskID {
		return nil, repository.NewError(repository.ErrCodeNotFound, "")
	}
	return r.task, nil
}

// List implements repository.Repository.
func (r *mockRepo) List(ctx context.Context) ([]repository.Task, error) {
	panic("unimplemented")
}

// Update implements repository.Repository.
func (r *mockRepo) Update(ctx context.Context, taskID uint64, update func(t repository.Task) (repository.Task, error)) (*repository.Task, error) {
	if r.task == nil || r.task.ID != taskID {
		return nil, repository.NewError(repository.ErrCodeNotFound, "")
	}
	task, err := update(*r.task)
	r.task = &task
	if err != nil {
		return nil, err
	}
	return r.task, nil
}

type mockExec struct {
	taskID         uint64
	task           executor.Task
	canceledTaskID uint64
}

// Cancel implements executor.Executor.
func (e *mockExec) Cancel(ctx context.Context, taskID uint64) error {
	e.canceledTaskID = taskID
	return nil
}

// Execute implements executor.Executor.
func (e *mockExec) Execute(ctx context.Context, taskID uint64, task executor.Task) error {
	e.taskID = taskID
	e.task = task
	return nil
}

// Results implements executor.Executor.
func (e *mockExec) Results(ctx context.Context) <-chan executor.TaskResult {
	ch := make(<-chan executor.TaskResult)
	return ch
}
