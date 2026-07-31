// Package repository provides implementations for working with the database.
package repository

import "github.com/jackc/pgx/v5/pgxpool"

// Repository stores shared dependencies for repository implementations.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates and returns a new Repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}
