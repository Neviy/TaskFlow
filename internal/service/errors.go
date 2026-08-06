// // Package service provides the business logic for user management.
package service

import "errors"

// Predefined errors for user service operations.
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrInvalidCredentials = errors.New("invalid email or password")
