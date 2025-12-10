package inmem

import (
	"context"
	"sort"
	"sync"

	"github.com/dsrosen6/yata/models"
)

type TaskRepo struct {
	mu     sync.RWMutex
	data   map[int]models.Task
	nextID int
}

func NewTaskRepo() *TaskRepo {
	return &TaskRepo{
		data:   make(map[int]models.Task),
		nextID: 1,
	}
}

func (r *TaskRepo) List(ctx context.Context, params *models.ListTasksParams) ([]*models.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks := make([]models.Task, 0, len(r.data))
	for _, t := range r.data {
		tasks = append(tasks, t)
	}

	tasks = filterTasks(tasks, params)
	return taskSliceToPtrs(tasks), nil
}

func filterTasks(tasks []models.Task, params *models.ListTasksParams) []models.Task {
	if params == nil {
		return tasks
	}

	if len(params.IDs) > 0 {
		tasks = filterTasksByID(tasks, params.IDs)
	}

	if len(params.ListIDs) > 0 {
		tasks = filterTasksByListID(tasks, params.ListIDs)
	}

	if len(params.Titles) > 0 {
		tasks = filterTasksByTitle(tasks, params.Titles)
	}

	return sortTasks(tasks, params.SortBy, params.SortDesc)
}

func filterTasksByID(tasks []models.Task, ids []int) []models.Task {
	idSet := make(map[int]struct{}, len(ids))
	for _, id := range ids {
		idSet[id] = struct{}{}
	}

	filtered := tasks[:0]
	for _, t := range tasks {
		if _, ok := idSet[t.ID]; ok {
			filtered = append(filtered, t)
		}
	}

	return filtered
}

func filterTasksByListID(tasks []models.Task, listIDs []int) []models.Task {
	idSet := make(map[int]struct{}, len(listIDs))
	for _, id := range listIDs {
		idSet[id] = struct{}{}
	}

	filtered := tasks[:0]
	for _, t := range tasks {
		if _, ok := idSet[t.ListID]; ok {
			filtered = append(filtered, t)
		}
	}

	return filtered
}

func filterTasksByTitle(tasks []models.Task, titles []string) []models.Task {
	titleSet := make(map[string]struct{})
	for _, n := range titles {
		titleSet[n] = struct{}{}
	}

	filtered := tasks[:0]
	for _, t := range tasks {
		if _, ok := titleSet[t.Title]; ok {
			filtered = append(filtered, t)
		}
	}

	return filtered
}

func sortTasks(tasks []models.Task, sortBy string, desc bool) []models.Task {
	if sortBy != "" {
		sort.SliceStable(tasks, func(i, j int) bool {
			var less bool
			switch sortBy {
			case "id":
				less = tasks[i].ID < tasks[j].ID
			case "title":
				less = tasks[i].Title < tasks[j].Title
			case "created-on":
				less = tasks[i].CreatedOn.Before(tasks[j].CreatedOn)
			case "updated-on":
				less = tasks[i].UpdatedOn.Before(tasks[j].UpdatedOn)
			default:
				return false
			}

			if desc {
				return !less
			}

			return less
		})
	}

	return tasks
}

func taskSliceToPtrs(tasks []models.Task) []*models.Task {
	out := make([]*models.Task, len(tasks))
	for i := range tasks {
		t := tasks[i]
		out[i] = &t
	}

	return out
}
