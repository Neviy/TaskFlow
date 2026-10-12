package service

import (
	"context"
	"errors"
	"strings"

	"taskflow/internal/model"

	"github.com/jackc/pgx/v5"
)

// TaskService handles task logic.
type TaskService struct {
	taskRepo    TaskRepository
	projectRepo ProjectRepository
	memberRepo  ProjectMemberRepository
	userRepo    UserRepository
}

func NewTaskService(taskRepo TaskRepository, projectRepo ProjectRepository, memberRepo ProjectMemberRepository, userRepo UserRepository) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		memberRepo:  memberRepo,
		userRepo:    userRepo,
	}
}

// checkProjectAccess checks if a user can access the project.
func (ts *TaskService) checkProjectAccess(ctx context.Context, projectID, userID int64) error {
	if projectID <= 0 || userID <= 0 {
		return ErrForbidden
	}
	project, err := ts.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProjectNotFound
		}
		return err
	}
	if project.OwnerID == userID {
		return nil
	}
	_, err = ts.memberRepo.GetByProjectAndUserID(ctx, projectID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrForbidden
		}
		return err
	}
	return nil
}

// Create creates a task in a project.
func (ts *TaskService) Create(ctx context.Context, title, description string, projectID, userID int64, assigneeID *int64) (*model.Task, error) {
	if projectID <= 0 {
		return nil, ErrInvalidProjectID
	}
	if err := ts.checkProjectAccess(ctx, projectID, userID); err != nil {
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
		if err := ts.checkProjectAccess(ctx, projectID, *assigneeID); err != nil {
			if errors.Is(err, ErrForbidden) {
				return nil, ErrInvalidUserID
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
func (ts *TaskService) GetByID(ctx context.Context, id, userID int64) (*model.Task, error) {
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
	if err := ts.checkProjectAccess(ctx, task.ProjectID, userID); err != nil {
		return nil, err
	}
	return task, nil
}

// GetByProjectID returns tasks from a project.
func (ts *TaskService) GetByProjectID(ctx context.Context, projectID, userID int64) ([]*model.Task, error) {
	if projectID <= 0 {
		return nil, ErrInvalidProjectID
	}
	if err := ts.checkProjectAccess(ctx, projectID, userID); err != nil {
		return nil, err
	}
	return ts.taskRepo.ListByProjectID(ctx, projectID)
}

// Update updates a task.
func (ts *TaskService) Update(ctx context.Context, task *model.Task, userID int64) error {
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
	existing, err := ts.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	if err := ts.checkProjectAccess(ctx, existing.ProjectID, userID); err != nil {
		return err
	}
	task.ProjectID = existing.ProjectID
	if err := ts.taskRepo.Update(ctx, task); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	return nil
}

// Delete removes a task.
func (ts *TaskService) Delete(ctx context.Context, id, userID int64) error {
	if id <= 0 {
		return ErrInvalidTaskID
	}
	task, err := ts.taskRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	if err := ts.checkProjectAccess(ctx, task.ProjectID, userID); err != nil {
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
func (ts *TaskService) AssignToUser(ctx context.Context, taskID, userID, assigneeID int64) error {
	if taskID <= 0 {
		return ErrInvalidTaskID
	}
	if assigneeID <= 0 {
		return ErrInvalidUserID
	}
	task, err := ts.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	if err := ts.checkProjectAccess(ctx, task.ProjectID, userID); err != nil {
		return err
	}
	if err := ts.checkProjectAccess(ctx, task.ProjectID, assigneeID); err != nil {
		if errors.Is(err, ErrForbidden) {
			return ErrInvalidUserID
		}
		return err
	}
	task.AssigneeID = &assigneeID
	if err := ts.taskRepo.Update(ctx, task); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	return nil
}

// UnassignFromUser removes the task assignee.
func (ts *TaskService) UnassignFromUser(ctx context.Context, taskID, userID int64) error {
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
	if err := ts.checkProjectAccess(ctx, task.ProjectID, userID); err != nil {
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
