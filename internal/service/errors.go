// Package service provides the business logic for user management.
package service

import "errors"

// User errors.
var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidUser        = errors.New("invalid user")
)

// Project errors.
var (
	ErrInvalidProjectName = errors.New("invalid project name")
	ErrInvalidProjectID   = errors.New("invalid project ID")
	ErrProjectNotFound    = errors.New("project not found")
)

// Task errors.
var (
	ErrInvalidTaskTitle = errors.New("invalid task title")
	ErrInvalidTaskID    = errors.New("invalid task ID")
	ErrTaskNotFound     = errors.New("task not found")
)

// Project member errors.
var (
	ErrProjectMemberAlreadyExists = errors.New("project member already exists")
	ErrInvalidProjectMemberID     = errors.New("invalid project member ID")
	ErrProjectMemberNotFound      = errors.New("project member not found")
	ErrInvalidProjectRole         = errors.New("invalid project role")
	ErrCannotChangeOwnerRole      = errors.New("cannot change the role of the project owner")
	ErrCannotRemoveOwner          = errors.New("cannot remove the project owner")
)

// Comment errors.
var (
	ErrInvalidCommentID   = errors.New("invalid comment ID")
	ErrCommentNotFound    = errors.New("comment not found")
	ErrInvalidCommentText = errors.New("invalid comment text")
)

// Task history errors.
var (
	ErrInvalidTaskHistory   = errors.New("invalid task history")
	ErrInvalidTaskHistoryID = errors.New("invalid task history ID")
	ErrTaskHistoryNotFound  = errors.New("task history not found")
)
