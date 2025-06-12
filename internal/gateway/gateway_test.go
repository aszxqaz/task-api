package gateway

import (
	"context"
	"task-api/api"
	"task-api/internal/operator"
	"task-api/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGatewayCreateTask(t *testing.T) {
	repo, oper, fact := setupDeps()
	task := &mockTask{}
	gat := New(repo, oper, fact)
	ctx := context.Background()

	var res api.CreateTaskResponse
	err := gat.CreateTask(ctx, &api.CreateTaskRequest{
		TaskType: task.Type(),
	}, &res)
	assert.Nil(t, err)
	assert.NotNil(t, oper.createdTask)
	assert.NotEmpty(t, res.CreatedAt)
	assert.Equal(t, res.Options, oper.createdTask.Options())
	assert.Equal(t, res.TaskID, 42)
	assert.Equal(t, res.TaskType, oper.createdTask.Type())
}

func TestGatewayListTasks(t *testing.T) {
	repo, oper, fact := setupDeps()
	gat := New(repo, oper, fact)
	ctx := context.Background()

	var res api.ListTasksResponse
	err := gat.ListTasks(ctx, &api.ListTasksRequest{}, &res)
	assert.Nil(t, err)
	assert.Len(t, res.Tasks, 3)
	assert.Equal(t, res.Tasks[0].Status, api.TaskStatusAborted)
	assert.Equal(t, res.Tasks[0].TaskID, 42)
	assert.Equal(t, res.Tasks[0].TaskType, "test42")
	assert.Equal(t, res.Tasks[1].Status, api.TaskStatusExecuted)
	assert.Equal(t, res.Tasks[1].TaskID, 43)
	assert.Equal(t, res.Tasks[1].TaskType, "test43")
	assert.Equal(t, res.Tasks[2].Status, api.TaskStatusCreated)
	assert.Equal(t, res.Tasks[2].TaskID, 44)
	assert.Equal(t, res.Tasks[2].TaskType, "test44")
}

func TestGatewayGetTaskResult(t *testing.T) {
	repo, oper, fact := setupDeps()
	gat := New(repo, oper, fact)
	ctx := context.Background()

	var res api.GetTaskResultResponse
	err := gat.GetTaskResult(ctx, &api.GetTaskResultRequest{
		TaskID: 1,
	}, &res)
	assert.Nil(t, err)
	assert.Empty(t, res.Error)
	assert.Equal(t, res.Result, 42)

	res = api.GetTaskResultResponse{}
	err = gat.GetTaskResult(ctx, &api.GetTaskResultRequest{
		TaskID: 2,
	}, &res)
	assert.Nil(t, err)
	assert.NotEmpty(t, res.Error)
	assert.Equal(t, res.Result, nil)

	res = api.GetTaskResultResponse{}
	err = gat.GetTaskResult(ctx, &api.GetTaskResultRequest{
		TaskID: 13,
	}, &res)
	assert.IsType(t, err, &Error{})
	gatErr := err.(*Error)
	assert.Equal(t, gatErr.code, ErrCodeNotFound)
}

func TestGatewayGetTaskDetails(t *testing.T) {
	repo, oper, fact := setupDeps()
	gat := New(repo, oper, fact)
	ctx := context.Background()
	var res api.GetTaskDetailsResponse

	err := gat.GetTaskDetails(ctx, &api.GetTaskDetailsRequest{
		TaskID: 1,
	}, &res)
	assert.Nil(t, err)
	assert.NotEmpty(t, res.AbortedAt)
	assert.NotEmpty(t, res.CreatedAt)
	assert.Empty(t, res.ExecutedAt)
	assert.Equal(t, res.ExecutionTime, "00:01:00")
	assert.Equal(t, res.Options, map[string]any{
		"test": 42,
	})
	assert.Equal(t, res.Status, api.TaskStatusAborted)
	assert.Equal(t, res.TaskType, "test42")
	assert.Equal(t, res.TaskID, 1)

	res = api.GetTaskDetailsResponse{}
	err = gat.GetTaskDetails(ctx, &api.GetTaskDetailsRequest{
		TaskID: 2,
	}, &res)
	assert.Nil(t, err)
	assert.Equal(t, res.ExecutionTime, "01:00:00")
	assert.Empty(t, res.AbortedAt)
	assert.NotEmpty(t, res.CreatedAt)
	assert.NotEmpty(t, res.ExecutedAt)
	assert.Equal(t, res.Status, api.TaskStatusExecuted)

	res = api.GetTaskDetailsResponse{}
	err = gat.GetTaskDetails(ctx, &api.GetTaskDetailsRequest{
		TaskID: 13,
	}, &res)
	assert.IsType(t, err, &Error{})
	gatErr := err.(*Error)
	assert.Equal(t, gatErr.code, ErrCodeNotFound)
}

func TestGatewayDeleteTask(t *testing.T) {
	repo, oper, fact := setupDeps()
	gat := New(repo, oper, fact)
	ctx := context.Background()

	var res api.DeleteTaskResponse
	err := gat.DeleteTask(ctx, &api.DeleteTaskRequest{TaskID: 42}, &res)
	assert.Nil(t, err)
	assert.Equal(t, oper.deletedTaskID, uint64(42))
	assert.Equal(t, res.TaskID, 42)

	res = api.DeleteTaskResponse{}
	err = gat.DeleteTask(ctx, &api.DeleteTaskRequest{TaskID: 13}, &res)
	assert.Error(t, err)
	assert.IsType(t, err, &Error{})
	gatErr := err.(*Error)
	assert.Equal(t, gatErr.code, ErrCodeNotFound)
}

func TestGatewayCancelTask(t *testing.T) {
	repo, oper, fact := setupDeps()
	gat := New(repo, oper, fact)
	ctx := context.Background()

	var res api.CancelTaskResponse
	err := gat.CancelTask(ctx, &api.CancelTaskRequest{TaskID: 42}, &res)
	assert.Nil(t, err)
	assert.Equal(t, oper.canceledTaskID, uint64(42))
	assert.NotEmpty(t, res.AbortedAt)
	assert.NotEmpty(t, res.CreatedAt)
	assert.Equal(t, res.ExecutionTime, "00:01:30")
	assert.Equal(t, res.Status, api.TaskStatusAborted)
	assert.Equal(t, res.TaskID, 42)

	res = api.CancelTaskResponse{}
	err = gat.CancelTask(ctx, &api.CancelTaskRequest{TaskID: 13}, &res)
	assert.Error(t, err)
	assert.IsType(t, err, &Error{})
	gatErr := err.(*Error)
	assert.Equal(t, gatErr.code, ErrCodeNotFound)
}

func setupDeps() (*mockRepo, *mockOper, *mockFact) {
	return &mockRepo{}, &mockOper{}, &mockFact{}
}

type mockFact struct{}

// Construct implements factory.Factory.
func (m *mockFact) Construct(taskType string, opts map[string]any) (operator.Task, error) {
	return &mockTask{}, nil
}

type mockRepo struct {
	task *repository.Task
}

// Create implements repository.Repository.
func (m *mockRepo) Create(ctx context.Context, task repository.Task) (*repository.Task, error) {
	panic("unimplemented")
}

// Delete implements repository.Repository.
func (m *mockRepo) Delete(ctx context.Context, taskID uint64) error {
	panic("unimplemented")
}

// Find implements repository.Repository.
func (m *mockRepo) Find(ctx context.Context, taskID uint64) (*repository.Task, error) {
	if taskID == 1 {
		return &repository.Task{
			ID:         1,
			CreatedAt:  0,
			FinishedAt: 60,
			Aborted:    true,
			Type:       "test42",
			Options: map[string]any{
				"test": 42,
			},
			Result: 42,
			Error:  "",
		}, nil
	}
	if taskID == 2 {
		return &repository.Task{
			ID:         1,
			CreatedAt:  0,
			FinishedAt: 3600,
			Error:      "test",
		}, nil
	}
	if taskID == 13 {
		return nil, repository.NewError(repository.ErrCodeNotFound, "")
	}
	return nil, nil
}

// List implements repository.Repository.
func (m *mockRepo) List(ctx context.Context) ([]repository.Task, error) {
	return []repository.Task{
		{
			ID:         42,
			Type:       "test42",
			Aborted:    true,
			FinishedAt: 1,
		},
		{
			ID:         43,
			Type:       "test43",
			FinishedAt: 1,
		},
		{
			ID:   44,
			Type: "test44",
		},
	}, nil
}

// Update implements repository.Repository.
func (m *mockRepo) Update(ctx context.Context, taskID uint64, update func(t repository.Task) (repository.Task, error)) (*repository.Task, error) {
	if m.task == nil || m.task.ID != taskID {
		return nil, repository.NewError(repository.ErrCodeNotFound, "")
	}
	updated, _ := update(*m.task)
	m.task = &updated
	return m.task, nil
}

type mockOper struct {
	createdTask    operator.Task
	deletedTaskID  uint64
	canceledTaskID uint64
}

// Cancel implements operator.Operator.
func (m *mockOper) Cancel(ctx context.Context, taskID uint64) (*repository.Task, error) {
	if taskID == 13 {
		return nil, operator.NewError(operator.ErrCodeNotFound, "")
	}
	m.canceledTaskID = taskID
	return &repository.Task{
		ID:         taskID,
		CreatedAt:  0,
		FinishedAt: 90,
		Aborted:    true,
	}, nil
}

// Create implements operator.Operator.
func (m *mockOper) Create(ctx context.Context, task operator.Task) (*repository.Task, error) {
	m.createdTask = task
	return &repository.Task{
		ID: 42,
	}, nil
}

// Delete implements operator.Operator.
func (m *mockOper) Delete(ctx context.Context, taskID uint64) error {
	if taskID == 13 {
		return operator.NewError(operator.ErrCodeNotFound, "")
	}
	m.deletedTaskID = taskID
	return nil
}

type mockTask struct{}

// Execute implements operator.Task.
func (m *mockTask) Execute(context.Context) (any, error) {
	return 42, nil
}

// Options implements operator.Task.
func (m *mockTask) Options() map[string]any {
	return map[string]any{
		"test": 42,
	}
}

// Type implements operator.Task.
func (m *mockTask) Type() string {
	return "test"
}
