package service

import (
	"context"
	"errors"
	"taskflow/internal/model"
	"testing"
)

type MockProjectRepository struct {
	CreateFunc        func(ctx context.Context, project *model.Project) error
	GetByIDFunc       func(ctx context.Context, id int64) (*model.Project, error)
	ListByOwnerIDFunc func(ctx context.Context, ownerID int64) ([]*model.Project, error)
	UpdateFunc        func(ctx context.Context, project *model.Project) error
	DeleteFunc        func(ctx context.Context, id int64) error
}

func (m *MockProjectRepository) Create(ctx context.Context, project *model.Project) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, project)
	}
	return nil
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockProjectRepository) ListByOwnerID(ctx context.Context, ownerID int64) ([]*model.Project, error) {
	if m.ListByOwnerIDFunc != nil {
		return m.ListByOwnerIDFunc(ctx, ownerID)
	}
	return nil, nil
}

func (m *MockProjectRepository) Update(ctx context.Context, project *model.Project) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, project)
	}
	return nil
}

func (m *MockProjectRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

type MockProjectMemberRepository struct {
	CreateFunc                func(ctx context.Context, member *model.ProjectMember) error
	GetByProjectAndUserIDFunc func(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error)
	ListByProjectIDFunc       func(ctx context.Context, projectID int64) ([]*model.ProjectMember, error)
	UpdateFunc                func(ctx context.Context, member *model.ProjectMember) error
	DeleteFunc                func(ctx context.Context, projectID, userID int64) error
}

func (m *MockProjectMemberRepository) Create(ctx context.Context, member *model.ProjectMember) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, member)
	}
	return nil
}

func (m *MockProjectMemberRepository) GetByProjectAndUserID(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error) {
	if m.GetByProjectAndUserIDFunc != nil {
		return m.GetByProjectAndUserIDFunc(ctx, projectID, userID)
	}
	return nil, nil
}

func (m *MockProjectMemberRepository) ListByProjectID(ctx context.Context, projectID int64) ([]*model.ProjectMember, error) {
	if m.ListByProjectIDFunc != nil {
		return m.ListByProjectIDFunc(ctx, projectID)
	}
	return nil, nil
}

func (m *MockProjectMemberRepository) Update(ctx context.Context, member *model.ProjectMember) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, member)
	}
	return nil
}

func (m *MockProjectMemberRepository) Delete(ctx context.Context, projectID, userID int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, projectID, userID)
	}
	return nil
}

type MockProjectUserRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*model.User, error)
}

func (m *MockProjectUserRepository) Create(ctx context.Context, user *model.User) error {
	return nil
}

func (m *MockProjectUserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockProjectUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}

func (m *MockProjectUserRepository) Update(ctx context.Context, user *model.User) error {
	return nil
}

func (m *MockProjectUserRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func TestProjectService_Create(t *testing.T) {
	projectRepo := &MockProjectRepository{
		CreateFunc: func(ctx context.Context, project *model.Project) error {
			project.ID = 1
			return nil
		},
	}

	memberRepo := &MockProjectMemberRepository{
		CreateFunc: func(ctx context.Context, member *model.ProjectMember) error {
			return nil
		},
	}

	userRepo := &MockProjectUserRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{ID: id}, nil
		},
	}

	ps := NewProjectService(projectRepo, memberRepo, userRepo)

	project, err := ps.Create(context.Background(), "TaskFlow", "Project", 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if project == nil {
		t.Fatal("expected project, got nil")
	}

	if project.ID != 1 {
		t.Fatalf("expected project ID 1, got %d", project.ID)
	}

	if project.OwnerID != 1 {
		t.Fatalf("expected owner ID 1, got %d", project.OwnerID)
	}
}

func TestProjectService_Create_InvalidOwner(t *testing.T) {
	ps := NewProjectService(
		&MockProjectRepository{},
		&MockProjectMemberRepository{},
		&MockProjectUserRepository{},
	)

	project, err := ps.Create(context.Background(), "TaskFlow", "Project", 0)

	if project != nil {
		t.Fatal("expected nil project")
	}

	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
}

func TestProjectService_Create_UserNotFound(t *testing.T) {
	userRepo := &MockProjectUserRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return nil, nil
		},
	}

	ps := NewProjectService(
		&MockProjectRepository{},
		&MockProjectMemberRepository{},
		userRepo,
	)

	_, err := ps.Create(context.Background(), "TaskFlow", "Project", 1)

	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestProjectService_GetByID(t *testing.T) {
	repo := &MockProjectRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Project, error) {
			return &model.Project{
				ID:   id,
				Name: "TaskFlow",
			}, nil
		},
	}

	ps := NewProjectService(
		repo,
		&MockProjectMemberRepository{},
		&MockProjectUserRepository{},
	)

	project, err := ps.GetByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if project == nil {
		t.Fatal("expected project, got nil")
	}

	if project.ID != 1 {
		t.Fatalf("expected ID 1, got %d", project.ID)
	}
}

func TestProjectService_GetByID_NotFound(t *testing.T) {
	repo := &MockProjectRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Project, error) {
			return nil, nil
		},
	}

	ps := NewProjectService(
		repo,
		&MockProjectMemberRepository{},
		&MockProjectUserRepository{},
	)

	_, err := ps.GetByID(context.Background(), 1)

	if !errors.Is(err, ErrProjectNotFound) {
		t.Fatalf("expected ErrProjectNotFound, got %v", err)
	}
}

func TestProjectService_GetByOwner(t *testing.T) {
	repo := &MockProjectRepository{
		ListByOwnerIDFunc: func(ctx context.Context, ownerID int64) ([]*model.Project, error) {
			return []*model.Project{
				{ID: 1, Name: "Project 1", OwnerID: ownerID},
				{ID: 2, Name: "Project 2", OwnerID: ownerID},
			}, nil
		},
	}

	ps := NewProjectService(
		repo,
		&MockProjectMemberRepository{},
		&MockProjectUserRepository{},
	)

	projects, err := ps.GetByOwner(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(projects) != 2 {
		t.Fatalf("expected 2 projects, got %d", len(projects))
	}

	if projects[0].OwnerID != 1 {
		t.Fatalf("expected owner ID 1, got %d", projects[0].OwnerID)
	}
}

func TestProjectService_GetByOwner_InvalidOwner(t *testing.T) {
	ps := NewProjectService(
		&MockProjectRepository{},
		&MockProjectMemberRepository{},
		&MockProjectUserRepository{},
	)

	projects, err := ps.GetByOwner(context.Background(), 0)

	if projects != nil {
		t.Fatal("expected nil projects")
	}

	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf("expected ErrInvalidUserID, got %v", err)
	}
}

func TestProjectService_Update(t *testing.T) {
	repo := &MockProjectRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Project, error) {
			return &model.Project{
				ID:      id,
				Name:    "Old name",
				OwnerID: 1,
			}, nil
		},
		UpdateFunc: func(ctx context.Context, project *model.Project) error {
			return nil
		},
	}

	ps := NewProjectService(
		repo,
		&MockProjectMemberRepository{},
		&MockProjectUserRepository{},
	)

	project := &model.Project{
		ID:          1,
		Name:        "New name",
		Description: "Updated description",
		OwnerID:     1,
	}

	err := ps.Update(context.Background(), project)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProjectService_Update_NotFound(t *testing.T) {
	repo := &MockProjectRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Project, error) {
			return nil, nil
		},
	}

	ps := NewProjectService(
		repo,
		&MockProjectMemberRepository{},
		&MockProjectUserRepository{},
	)

	project := &model.Project{
		ID:   1,
		Name: "TaskFlow",
	}

	err := ps.Update(context.Background(), project)

	if !errors.Is(err, ErrProjectNotFound) {
		t.Fatalf("expected ErrProjectNotFound, got %v", err)
	}
}

func TestProjectService_Update_InvalidName(t *testing.T) {
	ps := NewProjectService(
		&MockProjectRepository{},
		&MockProjectMemberRepository{},
		&MockProjectUserRepository{},
	)

	project := &model.Project{
		ID:   1,
		Name: "",
	}

	err := ps.Update(context.Background(), project)

	if !errors.Is(err, ErrInvalidProjectName) {
		t.Fatalf("expected ErrInvalidProjectName, got %v", err)
	}
}

func TestProjectService_Delete(t *testing.T) {
	repo := &MockProjectRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Project, error) {
			return &model.Project{ID: id}, nil
		},
		DeleteFunc: func(ctx context.Context, id int64) error {
			return nil
		},
	}

	ps := NewProjectService(
		repo,
		&MockProjectMemberRepository{},
		&MockProjectUserRepository{},
	)
	err := ps.Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
