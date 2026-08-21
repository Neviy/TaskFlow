package service

import (
	"context"
	"errors"
	"testing"

	"taskflow/internal/model"
)

type MockTaskRepository struct {
	CreateFunc          func(ctx context.Context, task *model.Task) error
	GetByIDFunc         func(ctx context.Context, id int64) (*model.Task, error)
	ListByProjectIDFunc func(ctx context.Context, projectID int64) ([]*model.Task, error)
	UpdateFunc          func(ctx context.Context, task *model.Task) error
	DeleteFunc          func(ctx context.Context, id int64) error
}

func (m *MockTaskRepository) Create(ctx context.Context, task *model.Task) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, task)
	}
	return nil
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id int64) (*model.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockTaskRepository) ListByProjectID(ctx context.Context, projectID int64) ([]*model.Task, error) {
	if m.ListByProjectIDFunc != nil {
		return m.ListByProjectIDFunc(ctx, projectID)
	}
	return nil, nil
}

func (m *MockTaskRepository) Update(ctx context.Context, task *model.Task) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, task)
	}
	return nil
}

func (m *MockTaskRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

type MockTaskProjectRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*model.Project, error)
}

func (m *MockTaskProjectRepository) Create(ctx context.Context, project *model.Project) error {
	return nil
}

func (m *MockTaskProjectRepository) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockTaskProjectRepository) ListByOwnerID(ctx context.Context, ownerID int64) ([]*model.Project, error) {
	return nil, nil
}

func (m *MockTaskProjectRepository) Update(ctx context.Context, project *model.Project) error {
	return nil
}

func (m *MockTaskProjectRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

type MockTaskUserRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*model.User, error)
}

func (m *MockTaskUserRepository) Create(ctx context.Context, user *model.User) error {
	return nil
}

func (m *MockTaskUserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockTaskUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}

func (m *MockTaskUserRepository) Update(ctx context.Context, user *model.User) error {
	return nil
}

func (m *MockTaskUserRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func TestTaskService_Create(t *testing.T) {
	taskRepo := &MockTaskRepository{
		CreateFunc: func(ctx context.Context, task *model.Task) error {
			task.ID = 1
			return nil
		},
	}

	projectRepo := &MockTaskProjectRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Project, error) {
			return &model.Project{ID: id}, nil
		},
	}

	userRepo := &MockTaskUserRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{ID: id}, nil
		},
	}

	ts := NewTaskService(taskRepo, projectRepo, userRepo)

	assigneeID := int64(5)

	task, err := ts.Create(
		context.Background(),
		"Test task",
		"Description",
		1,
		&assigneeID,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if task == nil {
		t.Fatal("expected task, got nil")
	}

	if task.ID != 1 {
		t.Fatalf("expected task ID 1, got %d", task.ID)
	}
}

func TestTaskService_Create_ProjectNotFound(t *testing.T) {
	projectRepo := &MockTaskProjectRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Project, error) {
			return nil, nil
		},
	}

	ts := NewTaskService(
		&MockTaskRepository{},
		projectRepo,
		&MockTaskUserRepository{},
	)

	task, err := ts.Create(
		context.Background(),
		"Test task",
		"Description",
		1,
		nil,
	)

	if task != nil {
		t.Fatal("expected nil task")
	}

	if !errors.Is(err, ErrProjectNotFound) {
		t.Fatalf("expected ErrProjectNotFound, got %v", err)
	}
}
