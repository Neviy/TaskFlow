package repository

import (
	"context"
	"errors"
	"fmt"

	"taskflow/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository works with users in the database.
type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// Create adds a new user to the database.
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (username,email,password_hash)
						VALUES ($1, $2, $3)
						RETURNING id,created_at,updated_at`
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// GetByEmail finds a user by their email address.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id,username,email,password_hash,created_at,updated_at
						FROM users
						WHERE email = $1`
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}
	return user, nil
}

// GetByID finds a user by their ID.
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user := &model.User{}
	query := `SELECT id,username,email,password_hash,created_at,updated_at
						FROM users
						WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return user, nil
}

// Update changes an existing user's information in the database.
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	query := `UPDATE users
						SET username = $1,
								email = $2,
								password_hash = $3,
								updated_at = NOW()
						WHERE id = $4
						RETURNING updated_at`
	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.PasswordHash, user.ID).Scan(&user.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// Delete removes a user from the database.
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users
						WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
