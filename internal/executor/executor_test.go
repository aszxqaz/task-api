package executor

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type successTask struct{}

type failureTask struct{}

func (f successTask) Execute(context.Context) (any, error) {
	time.Sleep(time.Second)
	return 42, nil
}

func (f failureTask) Execute(context.Context) (any, error) {
	time.Sleep(time.Second)
	return nil, fmt.Errorf("error")
}

func TestExecutor(t *testing.T) {
	exec := New()
	ctx := context.Background()
	task1 := successTask{}
	err := exec.Execute(ctx, 1, task1)
	assert.NoError(t, err)
	result := <-exec.Results(ctx)
	assert.Equal(t, result.Data, 42)
	assert.Equal(t, result.Error, nil)
	task2 := failureTask{}
	err = exec.Execute(ctx, 1, task2)
	assert.NoError(t, err)
	result = <-exec.Results(ctx)
	assert.Equal(t, result.Data, nil)
	assert.Equal(t, result.Error.Error(), "error")
}
