package service

import (
	"context"
	"fmt"
	"taskflow/internal/model"
)

type ProjectMemberService struct {
	memberRepo  ProjectMemberRepository
	projectRepo ProjectRepository
	userRepo    UserRepository
}

func NewProjectMemberService(memberRepo ProjectMemberRepository, projectRepo ProjectRepository, userRepo UserRepository) *ProjectMemberService {
	return &ProjectMemberService{
		memberRepo:  memberRepo,
		projectRepo: projectRepo,
		userRepo:    userRepo,
	}
}

func (pms *ProjectMemberService) AddMember(ctx context.Context, projectID int64, userID int64, role model.ProjectRole) error {
	if projectID <= 0 {
		return ErrInvalidProjectID
	}
	if userID <= 0 {
		return ErrInvalidUserID
	}
	existingProject, err := pms.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return fmt.Errorf("get project by id: %w", err)
	}
	if existingProject == nil {
		return ErrProjectNotFound
	}
	existingUser, err := pms.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}
	if existingUser == nil {
		return ErrUserNotFound
	}
	member, err := pms.memberRepo.GetByProjectAndUserID(ctx, projectID, userID)
	if err != nil {
		return fmt.Errorf("get project member: %w", err)
	}
	if member != nil {
		return ErrProjectMemberAlreadyExists
	}
	member = &model.ProjectMember{
		ProjectID: projectID,
		UserID:    userID,
		Role:      role,
	}
	if err := pms.memberRepo.Create(ctx, member); err != nil {
		return fmt.Errorf("create project member: %w", err)
	}
	return nil
}

func (pms *ProjectMemberService) GetMember(ctx context.Context, projectID int64, userID int64) (*model.ProjectMember, error) {
	if projectID <= 0 {
		return nil, ErrInvalidProjectID
	}
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	member, err := pms.memberRepo.GetByProjectAndUserID(ctx, projectID, userID)
	if err != nil {
		return nil, fmt.Errorf("get project member: %w", err)
	}
	if member == nil {
		return nil, ErrProjectMemberNotFound
	}
	return member, nil
}

func (pms *ProjectMemberService) ListMembers(ctx context.Context, projectID int64) ([]*model.ProjectMember, error) {
	if projectID <= 0 {
		return nil, ErrInvalidProjectID
	}
	existingProject, err := pms.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("get project by id: %w", err)
	}
	if existingProject == nil {
		return nil, ErrProjectNotFound
	}
	members, err := pms.memberRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list project members: %w", err)
	}
	return members, nil
}

func (pms *ProjectMemberService) UpdateRole(ctx context.Context, projectID int64, userID int64, role model.ProjectRole,
) error {
	if projectID <= 0 {
		return ErrInvalidProjectID
	}
	if userID <= 0 {
		return ErrInvalidUserID
	}
	if role == "" {
		return ErrInvalidProjectRole
	}
	member, err := pms.memberRepo.GetByProjectAndUserID(ctx, projectID, userID)
	if err != nil {
		return fmt.Errorf("get project member: %w", err)
	}
	if member == nil {
		return ErrProjectMemberNotFound
	}
	if member.Role == model.RoleOwner {
		return ErrCannotChangeOwnerRole
	}
	member.Role = role
	if err := pms.memberRepo.Update(ctx, member); err != nil {
		return fmt.Errorf("update project member role: %w", err)
	}
	return nil
}

func (pms *ProjectMemberService) RemoveMember(ctx context.Context, projectID int64, userID int64) error {
	if projectID <= 0 {
		return ErrInvalidProjectID
	}
	if userID <= 0 {
		return ErrInvalidUserID
	}
	member, err := pms.memberRepo.GetByProjectAndUserID(ctx, projectID, userID)
	if err != nil {
		return fmt.Errorf("get project member: %w", err)
	}
	if member == nil {
		return ErrProjectMemberNotFound
	}
	if member.Role == model.RoleOwner {
		return ErrCannotChangeOwnerRole
	}
	if err := pms.memberRepo.Delete(ctx, projectID, userID); err != nil {
		return fmt.Errorf("delete project member: %w", err)
	}
	return nil
}

func (pms *ProjectMemberService) IsMember(ctx context.Context, projectID int64, userID int64,
) (bool, error) {
	if projectID <= 0 {
		return false, ErrInvalidProjectID
	}
	if userID <= 0 {
		return false, ErrInvalidUserID
	}
	member, err := pms.memberRepo.GetByProjectAndUserID(ctx, projectID, userID)
	if err != nil {
		return false, fmt.Errorf("get project member: %w", err)
	}
	if member == nil {
		return false, nil
	}
	return true, nil
}
