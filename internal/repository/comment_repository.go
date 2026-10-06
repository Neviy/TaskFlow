package repository

import (
	"context"
	"errors"
	"fmt"

	"taskflow/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CommentRepository works with comments in the database.
type CommentRepository struct {
	db *pgxpool.Pool
}

func NewCommentRepository(db *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{
		db: db,
	}
}

// Create adds a new comment to the database.
func (r *CommentRepository) Create(ctx context.Context, comment *model.Comment) error {
	query := `INSERT INTO comments (task_id,author_id,text)
								VALUES ($1, $2, $3)
								RETURNING id,created_at`
	err := r.db.QueryRow(ctx, query, comment.TaskID, comment.AuthorID, comment.Text).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create comment: %w", err)
	}
	return nil
}

// GetByID finds a comment by its ID.
func (r *CommentRepository) GetByID(ctx context.Context, id int64) (*model.Comment, error) {
	comment := &model.Comment{}
	query := `SELECT id,task_id,author_id,text,created_at
								FROM comments
								WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&comment.ID, &comment.TaskID, &comment.AuthorID, &comment.Text, &comment.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get comment by id: %w", err)
	}
	return comment, nil
}

// Delete removes a comment from the database.
func (r *CommentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM comments
								WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// ListByTaskID returns all comments for a task.
func (r *CommentRepository) ListByTaskID(ctx context.Context, taskID int64) ([]*model.Comment, error) {
	query := `SELECT id,task_id,author_id,text,created_at
								FROM comments
								WHERE task_id = $1
								ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments by task: %w", err)
	}
	defer rows.Close()
	var comments []*model.Comment
	for rows.Next() {
		comment := &model.Comment{}
		err = rows.Scan(&comment.ID, &comment.TaskID, &comment.AuthorID, &comment.Text, &comment.CreatedAt)
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
