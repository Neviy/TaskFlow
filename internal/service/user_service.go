// Package service provides the business logic for user management.
package service

import (
	"context"
	"errors"
	"fmt"
	"taskflow/internal/auth"
	"taskflow/internal/model"

	"golang.org/x/crypto/bcrypt"
)

// UserService provides methods for user management.
type UserService struct {
	repo UserRepository
}

// NewUserService creates a new instance of UserService with the provided UserRepository.
func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register creates a new user if the email is not already in use.
func (us *UserService) Register(ctx context.Context, username, email, password string) (*model.User, error) {
	user, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	if user != nil {
		return nil, ErrUserAlreadyExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	users, err := model.NewUser(username, email, string(hash))
	if err != nil {
		return nil, fmt.Errorf("create user model: %w", err)
	}
	if err := us.repo.Create(ctx, users); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return users, nil
}

// Login authenticates a user and returns a JWT token.
func (us *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("get user by email: %w", err)
	}
	if user == nil {
		return "", ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", ErrInvalidCredentials
	}
	tokenString, err := auth.GenerateToken(user.ID)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return tokenString, nil
}

// GetByID retrieves a user by their ID.
func (us *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid id")
	}
	user, err := us.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user by id:%w", err)
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// Update modifies an existing user's information.
func (us *UserService) Update(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	existing, err := us.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existing != nil {
		return ErrUserAlreadyExists
	}
	if err := us.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// Delete removes a user by ID.
func (us *UserService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	if err := us.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

func (us *UserService) ChangePassword(ctx context.Context,
	userID int64, oldPassword string, newPassword string,
) error {
	if userID <= 0 {
		return ErrInvalidUserID
	}
	if oldPassword == "" || newPassword == "" {
		return errors.New("passwords cannot be empty")
	}
	user, err := us.repo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}
	if user == nil {
		return ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(oldPassword),
	); err != nil {
		return ErrInvalidCredentials
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash new password: %w", err)
	}
	user.PasswordHash = string(newHash)
	if err := us.repo.Update(ctx, user); err != nil {
		return fmt.Errorf("update user password: %w", err)
	}
	return nil
}
