package api

import "fmt"

type TaskStatus string

const (
	TaskStatusCreated  TaskStatus = "created"
	TaskStatusAborted  TaskStatus = "aborted"
	TaskStatusExecuted TaskStatus = "executed"
)

type Options struct {
	RetainResult bool `json:"retain_result"`
	TimeoutSec   int  `json:"timeout_sec"`
}

// Request header `Endpoint: Tasks.Create`
type CreateTaskRequest struct {
	TaskType string         `json:"task_type"`
	Options  map[string]any `json:"options"`
}

func (r CreateTaskRequest) Validate() error {
	if r.TaskType == "" {
		return fmt.Errorf("тело запроса не содержит поле `task_type`")
	}
	return nil
}

type CreateTaskResponse struct {
	TaskID    int            `json:"task_id"`
	TaskType  string         `json:"task_type"`
	Options   map[string]any `json:"options"`
	CreatedAt string         `json:"created_at"`
}

// Request header `Endpoint: Tasks.TaskDetails`
type GetTaskDetailsRequest struct {
	TaskID int `json:"task_id"`
}

func (r GetTaskDetailsRequest) Validate() error {
	if r.TaskID == 0 {
		return fmt.Errorf("тело запроса не содержит поле  `task_id`")
	}
	return nil
}

type GetTaskDetailsResponse struct {
	TaskID        int            `json:"task_id"`
	TaskType      string         `json:"task_type"`
	Options       map[string]any `json:"options"`
	CreatedAt     string         `json:"created_at"`
	Status        TaskStatus     `json:"status"`
	ExecutedAt    string         `json:"executed_at,omitempty"`
	AbortedAt     string         `json:"aborted_at,omitempty"`
	ExecutionTime string         `json:"execution_time"`
}

// Request header `Endpoint: Tasks.List`
type ListTasksRequest struct{}

type TaskSummary struct {
	TaskID   int        `json:"task_id"`
	TaskType string     `json:"task_type"`
	Status   TaskStatus `json:"status"`
}

type ListTasksResponse struct {
	Tasks []TaskSummary `json:"tasks"`
}

// Request header `Endpoint: Tasks.Cancel`
type CancelTaskRequest struct {
	TaskID uint64 `json:"task_id"`
}

func (r CancelTaskRequest) Validate() error {
	if r.TaskID == 0 {
		return fmt.Errorf("тело запроса не содержит поле `task_id`")
	}
	return nil
}

type CancelTaskResponse struct {
	TaskID        int        `json:"id"`
	Status        TaskStatus `json:"status"`
	CreatedAt     string     `json:"created_at"`
	AbortedAt     string     `json:"aborted_at"`
	ExecutionTime string     `json:"execution_time"`
}

// Request header `Endpoint: Tasks.Delete`
type DeleteTaskRequest struct {
	TaskID uint64 `json:"task_id"`
}

func (r DeleteTaskRequest) Validate() error {
	if r.TaskID == 0 {
		return fmt.Errorf("тело запроса не содержит поле `task_id`")
	}
	return nil
}

type DeleteTaskResponse struct {
	TaskID int `json:"task_id"`
}

// Request header `Endpoint: Tasks.TaskResult`
type GetTaskResultRequest struct {
	TaskID int `json:"task_id"`
}

func (r GetTaskResultRequest) Validate() error {
	if r.TaskID == 0 {
		return fmt.Errorf("тело запроса не содержит поле `task_id`")
	}
	return nil
}

type GetTaskResultResponse struct {
	TaskID int    `json:"task_id"`
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}
