package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepositoryConstructor(t *testing.T) {
	repo := New()
	assert.Equal(t, repo.currentTaskID, uint64(1))
	assert.NotNil(t, repo.store)
}

func TestRepositoryCreate(t *testing.T) {
	repo := New()
	task := Task{}
	created, err := repo.Create(context.Background(), task)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, uint64(1))
}

func TestRepositoryFind(t *testing.T) {
	repo := New()
	task := Task{}
	ctx := context.Background()

	created, err := repo.Create(ctx, task)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, uint64(1))

	found, err := repo.Find(ctx, uint64(1))
	assert.NoError(t, err)
	assert.Equal(t, found.ID, uint64(1))

	_, err = repo.Find(ctx, uint64(2))
	assert.Error(t, err)
}

func TestRepositoryList(t *testing.T) {
	repo := New()
	task := Task{}
	ctx := context.Background()

	first, err := repo.Create(ctx, task)
	assert.NoError(t, err)
	assert.Equal(t, first.ID, uint64(1))

	second, err := repo.Create(ctx, task)
	assert.NoError(t, err)
	assert.Equal(t, second.ID, uint64(2))

	tasks, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Equal(t, len(tasks), 2)
}

func TestRepositoryDelete(t *testing.T) {
	repo := New()
	task := Task{}
	ctx := context.Background()

	created, err := repo.Create(ctx, task)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, uint64(1))

	err = repo.Delete(ctx, created.ID)
	assert.NoError(t, err)

	_, err = repo.Find(ctx, created.ID)
	assert.Error(t, err)

	err = repo.Delete(ctx, created.ID)
	assert.NoError(t, err)
}

func TestRepositoryUpdate(t *testing.T) {
	repo := New()
	task := Task{}
	ctx := context.Background()
	created, err := repo.Create(ctx, task)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, uint64(1))

	updated, err := repo.Update(ctx, created.ID, func(t Task) (Task, error) {
		t.Aborted = true
		return t, nil
	})
	assert.NoError(t, err)
	assert.True(t, updated.Aborted)

	updated, err = repo.Find(ctx, updated.ID)
	assert.NoError(t, err)
	assert.True(t, updated.Aborted)
}
