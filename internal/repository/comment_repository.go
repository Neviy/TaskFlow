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

// CommentRepository provides methods for working with comments.
type CommentRepository struct {
	db *pgxpool.Pool
}

// NewCommentRepository creates a new CommentRepository.
func NewCommentRepository(db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

// Create creates a new comment in the database.
func (r *CommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	const query = `
		INSERT INTO comments (
			task_id,
			author_id,
			text
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			created_at
	`
	err := r.db.QueryRow(
		ctx,
		query,
		comment.TaskID,
		comment.AuthorID,
		comment.Text,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("create comment: %w", err)
	}
	return nil
}

// GetByID returns a comment by id.
func (r *CommentRepository) GetByID(ctx context.Context, id int64) (*model.Comment, error) {
	comment := &model.Comment{}
	const query = `
		SELECT
			id,
			task_id,
			author_id,
			text,
			created_at
		FROM comments
		WHERE id = $1
	`
	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&comment.ID,
		&comment.TaskID,
		&comment.AuthorID,
		&comment.Text,
		&comment.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get comment by id: %w", err)
	}
	return comment, nil
}

// Delete deletes a comment from the database.
func (r *CommentRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM comments
		WHERE id = $1
	`
	result, err := r.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("comment not found")
	}
	return nil
}

// ListByTaskID returns all comments belonging to a task.
func (r *CommentRepository) ListByTaskID(ctx context.Context, taskID int64) ([]*model.Comment, error) {
	const query = `
		SELECT
			id,
			task_id,
			author_id,
			text,
			created_at
		FROM comments
		WHERE task_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(
		ctx,
		query,
		taskID,
	)
	if err != nil {
		return nil, fmt.Errorf("list comments by task: %w", err)
	}
	defer rows.Close()
	var comments []*model.Comment
	for rows.Next() {
		comment := &model.Comment{}
		err = rows.Scan(
			&comment.ID,
			&comment.TaskID,
			&comment.AuthorID,
			&comment.Text,
			&comment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate comments: %w", err)
	}
	return comments, nil
}
