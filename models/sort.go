package models

import "sort"

type (
	SortParams struct {
		SortBy
		SortOrder
	}

	SortBy    int
	SortOrder int
)

const (
	SortByTitle SortBy = iota
	SortByComplete
	SortByCreatedAt
	SortByUpdatedAt
)

const (
	SortOrderAsc SortOrder = iota
	SortOrderDesc
)

func SortTasks(tasks []*Task, params SortParams) {
	sort.SliceStable(tasks, func(i, j int) bool {
		var less bool
		switch params.SortBy {
		case SortByTitle:
			less = tasks[i].Title < tasks[j].Title
		case SortByComplete:
			less = !tasks[i].Complete && tasks[j].Complete
		case SortByCreatedAt:
			less = tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
		case SortByUpdatedAt:
			less = tasks[i].UpdatedAt.Before(tasks[j].UpdatedAt)
		}

		if params.SortOrder == SortOrderDesc {
			return !less
		}
		return less
	})
}
