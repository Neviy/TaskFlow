// Package service provides the business logic for user management.
package service

import (
	"context"
	"errors"
	"fmt"
	"taskflow/internal/model"
)

// ProjectService contains project business logic.
type ProjectService struct {
	projectRepo ProjectRepository
	memberRepo  ProjectMemberRepository
	userRepo    UserRepository
}

// NewProjectService creates a new ProjectService.
func NewProjectService(projectRepo ProjectRepository, memberRepo ProjectMemberRepository, userRepo UserRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		memberRepo:  memberRepo,
		userRepo:    userRepo,
	}
}

// Create creates a project and adds the owner as a member.
func (ps *ProjectService) Create(ctx context.Context, name, description string, ownerID int64) (*model.Project, error) {
	if ownerID <= 0 {
		return nil, ErrInvalidUserID
	}
	if name == "" {
		return nil, ErrInvalidProjectName
	}

	// Check owner exists.
	user, err := ps.userRepo.GetByID(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("get owner: %w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Create project.
	project := &model.Project{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
	}
	if err := ps.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}

	// Add owner to project members.
	member := &model.ProjectMember{
		ProjectID: project.ID,
		UserID:    ownerID,
		Role:      model.RoleOwner,
	}
	if err := ps.memberRepo.Create(ctx, member); err != nil {
		return nil, fmt.Errorf("create project owner: %w", err)
	}
	return project, nil
}

// GetByID returns a project by ID.
func (ps *ProjectService) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	if id <= 0 {
		return nil, ErrInvalidProjectID
	}
	project, err := ps.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}
	if project == nil {
		return nil, ErrProjectNotFound
	}
	return project, nil
}

// GetByOwner returns all projects owned by a user.
func (ps *ProjectService) GetByOwner(ctx context.Context, ownerID int64) ([]*model.Project, error) {
	if ownerID <= 0 {
		return nil, ErrInvalidUserID
	}
	projects, err := ps.projectRepo.ListByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list projects by owner id: %w", err)
	}
	return projects, nil
}

// Update updates project information.
func (ps *ProjectService) Update(ctx context.Context, project *model.Project) error {
	if project == nil {
		return errors.New("project is nil")
	}
	if project.Name == "" {
		return ErrInvalidProjectName
	}
	existingProject, err := ps.projectRepo.GetByID(ctx, project.ID)
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}
	if existingProject == nil {
		return ErrProjectNotFound
	}
	if err := ps.projectRepo.Update(ctx, project); err != nil {
		return fmt.Errorf("update project: %w", err)
	}
	return nil
}

// Delete removes a project.
func (ps *ProjectService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidProjectID
	}
	project, err := ps.projectRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}
	if project == nil {
		return ErrProjectNotFound
	}
	if err := ps.projectRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete project: %w", err)
	}
	return nil
}
