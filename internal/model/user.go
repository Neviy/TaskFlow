// Package model contains domain models used by the application.
package model

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

// User represents an application user.
type User struct {
	ID           int64
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUser creates and returns a new User.
func NewUser(username, email, hash string) (*User, error) {
	var errs []error
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)
	if username == "" {
		errs = append(errs, errors.New("username cannot be empty"))
	}
	if _, err := mail.ParseAddress(email); err != nil {
		errs = append(errs, errors.New("invalid email"))
	}
	if strings.TrimSpace(hash) == "" {
		errs = append(errs, errors.New("password hash cannot be empty"))
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return &User{
		Username:     username,
		Email:        email,
		PasswordHash: hash,
	}, nil
}
