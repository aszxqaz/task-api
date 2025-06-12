package waiting

import (
	"context"
	"fmt"
	"reflect"
	"task-api/internal/operator"
	"task-api/pkg/fromjson"
	"time"
)

const TaskType string = "waiting"

const (
	defaultDurationSec = 10
)

type waitingTask struct {
	durationSec int
}

// "Dummy"-задача, которая спит некоторое время durationSec,
// прежде чем вернуть сообщение.
func New(opts map[string]any) (operator.Task, error) {
	durationSec := defaultDurationSec
	durationSecRaw, ok := opts["duration_sec"]
	if ok {
		if i, ok := fromjson.ParseInt(durationSecRaw); ok {
			durationSec = i
		} else {
			err := fmt.Errorf("параметр `duration_sec` должен быть целочисленным, получено: %v", reflect.TypeOf(durationSecRaw))
			return nil, err
		}
	}

	if durationSec <= 0 {
		err := fmt.Errorf("параметр `duration_sec` для WaitingTask не может быть <= 0. Получено: %d", durationSec)
		return nil, err
	}

	return &waitingTask{durationSec}, nil
}

// Execute implements operator.Task.
func (w *waitingTask) Execute(ctx context.Context) (any, error) {
	dur := time.Duration(w.durationSec) * time.Second
	select {
	case <-time.After(dur):
		msg := fmt.Sprintf("задача говорит \"привет\" спустя %d секунд", w.durationSec)
		return msg, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("задача отменена")
	}
}

// Type implements operator.Task.
func (w *waitingTask) Type() string {
	return string(TaskType)
}

// Options implements operator.Task.
func (t *waitingTask) Options() map[string]any {
	return map[string]any{
		"duration_sec": t.durationSec,
	}
}
