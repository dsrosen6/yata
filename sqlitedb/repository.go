package sqlitedb

import (
	"context"

	"github.com/dsrosen6/yata/models"
)

type TaskRepo struct {
	q *Queries
}

func NewRepos(q *Queries) *models.AllRepos {
	return &models.AllRepos{
		Tasks: NewTaskRepo(q),
	}
}

func NewTaskRepo(q *Queries) *TaskRepo {
	return &TaskRepo{
		q: q,
	}
}

func (tr *TaskRepo) ListAll(ctx context.Context) ([]*models.Task, error) {
	dt, err := tr.q.ListAllTasks(ctx)
	if err != nil {
		return nil, err
	}

	return dbTaskSliceToTaskSlice(dt), nil
}

func (tr *TaskRepo) Get(ctx context.Context, id int64) (*models.Task, error) {
	d, err := tr.q.GetTask(ctx, id)
	if err != nil {
		return nil, err
	}

	return dbTaskToTask(d), nil
}

func (tr *TaskRepo) Create(ctx context.Context, t *models.Task) (*models.Task, error) {
	d, err := tr.q.CreateTask(ctx, taskToCreateParams(t))
	if err != nil {
		return nil, err
	}

	return dbTaskToTask(d), nil
}

func (tr *TaskRepo) Update(ctx context.Context, t *models.Task) (*models.Task, error) {
	d, err := tr.q.UpdateTask(ctx, taskToUpdateParams(t))
	if err != nil {
		return nil, err
	}

	return dbTaskToTask(d), nil
}

func (tr *TaskRepo) Delete(ctx context.Context, id int64) error {
	return tr.q.DeleteTask(ctx, id)
}

func taskToCreateParams(t *models.Task) *CreateTaskParams {
	return &CreateTaskParams{
		Title:    t.Title,
		Complete: t.Complete,
	}
}

func taskToUpdateParams(t *models.Task) *UpdateTaskParams {
	return &UpdateTaskParams{
		ID:       t.ID,
		Title:    t.Title,
		Complete: t.Complete,
	}
}

func dbTaskSliceToTaskSlice(ds []*Task) []*models.Task {
	var t []*models.Task
	for _, d := range ds {
		t = append(t, dbTaskToTask(d))
	}

	return t
}

func dbTaskToTask(d *Task) *models.Task {
	return &models.Task{
		ID:        d.ID,
		Title:     d.Title,
		Complete:  d.Complete,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
