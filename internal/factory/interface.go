package factory

import "task-api/internal/operator"

// Фабрика задач.
type Factory interface {
	Construct(taskType string, opts map[string]any) (operator.Task, error)
}
