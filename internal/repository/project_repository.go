// Package repository provides database access layer.
package repository

import (
	"context"
	"errors"
	"fmt"

	"taskflow/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProjectRepository provides methods for working with projects.
type ProjectRepository struct {
	db *pgxpool.Pool
}

// NewProjectRepository creates a new ProjectRepository.
func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{
		db: db,
	}
}

// Create creates a new project in the database.
func (r *ProjectRepository) Create(ctx context.Context, project *model.Project) error {
	const query = `
		INSERT INTO projects (name, description, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(
		ctx,
		query,
		project.Name,
		project.Description,
		project.OwnerID,
	).Scan(
		&project.ID,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create project: %w", err)
	}
	return nil
}

// GetByID returns a project by id.
func (r *ProjectRepository) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	project := &model.Project{}
	const query = `
		SELECT
			id,
			name,
			description,
			owner_id,
			created_at,
			updated_at
		FROM projects
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.OwnerID,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}
	return project, nil
}

// Update updates an existing project in the database.
func (r *ProjectRepository) Update(ctx context.Context, project *model.Project) error {
	const query = `
		UPDATE projects
		SET
			name = $1,
			description = $2,
			updated_at = NOW()
		WHERE id = $3
		RETURNING updated_at
	`
	err := r.db.QueryRow(
		ctx,
		query,
		project.Name,
		project.Description,
		project.ID,
	).Scan(&project.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("project not found")
	}
	if err != nil {
		return fmt.Errorf("update project: %w", err)
	}
	return nil
}

// Delete deletes a project from the database.
func (r *ProjectRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM projects
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("project not found")
	}
	return nil
}

// ListByOwnerID returns all projects owned by a user.
func (r *ProjectRepository) ListByOwnerID(ctx context.Context, ownerID int64) ([]*model.Project, error) {
	const query = `
		SELECT
			id,
			name,
			description,
			owner_id,
			created_at,
			updated_at
		FROM projects
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list projects by owner: %w", err)
	}
	defer rows.Close()
	var projects []*model.Project
	for rows.Next() {
		project := &model.Project{}
		err = rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.OwnerID,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}
		projects = append(projects, project)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate projects: %w", err)
	}
	return projects, nil
}
