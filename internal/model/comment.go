package model

import (
	"time"
)

// Comment represents a comment left on a task.
type Comment struct {
	ID        int64
	TaskID    int64
	AuthorID  int64
	Text      string
	CreatedAt time.Time
}
