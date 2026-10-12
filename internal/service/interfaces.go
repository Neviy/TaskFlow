package service

import (
	"context"
	"taskflow/internal/model"
)

// UserRepository works with users.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

// ProjectRepository works with projects.
type ProjectRepository interface {
	Create(ctx context.Context, project *model.Project) error
	GetByID(ctx context.Context, id int64) (*model.Project, error)
	ListByOwnerID(ctx context.Context, ownerID int64) ([]*model.Project, error)
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id int64) error
}

// ProjectMemberRepository manages project members.
type ProjectMemberRepository interface {
	Create(ctx context.Context, member *model.ProjectMember) error
	GetByProjectAndUserID(ctx context.Context, projectID int64, userID int64) (*model.ProjectMember, error)
	ListByProjectID(ctx context.Context, projectID int64) ([]*model.ProjectMember, error)
	Update(ctx context.Context, member *model.ProjectMember) error
	Delete(ctx context.Context, projectID, userID int64) error
}

// TaskRepository works with tasks.
type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id int64) (*model.Task, error)
	ListByProjectID(ctx context.Context, projectID int64) ([]*model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id int64) error
}

// CommentRepository works with comments.
type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByID(ctx context.Context, id int64) (*model.Comment, error)
	ListByTaskID(ctx context.Context, taskID int64) ([]*model.Comment, error)
	Delete(ctx context.Context, id int64) error
}

// TaskHistoryRepository works with task history.
type TaskHistoryRepository interface {
	Create(ctx context.Context, history *model.TaskHistory) error
	GetByID(ctx context.Context, id int64) (*model.TaskHistory, error)
	ListByTaskID(ctx context.Context, taskID int64) ([]*model.TaskHistory, error)
	ListByProjectID(ctx context.Context, projectID int64) ([]*model.TaskHistory, error)
}
