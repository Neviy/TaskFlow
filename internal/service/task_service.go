// Package service provides the business logic for user management.
package service

import (
	"context"
	"fmt"
	"taskflow/internal/model"
)

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

func (ts *TaskService) Create(ctx context.Context, title, description string, projectID int64, assigneeID *int64) (*model.Task, error) {
	if projectID <= 0 {
		return nil, ErrInvalidProjectID
	}
	project, err := ts.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}
	if title == "" {
		return nil, ErrInvalidTaskTitle
	}
	if assigneeID != nil {
		if *assigneeID <= 0 {
			return nil, ErrInvalidUserID
		}
		user, err := ts.userRepo.GetByID(ctx, *assigneeID)
		if err != nil {
			return nil, fmt.Errorf("get user by id: %w", err)
		}
		if user == nil {
			return nil, ErrUserNotFound
		}
	}
	task := &model.Task{
		Title:       title,
		Description: description,
		ProjectID:   projectID,
		AssigneeID:  assigneeID,
	}
	if err := ts.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	return task, nil
}

func (ts *TaskService) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	if id <= 0 {
		return nil, ErrInvalidTaskID
	}
	task, err := ts.taskRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	return task, nil
}

func (ts *TaskService) GetByProjectID(ctx context.Context, projectID int64) ([]*model.Task, error) {
	if projectID <= 0 {
		return nil, ErrInvalidProjectID
	}
	tasks, err := ts.taskRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list tasks by project id: %w", err)
	}
	return tasks, nil
}

func (ts *TaskService) Update(ctx context.Context, task *model.Task) error {
	if task.ID <= 0 {
		return ErrInvalidTaskID
	}
	if task.Title == "" {
		return ErrInvalidTaskTitle
	}
	existing, err := ts.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		return fmt.Errorf("get task by id: %w", err)
	}
	if existing == nil {
		return ErrTaskNotFound
	}
	if err := ts.taskRepo.Update(ctx, task); err != nil {
		return fmt.Errorf("update task: %w", err)
	}
	return nil
}

func (ts *TaskService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidTaskID
	}
	existing, err := ts.taskRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get task by id: %w", err)
	}
	if existing == nil {
		return ErrTaskNotFound
	}
	if err := ts.taskRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	return nil
}

func (ts *TaskService) AssignToUser(ctx context.Context, taskID int64, userID int64) error {
	if taskID <= 0 {
		return ErrInvalidTaskID
	}
	if userID <= 0 {
		return ErrInvalidUserID
	}
	task, err := ts.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task by id: %w", err)
	}
	if task == nil {
		return ErrTaskNotFound
	}
	user, err := ts.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}
	task.AssigneeID = &userID
	if err := ts.taskRepo.Update(ctx, task); err != nil {
		return fmt.Errorf("assign task to user: %w", err)
	}
	return nil
}

func (ts *TaskService) UnassignFromUser(ctx context.Context, taskID int64) error {
	if taskID <= 0 {
		return ErrInvalidTaskID
	}
	task, err := ts.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task by id: %w", err)
	}
	if task == nil {
		return ErrTaskNotFound
	}
	task.AssigneeID = nil
	if err := ts.taskRepo.Update(ctx, task); err != nil {
		return fmt.Errorf("unassign task from user: %w", err)
	}
	return nil
}
