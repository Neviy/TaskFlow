package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"taskflow/internal/model"

	"github.com/jackc/pgx/v5"
)

// ProjectService handles project-related business logic.
type ProjectService struct {
	projectRepo ProjectRepository
	memberRepo  ProjectMemberRepository
	userRepo    UserRepository
}

func NewProjectService(projectRepo ProjectRepository, memberRepo ProjectMemberRepository, userRepo UserRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		memberRepo:  memberRepo,
		userRepo:    userRepo,
	}
}

// Create creates a project and adds its owner as a member.
func (ps *ProjectService) Create(ctx context.Context, name, description string, ownerID int64) (*model.Project, error) {
	if ownerID <= 0 {
		return nil, ErrInvalidUserID
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrInvalidProjectName
	}
	_, err := ps.userRepo.GetByID(ctx, ownerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", ErrUserNotFound, err)
		}
		return nil, fmt.Errorf("get owner: %w", err)
	}
	project := &model.Project{
		Name:        name,
		Description: strings.TrimSpace(description),
		OwnerID:     ownerID,
	}
	if err := ps.projectRepo.Create(ctx, project); err != nil {
		return nil, fmt.Errorf("create project: %w", err)
	}
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

// GetByID returns a project by its ID.
func (ps *ProjectService) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	if id <= 0 {
		return nil, ErrInvalidProjectID
	}
	project, err := ps.projectRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", ErrProjectNotFound, err)
		}
		return nil, fmt.Errorf("get project by id: %w", err)
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

// Update changes a project's name and description.
func (ps *ProjectService) Update(ctx context.Context, project *model.Project) error {
	if project == nil {
		return ErrInvalidProject
	}
	if project.ID <= 0 {
		return ErrInvalidProjectID
	}
	project.Name = strings.TrimSpace(project.Name)
	if project.Name == "" {
		return ErrInvalidProjectName
	}
	_, err := ps.projectRepo.GetByID(ctx, project.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: %v", ErrProjectNotFound, err)
		}
		return fmt.Errorf("get project by id: %w", err)
	}
	if err := ps.projectRepo.Update(ctx, project); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: %v", ErrProjectNotFound, err)
		}
		return fmt.Errorf("update project: %w", err)
	}
	return nil
}

// Delete removes a project.
func (ps *ProjectService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidProjectID
	}
	_, err := ps.projectRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: %v", ErrProjectNotFound, err)
		}
		return fmt.Errorf("get project by id: %w", err)
	}
	if err := ps.projectRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) || err.Error() == "project not found" {
			return fmt.Errorf("%w: %v", ErrProjectNotFound, err)
		}
		return fmt.Errorf("delete project: %w", err)
	}
	return nil
}
