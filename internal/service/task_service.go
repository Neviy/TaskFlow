package service

import (
	"context"
	"errors"
	"strings"
	"taskflow/internal/model"

	"github.com/jackc/pgx/v5"
)

// TaskService handles business logic for tasks.
type TaskService struct {
	taskRepo    TaskRepository
	projectRepo ProjectRepository
	userRepo    UserRepository
}

func NewTaskService(taskRepo TaskRepository, projectRepo ProjectRepository, userRepo UserRepository) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

// Create creates a new task in a project.
func (ts *TaskService) Create(ctx context.Context, title, description string, projectID int64, assigneeID *int64) (*model.Task, error) {
	if projectID <= 0 {
		return nil, ErrInvalidProjectID
	}
	_, err := ts.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, ErrInvalidTaskTitle
	}
	if assigneeID != nil {
		if *assigneeID <= 0 {
			return nil, ErrInvalidUserID
		}
		_, err := ts.userRepo.GetByID(ctx, *assigneeID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrUserNotFound
			}
			return nil, err
		}
	}
	task := &model.Task{
		Title:       title,
		Description: description,
		ProjectID:   projectID,
		AssigneeID:  assigneeID,
	}
	if err := ts.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

// GetByID returns a task by its ID.
func (ts *TaskService) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidTaskID
	}
	task, err := ts.taskRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}

// GetByProjectID returns all tasks in a project.
func (ts *TaskService) GetByProjectID(ctx context.Context, projectID int64) ([]*model.Task, error) {
	if projectID <= 0 {
		return nil, ErrInvalidProjectID
	}
	tasks, err := ts.taskRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// Update updates an existing task.
func (ts *TaskService) Update(ctx context.Context, task *model.Task) error {
	if task == nil {
		return ErrInvalidTask
	}
	if task.ID <= 0 {
		return ErrInvalidTaskID
	}
	task.Title = strings.TrimSpace(task.Title)
	if task.Title == "" {
		return ErrInvalidTaskTitle
	}
	_, err := ts.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	if err := ts.taskRepo.Update(ctx, task); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	return nil
}

// Delete removes a task by its ID.
func (ts *TaskService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidTaskID
	}
	_, err := ts.taskRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	if err := ts.taskRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	return nil
}

// AssignToUser assigns a task to a user.
func (ts *TaskService) AssignToUser(ctx context.Context, taskID int64, userID int64) error {
	if taskID <= 0 {
		return ErrInvalidTaskID
	}
	if userID <= 0 {
		return ErrInvalidUserID
	}
	task, err := ts.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	_, err = ts.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	task.AssigneeID = &userID
	if err := ts.taskRepo.Update(ctx, task); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	return nil
}

// UnassignFromUser removes the assignee from a task.
func (ts *TaskService) UnassignFromUser(ctx context.Context, taskID int64) error {
	if taskID <= 0 {
		return ErrInvalidTaskID
	}
	task, err := ts.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	task.AssigneeID = nil
	if err := ts.taskRepo.Update(ctx, task); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	return nil
}
