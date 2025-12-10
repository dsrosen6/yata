package models

import (
	"context"
	"time"
)

type TimeData struct {
	CreatedOn time.Time
	UpdatedOn time.Time
}

type TaskList struct {
	ID      int
	Title   string
	Details string
	TimeData
}

type ListTaskListsParams struct {
	IDs      []int
	Titles   []string
	Limit    int
	SortBy   string
	SortDesc bool
}

type TaskListRepo interface {
	ListAll(ctx context.Context, params *ListTaskListsParams) ([]*TaskList, error)
	Get(ctx context.Context, id int) (*TaskList, error)
	Create(ctx context.Context, tl *TaskList) (*TaskList, error)
	Update(ctx context.Context, id int, tl *TaskList) (*TaskList, error)
	Delete(ctx context.Context, id int) error
}

type Task struct {
	ID      int
	ListID  int
	Title   string
	Details string
	TimeData
}

type ListTasksParams struct {
	IDs      []int
	ListIDs  []int
	Titles   []string
	Limit    int
	SortBy   string
	SortDesc bool
}

type TaskRepo interface {
	List(ctx context.Context, params *ListTasksParams) ([]*Task, error)
	Get(ctx context.Context, id int) (*Task, error)
	Create(ctx context.Context, t *Task) (*Task, error)
	Update(ctx context.Context, id int, t *Task) (*Task, error)
	Delete(ctx context.Context, id int) error
}
