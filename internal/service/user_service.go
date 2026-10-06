package service

import (
	"context"
	"errors"
	"fmt"
	"taskflow/internal/auth"
	"taskflow/internal/model"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService contains the business logic for user management.
type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register creates a new user after checking that the email is available.
func (us *UserService) Register(ctx context.Context, username, email, password string) (*model.User, error) {
	user, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
		user = nil
	}
	if user != nil {
		return nil, ErrUserAlreadyExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user, err = model.NewUser(username, email, string(hash))
	if err != nil {
		return nil, fmt.Errorf("create user model: %w", err)
	}
	if err := us.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// Login checks the user's credentials and returns a JWT token.
func (us *UserService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := us.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", err
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

// GetByID returns a user by their ID.
func (us *UserService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	if id <= 0 {
		return nil, ErrInvalidUserID
	}
	user, err := us.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// Update changes the user's username, email, or password.
func (us *UserService) Update(ctx context.Context, user *model.User) error {
	if user == nil {
		return ErrInvalidUser
	}
	if user.ID <= 0 {
		return ErrInvalidUserID
	}
	if user.Username == "" || user.Email == "" {
		return ErrInvalidUser
	}
	_, err := us.repo.GetByID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	existingByEmail, err := us.repo.GetByEmail(ctx, user.Email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		existingByEmail = nil
	}
	if existingByEmail != nil && existingByEmail.ID != user.ID {
		return ErrUserAlreadyExists
	}
	if err := us.repo.Update(ctx, user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

// Delete removes a user by their ID.
func (us *UserService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidUserID
	}
	if err := us.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

// ChangePassword verifies the old password and saves a new one.
func (us *UserService) ChangePassword(ctx context.Context,
	userID int64, oldPassword string, newPassword string,
) error {
	if userID <= 0 {
		return ErrInvalidUserID
	}
	if oldPassword == "" || newPassword == "" {
		return ErrInvalidCredentials
	}
	user, err := us.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
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
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}
