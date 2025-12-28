package sqlitedb

import (
	"context"

	"github.com/dsrosen6/yata/models"
)

type ProjectRepo struct {
	q *Queries
}

func NewProjectRepo(q *Queries) *ProjectRepo {
	return &ProjectRepo{
		q: q,
	}
}

func (pr *ProjectRepo) ListAll(ctx context.Context) ([]*models.Project, error) {
	dp, err := pr.q.ListAllProjects(ctx)
	if err != nil {
		return nil, err
	}

	return dbProjectSliceToProjectSlice(dp), nil
}

func (pr *ProjectRepo) ListByParentID(ctx context.Context, parentID int64) ([]*models.Project, error) {
	dp, err := pr.q.ListProjectsByParentProjectID(ctx, &parentID)
	if err != nil {
		return nil, err
	}

	return dbProjectSliceToProjectSlice(dp), nil
}

func (pr *ProjectRepo) Get(ctx context.Context, id int64) (*models.Project, error) {
	d, err := pr.q.GetProject(ctx, id)
	if err != nil {
		return nil, err
	}

	return dbProjectToProject(d), nil
}

func (pr *ProjectRepo) GetByTitle(ctx context.Context, title string) (*models.Project, error) {
	d, err := pr.q.GetProjectByTitle(ctx, title)
	if err != nil {
		return nil, err
	}

	return dbProjectToProject(d), nil
}

func (pr *ProjectRepo) Create(ctx context.Context, p *models.Project) (*models.Project, error) {
	d, err := pr.q.CreateProject(ctx, projectToCreateParams(p))
	if err != nil {
		return nil, err
	}

	return dbProjectToProject(d), nil
}

func (pr *ProjectRepo) Update(ctx context.Context, p *models.Project) (*models.Project, error) {
	d, err := pr.q.UpdateProject(ctx, projectToUpdateParams(p))
	if err != nil {
		return nil, err
	}

	return dbProjectToProject(d), nil
}

func (pr *ProjectRepo) Delete(ctx context.Context, id int64) error {
	return pr.q.DeleteProject(ctx, id)
}

func projectToCreateParams(p *models.Project) *CreateProjectParams {
	return &CreateProjectParams{
		Title:           p.Title,
		ParentProjectID: p.ParentID,
	}
}

func projectToUpdateParams(p *models.Project) *UpdateProjectParams {
	return &UpdateProjectParams{
		ID:              p.ID,
		Title:           p.Title,
		ParentProjectID: p.ParentID,
	}
}

func dbProjectSliceToProjectSlice(ds []*Project) []*models.Project {
	var p []*models.Project
	for _, d := range ds {
		p = append(p, dbProjectToProject(d))
	}

	return p
}

func dbProjectToProject(d *Project) *models.Project {
	return &models.Project{
		ID:        d.ID,
		Title:     d.Title,
		ParentID:  d.ParentProjectID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
