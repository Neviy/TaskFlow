// Package model contains domain models used by the application.
package model

import "time"

// TaskHistory represents the history of task status changes.
type TaskHistory struct {
	ID        int64
	TaskID    int64
	ChangedBy int64
	OldStatus TaskStatus
	NewStatus TaskStatus
	CreatedAt time.Time
}
