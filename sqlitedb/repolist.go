package sqlitedb

import (
	"context"

	"github.com/dsrosen6/yata/models"
)

type ListRepo struct {
	q *Queries
}

func NewListRepo(q *Queries) *ListRepo {
	return &ListRepo{
		q: q,
	}
}

func (lr *ListRepo) ListAll(ctx context.Context) ([]*models.List, error) {
	dl, err := lr.q.ListAllLists(ctx)
	if err != nil {
		return nil, err
	}

	return dbListSliceToListSlice(dl), nil
}

func (lr *ListRepo) ListByParentID(ctx context.Context, parentID int64) ([]*models.List, error) {
	dl, err := lr.q.ListListsByParentListID(ctx, &parentID)
	if err != nil {
		return nil, err
	}

	return dbListSliceToListSlice(dl), nil
}

func (lr *ListRepo) Get(ctx context.Context, id int64) (*models.List, error) {
	d, err := lr.q.GetList(ctx, id)
	if err != nil {
		return nil, err
	}

	return dbListToList(d), nil
}

func (lr *ListRepo) Create(ctx context.Context, l *models.List) (*models.List, error) {
	d, err := lr.q.CreateList(ctx, listToCreateParams(l))
	if err != nil {
		return nil, err
	}

	return dbListToList(d), nil
}

func (lr *ListRepo) Update(ctx context.Context, l *models.List) (*models.List, error) {
	d, err := lr.q.UpdateList(ctx, listToUpdateParams(l))
	if err != nil {
		return nil, err
	}

	return dbListToList(d), nil
}

func (lr *ListRepo) Delete(ctx context.Context, id int64) error {
	return lr.q.DeleteList(ctx, id)
}

func listToCreateParams(l *models.List) *CreateListParams {
	return &CreateListParams{
		Title:        l.Title,
		ParentListID: l.ParentID,
	}
}

func listToUpdateParams(l *models.List) *UpdateListParams {
	return &UpdateListParams{
		ID:           l.ID,
		Title:        l.Title,
		ParentListID: l.ParentID,
	}
}

func dbListSliceToListSlice(ds []*List) []*models.List {
	var l []*models.List
	for _, d := range ds {
		l = append(l, dbListToList(d))
	}

	return l
}

func dbListToList(d *List) *models.List {
	return &models.List{
		ID:        d.ID,
		Title:     d.Title,
		ParentID:  d.ParentListID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
