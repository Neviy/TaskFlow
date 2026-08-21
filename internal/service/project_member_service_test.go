package service

import (
	"context"
	"errors"
	"taskflow/internal/model"
	"testing"
)

type MockMemberRepository struct {
	CreateFunc                func(ctx context.Context, member *model.ProjectMember) error
	GetByProjectAndUserIDFunc func(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error)
	ListByProjectIDFunc       func(ctx context.Context, projectID int64) ([]*model.ProjectMember, error)
	UpdateFunc                func(ctx context.Context, member *model.ProjectMember) error
	DeleteFunc                func(ctx context.Context, projectID, userID int64) error
}

func (m *MockMemberRepository) Create(ctx context.Context, member *model.ProjectMember) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, member)
	}
	return nil
}

func (m *MockMemberRepository) GetByProjectAndUserID(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error) {
	if m.GetByProjectAndUserIDFunc != nil {
		return m.GetByProjectAndUserIDFunc(ctx, projectID, userID)
	}
	return nil, nil
}

func (m *MockMemberRepository) ListByProjectID(ctx context.Context, projectID int64) ([]*model.ProjectMember, error) {
	if m.ListByProjectIDFunc != nil {
		return m.ListByProjectIDFunc(ctx, projectID)
	}
	return nil, nil
}

func (m *MockMemberRepository) Update(ctx context.Context, member *model.ProjectMember) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, member)
	}
	return nil
}

func (m *MockMemberRepository) Delete(ctx context.Context, projectID, userID int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, projectID, userID)
	}
	return nil
}

type MockMemberProjectRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*model.Project, error)
}

func (m *MockMemberProjectRepository) Create(ctx context.Context, project *model.Project) error {
	return nil
}

func (m *MockMemberProjectRepository) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockMemberProjectRepository) ListByOwnerID(ctx context.Context, ownerID int64) ([]*model.Project, error) {
	return nil, nil
}

func (m *MockMemberProjectRepository) Update(ctx context.Context, project *model.Project) error {
	return nil
}

func (m *MockMemberProjectRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

type MockMemberUserRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*model.User, error)
}

func (m *MockMemberUserRepository) Create(ctx context.Context, user *model.User) error {
	return nil
}

func (m *MockMemberUserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockMemberUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	return nil, nil
}

func (m *MockMemberUserRepository) Update(ctx context.Context, user *model.User) error {
	return nil
}

func (m *MockMemberUserRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func TestProjectMemberService_AddMember(t *testing.T) {
	memberRepo := &MockMemberRepository{
		GetByProjectAndUserIDFunc: func(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error) {
			return nil, nil
		},
		CreateFunc: func(ctx context.Context, member *model.ProjectMember) error {
			return nil
		},
	}

	projectRepo := &MockMemberProjectRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Project, error) {
			return &model.Project{ID: id}, nil
		},
	}

	userRepo := &MockMemberUserRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{ID: id}, nil
		},
	}

	pms := NewProjectMemberService(memberRepo, projectRepo, userRepo)

	err := pms.AddMember(
		context.Background(),
		1,
		2,
		model.RoleMember,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProjectMemberService_AddMember_AlreadyExists(t *testing.T) {
	memberRepo := &MockMemberRepository{
		GetByProjectAndUserIDFunc: func(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error) {
			return &model.ProjectMember{
				ProjectID: projectID,
				UserID:    userID,
			}, nil
		},
	}

	projectRepo := &MockMemberProjectRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Project, error) {
			return &model.Project{ID: id}, nil
		},
	}

	userRepo := &MockMemberUserRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.User, error) {
			return &model.User{ID: id}, nil
		},
	}

	pms := NewProjectMemberService(memberRepo, projectRepo, userRepo)

	err := pms.AddMember(
		context.Background(),
		1,
		2,
		model.RoleMember,
	)

	if !errors.Is(err, ErrProjectMemberAlreadyExists) {
		t.Fatalf("expected ErrProjectMemberAlreadyExists, got %v", err)
	}
}

func TestProjectMemberService_GetMember(t *testing.T) {
	repo := &MockMemberRepository{
		GetByProjectAndUserIDFunc: func(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error) {
			return &model.ProjectMember{
				ProjectID: projectID,
				UserID:    userID,
				Role:      model.RoleMember,
			}, nil
		},
	}

	pms := NewProjectMemberService(
		repo,
		&MockMemberProjectRepository{},
		&MockMemberUserRepository{},
	)

	member, err := pms.GetMember(context.Background(), 1, 2)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if member == nil {
		t.Fatal("expected member, got nil")
	}

	if member.UserID != 2 {
		t.Fatalf("expected user ID 2, got %d", member.UserID)
	}
}

func TestProjectMemberService_GetMember_NotFound(t *testing.T) {
	repo := &MockMemberRepository{
		GetByProjectAndUserIDFunc: func(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error) {
			return nil, nil
		},
	}

	pms := NewProjectMemberService(
		repo,
		&MockMemberProjectRepository{},
		&MockMemberUserRepository{},
	)

	_, err := pms.GetMember(context.Background(), 1, 2)

	if !errors.Is(err, ErrProjectMemberNotFound) {
		t.Fatalf("expected ErrProjectMemberNotFound, got %v", err)
	}
}

func TestProjectMemberService_ListMembers(t *testing.T) {
	repo := &MockMemberRepository{
		ListByProjectIDFunc: func(ctx context.Context, projectID int64) ([]*model.ProjectMember, error) {
			return []*model.ProjectMember{
				{ProjectID: projectID, UserID: 1, Role: model.RoleOwner},
				{ProjectID: projectID, UserID: 2, Role: model.RoleMember},
			}, nil
		},
	}

	projectRepo := &MockMemberProjectRepository{
		GetByIDFunc: func(ctx context.Context, id int64) (*model.Project, error) {
			return &model.Project{ID: id}, nil
		},
	}

	pms := NewProjectMemberService(
		repo,
		projectRepo,
		&MockMemberUserRepository{},
	)

	members, err := pms.ListMembers(context.Background(), 1)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(members))
	}
}

func TestProjectMemberService_UpdateRole(t *testing.T) {
	repo := &MockMemberRepository{
		GetByProjectAndUserIDFunc: func(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error) {
			return &model.ProjectMember{
				ProjectID: projectID,
				UserID:    userID,
				Role:      model.RoleMember,
			}, nil
		},
		UpdateFunc: func(ctx context.Context, member *model.ProjectMember) error {
			if member.Role != model.RoleAdmin {
				t.Fatalf("expected admin role, got %v", member.Role)
			}
			return nil
		},
	}

	pms := NewProjectMemberService(
		repo,
		&MockMemberProjectRepository{},
		&MockMemberUserRepository{},
	)

	err := pms.UpdateRole(
		context.Background(),
		1,
		2,
		model.RoleAdmin,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProjectMemberService_UpdateRole_CannotChangeOwner(t *testing.T) {
	repo := &MockMemberRepository{
		GetByProjectAndUserIDFunc: func(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error) {
			return &model.ProjectMember{
				ProjectID: projectID,
				UserID:    userID,
				Role:      model.RoleOwner,
			}, nil
		},
	}

	pms := NewProjectMemberService(
		repo,
		&MockMemberProjectRepository{},
		&MockMemberUserRepository{},
	)

	err := pms.UpdateRole(
		context.Background(),
		1,
		2,
		model.RoleAdmin,
	)

	if !errors.Is(err, ErrCannotChangeOwnerRole) {
		t.Fatalf("expected ErrCannotChangeOwnerRole, got %v", err)
	}
}

func TestProjectMemberService_RemoveMember(t *testing.T) {
	repo := &MockMemberRepository{
		GetByProjectAndUserIDFunc: func(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error) {
			return &model.ProjectMember{
				ProjectID: projectID,
				UserID:    userID,
				Role:      model.RoleMember,
			}, nil
		},
		DeleteFunc: func(ctx context.Context, projectID, userID int64) error {
			return nil
		},
	}

	pms := NewProjectMemberService(
		repo,
		&MockMemberProjectRepository{},
		&MockMemberUserRepository{},
	)

	err := pms.RemoveMember(
		context.Background(),
		1,
		2,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProjectMemberService_RemoveMember_CannotRemoveOwner(t *testing.T) {
	repo := &MockMemberRepository{
		GetByProjectAndUserIDFunc: func(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error) {
			return &model.ProjectMember{
				ProjectID: projectID,
				UserID:    userID,
				Role:      model.RoleOwner,
			}, nil
		},
	}

	pms := NewProjectMemberService(
		repo,
		&MockMemberProjectRepository{},
		&MockMemberUserRepository{},
	)

	err := pms.RemoveMember(
		context.Background(),
		1,
		2,
	)

	if !errors.Is(err, ErrCannotRemoveOwner) {
		t.Fatalf("expected ErrCannotRemoveOwner, got %v", err)
	}
}

func TestProjectMemberService_IsMember(t *testing.T) {
	repo := &MockMemberRepository{
		GetByProjectAndUserIDFunc: func(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error) {
			return &model.ProjectMember{
				ProjectID: projectID,
				UserID:    userID,
			}, nil
		},
	}

	pms := NewProjectMemberService(
		repo,
		&MockMemberProjectRepository{},
		&MockMemberUserRepository{},
	)

	isMember, err := pms.IsMember(
		context.Background(),
		1,
		2,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !isMember {
		t.Fatal("expected user to be a member")
	}
}

func TestProjectMemberService_IsMember_NotMember(t *testing.T) {
	repo := &MockMemberRepository{
		GetByProjectAndUserIDFunc: func(ctx context.Context, projectID, userID int64) (*model.ProjectMember, error) {
			return nil, nil
		},
	}

	pms := NewProjectMemberService(
		repo,
		&MockMemberProjectRepository{},
		&MockMemberUserRepository{},
	)

	isMember, err := pms.IsMember(
		context.Background(),
		1,
		2,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if isMember {
		t.Fatal("expected user not to be a member")
	}
}
