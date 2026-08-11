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

// ProjectMemberRepository provides methods for working with project members.
type ProjectMemberRepository struct {
	db *pgxpool.Pool
}

// NewProjectMemberRepository creates a new ProjectMemberRepository.
func NewProjectMemberRepository(db *pgxpool.Pool) *ProjectMemberRepository {
	return &ProjectMemberRepository{
		db: db,
	}
}

// Create adds a user to a project.
func (r *ProjectMemberRepository) Create(
	ctx context.Context,
	member *model.ProjectMember,
) error {
	const query = `
		INSERT INTO project_members (
			project_id,
			user_id,
			role
		)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		member.ProjectID,
		member.UserID,
		member.Role,
	).Scan(
		&member.ID,
	)

	if err != nil {
		return fmt.Errorf("create project member: %w", err)
	}

	return nil
}

// GetByProjectAndUserID returns a project member by project ID and user ID.
func (r *ProjectMemberRepository) GetByProjectAndUserID(
	ctx context.Context,
	projectID int64,
	userID int64,
) (*model.ProjectMember, error) {
	member := &model.ProjectMember{}

	const query = `
		SELECT
			id,
			project_id,
			user_id,
			role
		FROM project_members
		WHERE project_id = $1
		  AND user_id = $2
	`

	err := r.db.QueryRow(
		ctx,
		query,
		projectID,
		userID,
	).Scan(
		&member.ID,
		&member.ProjectID,
		&member.UserID,
		&member.Role,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("get project member: %w", err)
	}

	return member, nil
}

// ListByProjectID returns all members of a project.
func (r *ProjectMemberRepository) ListByProjectID(
	ctx context.Context,
	projectID int64,
) ([]*model.ProjectMember, error) {
	const query = `
		SELECT
			id,
			project_id,
			user_id,
			role
		FROM project_members
		WHERE project_id = $1
		ORDER BY id ASC
	`

	rows, err := r.db.Query(
		ctx,
		query,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("list project members: %w", err)
	}

	defer rows.Close()

	var members []*model.ProjectMember

	for rows.Next() {
		member := &model.ProjectMember{}

		err := rows.Scan(
			&member.ID,
			&member.ProjectID,
			&member.UserID,
			&member.Role,
		)
		if err != nil {
			return nil, fmt.Errorf("scan project member: %w", err)
		}

		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate project members: %w", err)
	}

	return members, nil
}

// Update changes the role of a project member.
func (r *ProjectMemberRepository) Update(
	ctx context.Context,
	member *model.ProjectMember,
) error {
	const query = `
		UPDATE project_members
		SET role = $1
		WHERE project_id = $2
		  AND user_id = $3
		RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		member.Role,
		member.ProjectID,
		member.UserID,
	).Scan(
		&member.ID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("project member not found")
	}

	if err != nil {
		return fmt.Errorf("update project member: %w", err)
	}

	return nil
}

// Delete removes a user from a project.
func (r *ProjectMemberRepository) Delete(
	ctx context.Context,
	projectID int64,
	userID int64,
) error {
	const query = `
		DELETE FROM project_members
		WHERE project_id = $1
		  AND user_id = $2
	`

	result, err := r.db.Exec(
		ctx,
		query,
		projectID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("delete project member: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("project member not found")
	}

	return nil
}
