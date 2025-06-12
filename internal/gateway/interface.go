package gateway

import (
	"context"
	"task-api/api"
)

// Точка входа публичного апи (веб-сервис).
type Gateway interface {
	CreateTask(context.Context, *api.CreateTaskRequest, *api.CreateTaskResponse) error
	ListTasks(context.Context, *api.ListTasksRequest, *api.ListTasksResponse) error
	CancelTask(context.Context, *api.CancelTaskRequest, *api.CancelTaskResponse) error
	DeleteTask(context.Context, *api.DeleteTaskRequest, *api.DeleteTaskResponse) error
	GetTaskDetails(context.Context, *api.GetTaskDetailsRequest, *api.GetTaskDetailsResponse) error
	GetTaskResult(context.Context, *api.GetTaskResultRequest, *api.GetTaskResultResponse) error
}
