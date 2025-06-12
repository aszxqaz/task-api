package executor

import "context"

type Task interface {
	Execute(context.Context) (any, error)
}

// Берет задачи на исполнение, исполняет, отменяет и возвращает результаты.
type Executor interface {
	Execute(ctx context.Context, taskID uint64, task Task) error
	Cancel(ctx context.Context, taskID uint64) error
	Results(ctx context.Context) <-chan TaskResult
}
