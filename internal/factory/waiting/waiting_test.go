package waiting

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWaitingConstructor(t *testing.T) {
	w, err := New(map[string]any{
		"duration_sec": 20,
		"test":         123,
	})
	assert.NoError(t, err)
	assert.Equal(t, w.Type(), TaskType)
	assert.Equal(t, w.Options(), map[string]any{
		"duration_sec": 20,
	})

	w, err = New(map[string]any{
		"duration_sec": "wrong parameter",
	})
	assert.Error(t, err)
	assert.Nil(t, w)

	w, err = New(map[string]any{})
	assert.NoError(t, err)
	assert.Equal(t, w.Options(), map[string]any{
		"duration_sec": defaultDurationSec,
	})
}

func TestWaitingExecute(t *testing.T) {
	w, _ := New(map[string]any{
		"duration_sec": 1,
	})
	ctx := context.Background()
	res, err := w.Execute(ctx)
	assert.NotNil(t, res)
	assert.Nil(t, err)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Millisecond)
	defer cancel()
	res, err = w.Execute(ctx)
	assert.Nil(t, res)
	assert.Error(t, err)
}
