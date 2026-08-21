package service

import (
	"context"
	"errors"
	"taskflow/internal/model"
	"testing"
)

type MockTaskHistoryRepository struct {
	CreateFunc          func(ctx context.Context, history *model.TaskHistory) error
	GetByIDFunc         func(ctx context.Context, id int64) (*model.TaskHistory, error)
	ListByTaskIDFunc    func(ctx context.Context, taskID int64) ([]*model.TaskHistory, error)
	ListByProjectIDFunc func(ctx context.Context, projectID int64) ([]*model.TaskHistory, error)
}

func (m *MockTaskHistoryRepository) Create(ctx context.Context, history *model.TaskHistory) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, history)
	}
	return nil
}

func (m *MockTaskHistoryRepository) GetByID(ctx context.Context, id int64) (*model.TaskHistory, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockTaskHistoryRepository) ListByTaskID(ctx context.Context, taskID int64) ([]*model.TaskHistory, error) {
	if m.ListByTaskIDFunc != nil {
		return m.ListByTaskIDFunc(ctx, taskID)
	}
	return nil, nil
}

func (m *MockTaskHistoryRepository) ListByProjectID(ctx context.Context, projectID int64) ([]*model.TaskHistory, error) {
	if m.ListByProjectIDFunc != nil {
		return m.ListByProjectIDFunc(ctx, projectID)
	}
	return nil, nil
}

type MockHistoryTaskRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*model.Task, error)
}

func (m *MockHistoryTaskRepository) Create(ctx context.Context, task *model.Task) error {
	return nil
}

func (m *MockHistoryTaskRepository) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockHistoryTaskRepository) ListByProjectID(ctx context.Context, projectID int64) ([]*model.Task, error) {
	return nil, nil
}

func (m *MockHistoryTaskRepository) Update(ctx context.Context, task *model.Task) error {
	return nil
}

func (m *MockHistoryTaskRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

type MockHistoryProjectRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*model.Project, error)
}

func (m *MockHistoryProjectRepository) Create(ctx context.Context, project *model.Project) error {
	return nil
}

func (m *MockHistoryProjectRepository) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockHistoryProjectRepository) ListByOwnerID(ctx context.Context, ownerID int64) ([]*model.Project, error) {
	return nil, nil
}

func (m *MockHistoryProjectRepository) Update(ctx context.Context, project *model.Project) error {
	return nil
}

func (m *MockHistoryProjectRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func TestTaskHistoryService_Create(t *testing.T) {
	historyRepo := &MockTaskHistoryRepository{
		CreateFunc: func(ctx context.Context, history *model.TaskHistory) error {
			history.ID = 1
			return nil
		},
	}

	taskRepo := &MockHistoryTaskRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Task, error) {
			return &model.Task{ID: id}, nil
		},
	}

	ths := NewTaskHistoryService(
		historyRepo,
		taskRepo,
		&MockHistoryProjectRepository{},
	)

	history := &model.TaskHistory{
		TaskID:    1,
		ChangedBy: 2,
	}

	err := ths.Create(context.Background(), history)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if history.ID != 1 {
		t.Fatalf("expected history ID 1, got %d", history.ID)
	}
}

func TestTaskHistoryService_Create_InvalidHistory(t *testing.T) {
	ths := NewTaskHistoryService(
		&MockTaskHistoryRepository{},
		&MockHistoryTaskRepository{},
		&MockHistoryProjectRepository{},
	)

	err := ths.Create(context.Background(), nil)

	if !errors.Is(err, ErrInvalidTaskHistory) {
		t.Fatalf("expected ErrInvalidTaskHistory, got %v", err)
	}
}

func TestTaskHistoryService_Create_TaskNotFound(t *testing.T) {
	taskRepo := &MockHistoryTaskRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Task, error) {
			return nil, nil
		},
	}

	ths := NewTaskHistoryService(
		&MockTaskHistoryRepository{},
		taskRepo,
		&MockHistoryProjectRepository{},
	)

	history := &model.TaskHistory{
		TaskID:    1,
		ChangedBy: 2,
	}

	err := ths.Create(context.Background(), history)

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestTaskHistoryService_GetByID(t *testing.T) {
	historyRepo := &MockTaskHistoryRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.TaskHistory, error) {
			return &model.TaskHistory{
				ID:        id,
				TaskID:    1,
				ChangedBy: 2,
			}, nil
		},
	}

	ths := NewTaskHistoryService(
		historyRepo,
		&MockHistoryTaskRepository{},
		&MockHistoryProjectRepository{},
	)

	history, err := ths.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if history == nil {
		t.Fatal("expected history, got nil")
	}

	if history.ID != 1 {
		t.Fatalf("expected history ID 1, got %d", history.ID)
	}
}

func TestTaskHistoryService_GetByID_NotFound(t *testing.T) {
	historyRepo := &MockTaskHistoryRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.TaskHistory, error) {
			return nil, nil
		},
	}

	ths := NewTaskHistoryService(
		historyRepo,
		&MockHistoryTaskRepository{},
		&MockHistoryProjectRepository{},
	)

	_, err := ths.GetByID(context.Background(), 1)

	if !errors.Is(err, ErrTaskHistoryNotFound) {
		t.Fatalf("expected ErrTaskHistoryNotFound, got %v", err)
	}
}

func TestTaskHistoryService_GetByTaskID(t *testing.T) {
	historyRepo := &MockTaskHistoryRepository{
		ListByTaskIDFunc: func(ctx context.Context, taskID int64) ([]*model.TaskHistory, error) {
			return []*model.TaskHistory{
				{ID: 1, TaskID: taskID},
				{ID: 2, TaskID: taskID},
			}, nil
		},
	}

	taskRepo := &MockHistoryTaskRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Task, error) {
			return &model.Task{ID: id}, nil
		},
	}

	ths := NewTaskHistoryService(
		historyRepo,
		taskRepo,
		&MockHistoryProjectRepository{},
	)

	history, err := ths.GetByTaskID(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(history) != 2 {
		t.Fatalf("expected 2 history records, got %d", len(history))
	}
}

func TestTaskHistoryService_GetByTaskID_TaskNotFound(t *testing.T) {
	taskRepo := &MockHistoryTaskRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Task, error) {
			return nil, nil
		},
	}

	ths := NewTaskHistoryService(
		&MockTaskHistoryRepository{},
		taskRepo,
		&MockHistoryProjectRepository{},
	)

	_, err := ths.GetByTaskID(context.Background(), 1)

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestTaskHistoryService_GetByProjectID(t *testing.T) {
	historyRepo := &MockTaskHistoryRepository{
		ListByProjectIDFunc: func(ctx context.Context, projectID int64) ([]*model.TaskHistory, error) {
			return []*model.TaskHistory{
				{ID: 1, TaskID: 10},
				{ID: 2, TaskID: 20},
			}, nil
		},
	}

	projectRepo := &MockHistoryProjectRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Project, error) {
			return &model.Project{ID: id}, nil
		},
	}

	ths := NewTaskHistoryService(
		historyRepo,
		&MockHistoryTaskRepository{},
		projectRepo,
	)

	history, err := ths.GetByProjectID(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(history) != 2 {
		t.Fatalf("expected 2 history records, got %d", len(history))
	}
}

func TestTaskHistoryService_GetByProjectID_ProjectNotFound(t *testing.T) {
	projectRepo := &MockHistoryProjectRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Project, error) {
			return nil, nil
		},
	}

	ths := NewTaskHistoryService(
		&MockTaskHistoryRepository{},
		&MockHistoryTaskRepository{},
		projectRepo,
	)

	_, err := ths.GetByProjectID(context.Background(), 1)

	if !errors.Is(err, ErrProjectNotFound) {
		t.Fatalf("expected ErrProjectNotFound, got %v", err)
	}
}
