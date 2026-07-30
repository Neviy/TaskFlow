// Package model contains domain models used by the application.
package model

import (
	"time"
)

// Project represents a project in the task management system.
type Project struct {
	ID          int64
	Name        string
	Description string
	OwnerID     int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
