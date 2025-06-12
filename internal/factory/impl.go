package factory

import (
	"fmt"
	"task-api/internal/factory/waiting"
	"task-api/internal/operator"
)

type factory struct {
	ctorMap CtorMap
}

type CtorMap = map[string]func(opts map[string]any) (operator.Task, error)

var defaultCtorMap CtorMap = CtorMap{
	waiting.TaskType: waiting.New,
}

type Config struct {
	CtorMap CtorMap
}

type Option func(f *factory)

func New(opts ...Option) *factory {
	f := &factory{}
	for _, opt := range opts {
		opt(f)
	}
	if f.ctorMap == nil {
		f.ctorMap = defaultCtorMap
	}
	return f
}

func WithCtorMap(ctorMap CtorMap) Option {
	return func(f *factory) {
		f.ctorMap = ctorMap
	}
}

func (f *factory) Construct(taskType string, opts map[string]any) (operator.Task, error) {
	ctor, ok := f.ctorMap[taskType]
	if !ok {
		msg := fmt.Sprintf("тип задачи неизвестен: %s", taskType)
		return nil, NewError(ErrCodeUnknownTaskType, msg)
	}
	task, err := ctor(opts)
	if err != nil {
		msg := fmt.Sprintf("параметры задачи неверны: %s", err.Error())
		return nil, NewError(ErrCodeBadInput, msg)
	}
	return task, nil
}
