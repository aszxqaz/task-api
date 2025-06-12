package gateway

import (
	"context"
	"fmt"
	"task-api/api"
	"task-api/internal/factory"
	"task-api/internal/operator"
	"task-api/internal/repository"
	"task-api/pkg/timing"
)

type gateway struct {
	repo     repository.Repository
	operator operator.Operator
	factory  factory.Factory
}

func (g *gateway) CancelTask(ctx context.Context, req *api.CancelTaskRequest, res *api.CancelTaskResponse) error {
	task, err := g.operator.Cancel(ctx, req.TaskID)
	if err != nil {
		if operErr, ok := err.(*operator.Error); ok {
			switch operErr.Code() {
			case operator.ErrCodeBadInput:
				return NewError(ErrCodeBadInput, operErr.Error())
			case operator.ErrCodeNotFound:
				return NewError(ErrCodeNotFound, operErr.Error())
			}
		}
		return err
	}
	res.AbortedAt = timing.Format(task.FinishedAt)
	res.CreatedAt = timing.Format(task.CreatedAt)
	res.ExecutionTime = timing.Elapsed(task.FinishedAt, task.CreatedAt)
	res.Status = taskApiStatus(*task)
	res.TaskID = int(task.ID)
	return nil
}

func (g *gateway) CreateTask(ctx context.Context, req *api.CreateTaskRequest, res *api.CreateTaskResponse) error {
	optask, err := g.factory.Construct(req.TaskType, req.Options)
	if err != nil {
		if factoryErr, ok := err.(*factory.Error); ok {
			switch factoryErr.Code() {
			case factory.ErrCodeBadInput, factory.ErrCodeUnknownTaskType:
				return NewError(ErrCodeBadInput, factoryErr.Error())
			}
		}
		return err
	}
	task, err := g.operator.Create(ctx, optask)
	if err != nil {
		if operErr, ok := err.(*operator.Error); ok {
			switch operErr.Code() {
			case operator.ErrCodeBadInput:
				return NewError(ErrCodeBadInput, operErr.Error())
			case operator.ErrCodeNotFound:
				return NewError(ErrCodeNotFound, operErr.Error())
			}
		}
		return err
	}
	res.TaskID = int(task.ID)
	res.CreatedAt = timing.Format(task.CreatedAt)
	res.Options = optask.Options()
	res.TaskType = optask.Type()
	return nil
}

func (g *gateway) DeleteTask(ctx context.Context, req *api.DeleteTaskRequest, res *api.DeleteTaskResponse) error {
	err := g.operator.Delete(ctx, req.TaskID)
	if err != nil {
		if operErr, ok := err.(*operator.Error); ok {
			switch operErr.Code() {
			case operator.ErrCodeBadInput:
				return NewError(ErrCodeBadInput, operErr.Error())
			case operator.ErrCodeNotFound:
				return NewError(ErrCodeNotFound, operErr.Error())
			}
		}
		return err
	}
	res.TaskID = int(req.TaskID)
	return nil
}

func (g *gateway) GetTaskDetails(ctx context.Context, req *api.GetTaskDetailsRequest, res *api.GetTaskDetailsResponse) error {
	task, err := g.repo.Find(ctx, uint64(req.TaskID))
	if err != nil {
		if repoErr, ok := err.(*repository.Error); ok {
			if repoErr.Code() == repository.ErrCodeNotFound {
				return NewError(ErrCodeNotFound, repoErr.Error())
			}
		}
		return err
	}
	res.TaskID = int(task.ID)
	res.Options = task.Options
	res.TaskType = task.Type
	res.CreatedAt = timing.Format(task.CreatedAt)
	res.Status = taskApiStatus(*task)
	if task.FinishedAt != 0 {
		res.ExecutionTime = timing.Elapsed(task.FinishedAt, task.CreatedAt)
		if task.Aborted {
			res.AbortedAt = timing.Format(task.FinishedAt)
		} else {
			res.ExecutedAt = timing.Format(task.FinishedAt)
		}
	} else {
		res.ExecutionTime = timing.Elapsed(timing.Timestamp(), task.CreatedAt)
	}
	return nil
}

func (g *gateway) GetTaskResult(ctx context.Context, req *api.GetTaskResultRequest, res *api.GetTaskResultResponse) error {
	task, err := g.repo.Find(ctx, uint64(req.TaskID))
	if err != nil {
		if repoErr, ok := err.(*repository.Error); ok {
			if repoErr.Code() == repository.ErrCodeNotFound {
				return NewError(ErrCodeNotFound, repoErr.Error())
			}
		}
		return err
	}
	if task.Result == nil && task.Error == "" {
		msg := fmt.Sprintf("задача с id %d не выполнена", req.TaskID)
		return NewError(ErrCodeNotFound, msg)
	}
	res.TaskID = int(task.ID)
	res.Result = task.Result
	res.Error = task.Error
	return nil
}

func (g *gateway) ListTasks(ctx context.Context, req *api.ListTasksRequest, res *api.ListTasksResponse) error {
	rTasks, err := g.repo.List(ctx)
	if err != nil {
		return err
	}
	tasks := make([]api.TaskSummary, 0, len(rTasks))
	for _, task := range rTasks {
		summary := api.TaskSummary{
			TaskID:   int(task.ID),
			TaskType: task.Type,
			Status:   taskApiStatus(task),
		}
		tasks = append(tasks, summary)
	}
	res.Tasks = tasks
	return nil
}

func New(r repository.Repository, o operator.Operator, f factory.Factory) Gateway {
	return &gateway{r, o, f}
}

func taskApiStatus(task repository.Task) api.TaskStatus {
	status := api.TaskStatusCreated
	if task.FinishedAt != 0 {
		if task.Aborted {
			status = api.TaskStatusAborted
		} else {
			status = api.TaskStatusExecuted
		}
	}
	return status
}
