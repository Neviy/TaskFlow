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

// TaskRepository provides methods for working with tasks.
type TaskRepository struct {
	db *pgxpool.Pool
}

// NewTaskRepository creates a new TaskRepository.
func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

// Create creates a new task in the database.
func (r *TaskRepository) Create(ctx context.Context, task *model.Task) error {
	const query = `
		INSERT INTO tasks (
			project_id,
			title,
			description,
			status,
			priority,
			assignee_id,
			deadline
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(
		ctx,
		query,
		task.ProjectID,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.AssigneeID,
		task.Deadline,
	).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create task: %w", err)
	}
	return nil
}

// GetByID returns a task by id.
func (r *TaskRepository) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	task := &model.Task{}
	const query = `
		SELECT
			id,
			project_id,
			title,
			description,
			status,
			priority,
			assignee_id,
			deadline,
			created_at,
			updated_at
		FROM tasks
		WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&task.ID,
		&task.ProjectID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.AssigneeID,
		&task.Deadline,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	return task, nil
}

// Update updates an existing task in the database.
func (r *TaskRepository) Update(ctx context.Context, task *model.Task) error {
	const query = `
		UPDATE tasks
		SET
			title = $1,
			description = $2,
			status = $3,
			priority = $4,
			assignee_id = $5,
			deadline = $6,
			updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at
	`
	err := r.db.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.AssigneeID,
		task.Deadline,
		task.ID,
	).Scan(&task.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("task not found")
	}
	if err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	return nil
}

// Delete deletes a task from the database.
func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM tasks
		WHERE id = $1
	`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("task not found")
	}
	return nil
}

// ListByProjectID returns all tasks belonging to a project.
func (r *TaskRepository) ListByProjectID(ctx context.Context, projectID int64) ([]*model.Task, error) {
	const query = `
		SELECT
			id,
			project_id,
			title,
			description,
			status,
			priority,
			assignee_id,
			deadline,
			created_at,
			updated_at
		FROM tasks
		WHERE project_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("list tasks by project: %w", err)
	}
	defer rows.Close()
	var tasks []*model.Task
	for rows.Next() {
		task := &model.Task{}
		err = rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.AssigneeID,
			&task.Deadline,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}
	return tasks, nil
}

// ListByAssigneeID returns all tasks assigned to a user.
func (r *TaskRepository) ListByAssigneeID(ctx context.Context, assigneeID int64) ([]*model.Task, error) {
	const query = `
		SELECT
			id,
			project_id,
			title,
			description,
			status,
			priority,
			assignee_id,
			deadline,
			created_at,
			updated_at
		FROM tasks
		WHERE assignee_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, assigneeID)
	if err != nil {
		return nil, fmt.Errorf("list tasks by assignee: %w", err)
	}
	defer rows.Close()
	var tasks []*model.Task
	for rows.Next() {
		task := &model.Task{}
		err = rows.Scan(
			&task.ID,
			&task.ProjectID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.AssigneeID,
			&task.Deadline,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}
	return tasks, nil
}
