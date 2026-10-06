package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"taskflow/internal/model"

	"github.com/jackc/pgx/v5"
)

// CommentService handles comment-related business logic.
type CommentService struct {
	commentRepo CommentRepository
	taskRepo    TaskRepository
	userRepo    UserRepository
}

func NewCommentService(commentRepo CommentRepository, taskRepo TaskRepository, userRepo UserRepository) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		taskRepo:    taskRepo,
		userRepo:    userRepo,
	}
}

// Create creates a comment after validating the user, task, and comment text.
func (cs *CommentService) Create(ctx context.Context, taskID int64, userID int64, content string,
) (*model.Comment, error) {
	if taskID <= 0 {
		return nil, ErrInvalidTaskID
	}
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	content = strings.TrimSpace(content)
	if len(content) == 0 {
		return nil, ErrInvalidCommentText
	}
	if _, err := cs.userRepo.GetByID(ctx, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	if _, err := cs.taskRepo.GetByID(ctx, taskID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	comment := &model.Comment{
		TaskID:   taskID,
		AuthorID: userID,
		Text:     content,
	}
	if err := cs.commentRepo.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}
	return comment, nil
}

// GetByID returns a comment by its ID.
func (cs *CommentService) GetByID(ctx context.Context, id int64) (*model.Comment, error) {
	if id <= 0 {
		return nil, ErrInvalidCommentID
	}
	comment, err := cs.commentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCommentNotFound
		}
		return nil, fmt.Errorf("get comment by id: %w", err)
	}
	if comment == nil {
		return nil, ErrCommentNotFound
	}
	return comment, nil
}

// GetByTaskID returns all comments for a task.
func (cs *CommentService) GetByTaskID(ctx context.Context, taskID int64) ([]*model.Comment, error) {
	if taskID <= 0 {
		return nil, ErrInvalidTaskID
	}
	_, err := cs.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	comments, err := cs.commentRepo.ListByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("list comments by task id: %w", err)
	}
	return comments, nil
}

// Delete removes a comment by its ID.
func (cs *CommentService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidCommentID
	}
	comment, err := cs.commentRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCommentNotFound
		}
		return fmt.Errorf("get comment by id: %w", err)
	}
	if comment == nil {
		return ErrCommentNotFound
	}
	if err := cs.commentRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}
	return nil
}
