package repository

import (
	"context"
	"errors"
	"fmt"

	"taskflow/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProjectRepository works with projects in the database.
type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

// Create adds a new project to the database.
func (r *ProjectRepository) Create(ctx context.Context, project *model.Project) error {
	query := `INSERT INTO projects (name,description,owner_id)
						VALUES ($1, $2, $3)
						RETURNING id,created_at,updated_at`
	err := r.db.QueryRow(ctx, query, project.Name, project.Description, project.OwnerID).Scan(&project.ID, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

// GetByID finds a project by its ID.
func (r *ProjectRepository) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	project := &model.Project{}
	query := `SELECT id,name,description,owner_id,created_at,updated_at
						FROM projects
						WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get project by id: %w", err)
	}
	return project, nil
}

// Update changes an existing project in the database.
func (r *ProjectRepository) Update(ctx context.Context, project *model.Project) error {
	query := `UPDATE projects
						SET name = $1,
								description = $2,
								updated_at = NOW()
						WHERE id = $3
						RETURNING updated_at`
	err := r.db.QueryRow(ctx, query, project.Name, project.Description, project.ID).Scan(&project.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}
	return nil
}

// Delete removes a project from the database.
func (r *ProjectRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM projects
						WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// ListByOwnerID returns all projects owned by a user.
func (r *ProjectRepository) ListByOwnerID(ctx context.Context, ownerID int64) ([]*model.Project, error) {
	query := `SELECT id,name,description,owner_id,created_at,updated_at
						FROM projects
						WHERE owner_id = $1
						ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list projects by owner: %w", err)
	}
	defer rows.Close()
	var projects []*model.Project
	for rows.Next() {
		project := &model.Project{}
		err = rows.Scan(&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}
		projects = append(projects, project)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate projects: %w", err)
	}
	return projects, nil
}
