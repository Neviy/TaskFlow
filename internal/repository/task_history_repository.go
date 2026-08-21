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

// TaskHistoryRepository provides methods for working with task history.
type TaskHistoryRepository struct {
	db *pgxpool.Pool
}

// NewTaskHistoryRepository creates a new TaskHistoryRepository.
func NewTaskHistoryRepository(db *pgxpool.Pool) *TaskHistoryRepository {
	return &TaskHistoryRepository{
		db: db,
	}
}

// Create creates a new task history record.
func (r *TaskHistoryRepository) Create(
	ctx context.Context,
	history *model.TaskHistory,
) error {
	const query = `
		INSERT INTO task_history (
			task_id,
			changed_by,
			old_status,
			new_status
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := r.db.QueryRow(
		ctx,
		query,
		history.TaskID,
		history.ChangedBy,
		history.OldStatus,
		history.NewStatus,
	).Scan(
		&history.ID,
		&history.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create task history: %w", err)
	}
	return nil
}

// ListByTaskID returns history records for a task.
func (r *TaskHistoryRepository) ListByTaskID(
	ctx context.Context,
	taskID int64,
) ([]*model.TaskHistory, error) {
	const query = `
		SELECT
			id,
			task_id,
			changed_by,
			old_status,
			new_status,
			created_at
		FROM task_history
		WHERE task_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("list task history: %w", err)
	}
	defer rows.Close()
	var history []*model.TaskHistory
	for rows.Next() {
		item := &model.TaskHistory{}
		err := rows.Scan(
			&item.ID,
			&item.TaskID,
			&item.ChangedBy,
			&item.OldStatus,
			&item.NewStatus,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan task history: %w", err)
		}
		history = append(history, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate task history: %w", err)
	}
	return history, nil
}

// GetByID returns a task history record by ID.
func (r *TaskHistoryRepository) GetByID(
	ctx context.Context,
	id int64,
) (*model.TaskHistory, error) {
	history := &model.TaskHistory{}

	const query = `
		SELECT
			id,
			task_id,
			changed_by,
			old_status,
			new_status,
			created_at
		FROM task_history
		WHERE id = $1
	`

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&history.ID,
		&history.TaskID,
		&history.ChangedBy,
		&history.OldStatus,
		&history.NewStatus,
		&history.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("get task history by id: %w", err)
	}

	return history, nil
}
