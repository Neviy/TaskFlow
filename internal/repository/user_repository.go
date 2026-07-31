// Package repository provides database access layer.
package repository

import (
	"context"
	"errors"
	"fmt"

	"taskflow/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository provides methods for working with users.
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create creates a new user in the database.
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	const query = `
		INSERT INTO users (
			username,
			email,
			password_hash
		)
		VALUES ($1, $2, $3)
		RETURNING
			id,
			created_at,
			updated_at
	`
	err := r.db.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// GetByEmail returns a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	const query = `
		SELECT
			id,
			username,
			email,
			password_hash,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

// GetByID returns a user by id.
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}
	const query = `
		SELECT
			id,
			username,
			email,
			password_hash,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

// Update updates an existing user in the database.
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	const query = `
		UPDATE users
		SET
			username = $1,
			email = $2,
			password_hash = $3,
			updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`
	err := r.db.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.ID,
	).Scan(
		&user.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return errors.New("user not found")
	}
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// Delete deletes a user from the database.
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM users
		WHERE id = $1
	`
	result, err := r.db.Exec(
		ctx,
		query,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}
	return nil
}
