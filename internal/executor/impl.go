package executor

import (
	"context"
	"fmt"
	"task-api/pkg/syncmap"
	"task-api/pkg/timing"
)

type TaskResult struct {
	TaskID    uint64
	Timestamp int64
	Error     error
	Data      any
}

func New() *executor {
	return &executor{results: make(chan TaskResult)}
}

type executor struct {
	results chan TaskResult
	aborts  syncmap.Map[uint64, chan struct{}]
}

// Results implements TaskExecutor.
func (e *executor) Results(context.Context) <-chan TaskResult {
	return e.results
}

var _ Executor = (*executor)(nil)

// Execute implements TaskExecutor.
func (e *executor) Execute(ctx context.Context, taskID uint64, task Task) error {
	abort := make(chan struct{})
	e.aborts.Set(taskID, abort)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		select {
		case ev := <-e.execute(ctx, taskID, task):
			e.aborts.Delete(taskID)
			e.results <- ev
		case <-abort:
			e.aborts.Delete(taskID)
			cancel()
		}
	}()
	return nil
}

// Cancel implements TaskExecutor.
func (e *executor) Cancel(ctx context.Context, taskID uint64) error {
	abort, ok := e.aborts.Get(taskID)
	if !ok {
		msg := fmt.Sprintf("задачи с id %d нет среди выполняемых", taskID)
		return NewError(ErrCodeBadInput, msg)
	}
	close(abort)
	return nil
}

func (e *executor) execute(ctx context.Context, taskID uint64, task Task) chan TaskResult {
	result := make(chan TaskResult)
	go func() {
		data, err := task.Execute(ctx)
		result <- TaskResult{
			TaskID:    taskID,
			Timestamp: timing.Timestamp(),
			Error:     err,
			Data:      data,
		}
		close(result)
	}()
	return result
}
