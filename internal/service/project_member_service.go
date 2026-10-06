package service

import (
	"context"
	"errors"
	"taskflow/internal/model"

	"github.com/jackc/pgx/v5"
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

// AddMember adds a user to a project.
func (pms *ProjectMemberService) AddMember(ctx context.Context, projectID int64, userID int64, role model.ProjectRole) error {
	if projectID <= 0 {
		return ErrInvalidProjectID
	}
	if userID <= 0 {
		return ErrInvalidUserID
	}
	_, err := pms.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProjectNotFound
		}
		return err
	}
	_, err = pms.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	member, err := pms.memberRepo.GetByProjectAndUserID(ctx, projectID, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
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
		return err
	}
	return nil
}

// GetMember returns a project member by project ID and user ID.
func (pms *ProjectMemberService) GetMember(ctx context.Context, projectID int64, userID int64) (*model.ProjectMember, error) {
	if projectID <= 0 {
		return nil, ErrInvalidProjectID
	}
	if userID <= 0 {
		return nil, ErrInvalidUserID
	}
	member, err := pms.memberRepo.GetByProjectAndUserID(ctx, projectID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectMemberNotFound
		}
		return nil, err
	}
	if member == nil {
		return nil, ErrProjectMemberNotFound
	}
	return member, nil
}

// ListMembers returns all members of a project.
func (pms *ProjectMemberService) ListMembers(ctx context.Context, projectID int64) ([]*model.ProjectMember, error) {
	if projectID <= 0 {
		return nil, ErrInvalidProjectID
	}
	_, err := pms.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	members, err := pms.memberRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return members, nil
}

// UpdateRole changes a project member's role.
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
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProjectMemberNotFound
		}
		return err
	}
	if member == nil {
		return ErrProjectMemberNotFound
	}
	if member.Role == model.RoleOwner {
		return ErrCannotChangeOwnerRole
	}
	member.Role = role
	if err := pms.memberRepo.Update(ctx, member); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProjectMemberNotFound
		}
		return err
	}
	return nil
}

// RemoveMember removes a user from a project.
func (pms *ProjectMemberService) RemoveMember(ctx context.Context, projectID, userID int64) error {
	if projectID <= 0 {
		return ErrInvalidProjectID
	}
	if userID <= 0 {
		return ErrInvalidUserID
	}
	member, err := pms.memberRepo.GetByProjectAndUserID(ctx, projectID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProjectMemberNotFound
		}
		return err
	}
	if member == nil {
		return ErrProjectMemberNotFound
	}
	if member.Role == model.RoleOwner {
		return ErrCannotRemoveOwner
	}
	if err := pms.memberRepo.Delete(ctx, projectID, userID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProjectMemberNotFound
		}
		return err
	}
	return nil
}

// IsMember checks whether a user belongs to a project.
func (pms *ProjectMemberService) IsMember(ctx context.Context, projectID int64, userID int64) (bool, error) {
	if projectID <= 0 {
		return false, ErrInvalidProjectID
	}
	if userID <= 0 {
		return false, ErrInvalidUserID
	}
	member, err := pms.memberRepo.GetByProjectAndUserID(ctx, projectID, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return member != nil, nil
}
