package service

import (
	"context"
	"fmt"
	"taskflow/internal/model"
)

// TaskHistoryService contains task history business logic.
type TaskHistoryService struct {
	historyRepo TaskHistoryRepository
	taskRepo    TaskRepository
	projectRepo ProjectRepository
}

// NewTaskHistoryService creates a new TaskHistoryService.
func NewTaskHistoryService(
	historyRepo TaskHistoryRepository,
	taskRepo TaskRepository,
	projectRepo ProjectRepository,
) *TaskHistoryService {
	return &TaskHistoryService{
		historyRepo: historyRepo,
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

// Create creates a new task history record.
func (s *TaskHistoryService) Create(ctx context.Context, history *model.TaskHistory) error {
	if history == nil {
		return ErrInvalidTaskHistory
	}
	if history.TaskID <= 0 {
		return ErrInvalidTaskID
	}
	task, err := s.taskRepo.GetByID(ctx, history.TaskID)
	if err != nil {
		return fmt.Errorf("get task by id: %w", err)
	}
	if task == nil {
		return ErrTaskNotFound
	}
	if err := s.historyRepo.Create(ctx, history); err != nil {
		return fmt.Errorf("create task history: %w", err)
	}
	return nil
}

// GetByID returns a task history record by ID.
func (s *TaskHistoryService) GetByID(ctx context.Context, id int64) (*model.TaskHistory, error) {
	if id <= 0 {
		return nil, ErrInvalidTaskHistoryID
	}
	history, err := s.historyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get task history by id: %w", err)
	}
	if history == nil {
		return nil, ErrTaskHistoryNotFound
	}
	return history, nil
}

// GetByTaskID returns the history of a task.
func (s *TaskHistoryService) GetByTaskID(ctx context.Context, taskID int64) ([]*model.TaskHistory, error) {
	if taskID <= 0 {
		return nil, ErrInvalidTaskID
	}
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("get task by id: %w", err)
	}
	if task == nil {
		return nil, ErrTaskNotFound
	}
	history, err := s.historyRepo.ListByTaskID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("list task history by task id: %w", err)
	}
	return history, nil
}

// GetByProjectID returns task history for all tasks in a project.
func (s *TaskHistoryService) GetByProjectID(ctx context.Context, projectID int64) ([]*model.TaskHistory, error) {
	if projectID <= 0 {
		return nil, ErrInvalidProjectID
	}
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}
	history, err := s.historyRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list task history by project id: %w", err)
	}
	return history, nil
}
