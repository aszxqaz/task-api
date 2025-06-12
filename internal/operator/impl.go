package operator

import (
	"context"
	"task-api/internal/executor"
	"task-api/internal/repository"
	"task-api/pkg/timing"
)

type operator struct {
	repo repository.Repository
	exec executor.Executor
}

func New(r repository.Repository, e executor.Executor) *operator {
	o := &operator{r, e}
	o.consumeResults(context.Background())
	return o
}

var _ Operator = (*operator)(nil)

// Cancel implements Operator.
func (h *operator) Cancel(ctx context.Context, taskID uint64) (*repository.Task, error) {
	err := h.exec.Cancel(ctx, taskID)
	if err != nil {
		if execError, ok := err.(*executor.Error); ok {
			if execError.Code() == executor.ErrCodeBadInput {
				return nil, NewError(ErrCodeBadInput, execError.Error())
			}
		}
		return nil, err
	}
	task, err := h.repo.Update(ctx, taskID, func(t repository.Task) (repository.Task, error) {
		t.FinishedAt = timing.Timestamp()
		t.Aborted = true
		return t, nil
	})
	if err != nil {
		if repoError, ok := err.(*repository.Error); ok {
			if repoError.Code() == repository.ErrCodeNotFound {
				return nil, NewError(ErrCodeNotFound, repoError.Error())
			}
		}
		return nil, err
	}
	return task, nil
}

// Create implements Operator.
func (o *operator) Create(ctx context.Context, t Task) (*repository.Task, error) {
	task, err := o.repo.Create(ctx, repository.Task{
		Type:      t.Type(),
		CreatedAt: timing.Timestamp(),
		Options:   t.Options(),
	})
	if err != nil {
		return nil, err
	}
	err = o.exec.Execute(ctx, task.ID, t)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Delete implements Operator.
func (t *operator) Delete(ctx context.Context, taskID uint64) error {
	_ = t.exec.Cancel(ctx, taskID)
	err := t.repo.Delete(ctx, taskID)
	if err != nil {
		if repoErr, ok := err.(*repository.Error); ok {
			if repoErr.Code() == repository.ErrCodeNotFound {
				return NewError(ErrCodeNotFound, repoErr.Error())
			}
		}
		return err
	}
	return nil
}

func (o *operator) consumeResults(ctx context.Context) {
	go func() {
		for result := range o.exec.Results(ctx) {
			o.repo.Update(ctx, result.TaskID, func(t repository.Task) (repository.Task, error) {
				t.FinishedAt = timing.Timestamp()
				t.Result = result.Data
				return t, nil
			})
		}
	}()
}
