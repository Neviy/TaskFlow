package service

import (
	"context"
	"fmt"
	"taskflow/internal/model"
)

// CommentService contains comment business logic.
type CommentService struct {
	commentRepo CommentRepository
	taskRepo    TaskRepository
	userRepo    UserRepository
}

// NewCommentService creates a new CommentService.
func NewCommentService(
	commentRepo CommentRepository,
	taskRepo TaskRepository,
	userRepo UserRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		taskRepo:    taskRepo,
		userRepo:    userRepo,
	}
}

// Create creates a new comment for a task.
func (cs *CommentService) Create(ctx context.Context, taskID int64, userID int64, content string,
) (*model.Comment, error) {
	if taskID <= 0 {
		return nil, ErrInvalidTaskID
	}
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	if content == "" {
		return nil, ErrInvalidCommentText
	}
	task, err := cs.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	user, err := cs.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
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

// GetByID returns a comment by ID.
func (cs *CommentService) GetByID(ctx context.Context, id int64) (*model.Comment, error) {
	if id <= 0 {
		return nil, ErrInvalidCommentID
	}
	comment, err := cs.commentRepo.GetByID(ctx, id)
	if err != nil {
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
	task, err := cs.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	comments, err := cs.commentRepo.ListByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("list comments by task id: %w", err)
	}
	return comments, nil
}

// Delete removes a comment.
func (cs *CommentService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidCommentID
	}
	comment, err := cs.commentRepo.GetByID(ctx, id)
	if err != nil {
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
