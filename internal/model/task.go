package model

import (
	"time"
)

type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusReview     TaskStatus = "review"
	StatusDone       TaskStatus = "done"
	StatusOnHold     TaskStatus = "on_hold"
	StatusCanceled   TaskStatus = "canceled"
)

// TaskPriority represents the priority of a task.
type TaskPriority string
const (
	PriorityLow      TaskPriority = "low"
	PriorityMedium   TaskPriority = "medium"
	PriorityHigh     TaskPriority = "high"
	PriorityCritical TaskPriority = "critical"
)

// Task represents a task within a project.
type Task struct {
	ID          int64
	ProjectID   int64
	Title       string
	Description string
	Status      TaskStatus
	Priority    TaskPriority
	AssigneeID  *int64
	Deadline    *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
