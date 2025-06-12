package factory

import (
	"context"
	"task-api/internal/factory/waiting"
	"task-api/internal/operator"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFactoryDefaultCtor(t *testing.T) {
	f := New()
	task, err := f.Construct(waiting.TaskType, map[string]any{})
	assert.NotNil(t, task)
	assert.Nil(t, err)

	task, err = f.Construct("unknown type", map[string]any{})
	assert.Nil(t, task)
	assert.IsType(t, &Error{}, err)
	facErr, _ := err.(*Error)
	assert.Equal(t, ErrCodeUnknownTaskType, facErr.Code())

	task, err = f.Construct(waiting.TaskType, map[string]any{
		"duration_sec": "wrong parameter",
	})
	assert.Nil(t, task)
	assert.IsType(t, &Error{}, err)
	facErr, _ = err.(*Error)
	assert.Equal(t, ErrCodeBadInput, facErr.Code())
}

func TestFactoryCustomCtor(t *testing.T) {
	f := New(
		WithCtorMap(CtorMap{
			"test": func(opts map[string]any) (operator.Task, error) {
				return &mockTask{}, nil
			},
		}),
	)
	task, err := f.Construct("test", map[string]any{})
	assert.NotNil(t, task)
	assert.Nil(t, err)

	task, err = f.Construct("unknown type", map[string]any{})
	assert.Nil(t, task)
	assert.NotNil(t, err)
}

type mockTask struct{}

// Execute implements operator.Task.
func (m *mockTask) Execute(context.Context) (any, error) {
	panic("unimplemented")
}

// Options implements operator.Task.
func (m *mockTask) Options() map[string]any {
	panic("unimplemented")
}

// Type implements operator.Task.
func (m *mockTask) Type() string {
	panic("unimplemented")
}
