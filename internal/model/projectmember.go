// Package model contains domain models used by the application.
package model

import (
	"time"
)

// ProjectRole represents a user's role in a project.
type ProjectRole string

const (
	RoleOwner  ProjectRole = "owner"
	RoleAdmin  ProjectRole = "admin"
	RoleMember ProjectRole = "member"
)

// ProjectMember represents a user's membership in a project.
type ProjectMember struct {
	ID        int64
	ProjectID int64
	UserID    int64
	Role      ProjectRole
	CreatedAt time.Time
}
