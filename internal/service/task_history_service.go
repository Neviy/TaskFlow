package service

import (
	"context"
	"errors"
	"taskflow/internal/model"

	"github.com/jackc/pgx/v5"
)

// TaskHistoryService handles business logic for task history.
type TaskHistoryService struct {
	historyRepo TaskHistoryRepository
	taskRepo    TaskRepository
	projectRepo ProjectRepository
}

func NewTaskHistoryService(historyRepo TaskHistoryRepository, taskRepo TaskRepository, projectRepo ProjectRepository) *TaskHistoryService {
	return &TaskHistoryService{
		historyRepo: historyRepo,
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

// Create adds a new history record for a task.
func (s *TaskHistoryService) Create(ctx context.Context, history *model.TaskHistory) error {
	if history == nil {
		return ErrInvalidTaskHistory
	}
	if history.TaskID <= 0 {
		return ErrInvalidTaskID
	}
	_, err := s.taskRepo.GetByID(ctx, history.TaskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrTaskNotFound
		}
		return err
	}
	if err := s.historyRepo.Create(ctx, history); err != nil {
		return err
	}
	return nil
}

// GetByID returns a history record by its ID.
func (s *TaskHistoryService) GetByID(ctx context.Context, id int64) (*model.TaskHistory, error) {
	if id <= 0 {
		return nil, ErrInvalidTaskHistoryID
	}
	history, err := s.historyRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskHistoryNotFound
		}
		return nil, err
	}
	return history, nil
}

// GetByTaskID returns all history records for a task.
func (s *TaskHistoryService) GetByTaskID(ctx context.Context, taskID int64) ([]*model.TaskHistory, error) {
	if taskID <= 0 {
		return nil, ErrInvalidTaskID
	}
	_, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	history, err := s.historyRepo.ListByTaskID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	return history, nil
}

// GetByProjectID returns history records for all tasks in a project.
func (s *TaskHistoryService) GetByProjectID(ctx context.Context, projectID int64) ([]*model.TaskHistory, error) {
	if projectID <= 0 {
		return nil, ErrInvalidProjectID
	}
	_, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	history, err := s.historyRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return history, nil
}
