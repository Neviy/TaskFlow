package service

import (
	"context"
	"errors"

	"taskflow/internal/model"

	"github.com/jackc/pgx/v5"
)

// TaskHistoryService handles task history.
type TaskHistoryService struct {
	historyRepo TaskHistoryRepository
	taskRepo    TaskRepository
	projectRepo ProjectRepository
	memberRepo  ProjectMemberRepository
}

func NewTaskHistoryService(historyRepo TaskHistoryRepository, taskRepo TaskRepository, projectRepo ProjectRepository, memberRepo ProjectMemberRepository) *TaskHistoryService {
	return &TaskHistoryService{
		historyRepo: historyRepo,
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		memberRepo:  memberRepo,
	}
}

func (s *TaskHistoryService) checkProjectAccess(ctx context.Context, projectID, userID int64) error {
	if projectID <= 0 || userID <= 0 {
		return ErrForbidden
	}
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProjectNotFound
		}
		return err
	}
	if project.OwnerID == userID {
		return nil
	}
	_, err = s.memberRepo.GetByProjectAndUserID(ctx, projectID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrForbidden
		}
		return err
	}
	return nil
}

// Create adds a history record.
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

// GetByID returns a history record.
func (s *TaskHistoryService) GetByID(ctx context.Context, id, userID int64) (*model.TaskHistory, error) {
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
	task, err := s.taskRepo.GetByID(ctx, history.TaskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	if err := s.checkProjectAccess(ctx, task.ProjectID, userID); err != nil {
		return nil, err
	}
	return history, nil
}

// GetByTaskID returns task history.
func (s *TaskHistoryService) GetByTaskID(ctx context.Context, taskID, userID int64) ([]*model.TaskHistory, error) {
	if taskID <= 0 {
		return nil, ErrInvalidTaskID
	}
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	if err := s.checkProjectAccess(ctx, task.ProjectID, userID); err != nil {
		return nil, err
	}
	return s.historyRepo.ListByTaskID(ctx, taskID)
}

// GetByProjectID returns project task history.
func (s *TaskHistoryService) GetByProjectID(ctx context.Context, projectID, userID int64) ([]*model.TaskHistory, error) {
	if projectID <= 0 {
		return nil, ErrInvalidProjectID
	}
	if err := s.checkProjectAccess(ctx, projectID, userID); err != nil {
		return nil, err
	}
	return s.historyRepo.ListByProjectID(ctx, projectID)
}
