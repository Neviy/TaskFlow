package service

import (
	"context"
	"taskflow/internal/model"
)

// UserRepository provides access to user storage.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id int64) error
}

// ProjectRepository provides access to project storage.
type ProjectRepository interface {
	Create(ctx context.Context, project *model.Project) error
	GetByID(ctx context.Context, id int64) (*model.Project, error)
	ListByOwnerID(ctx context.Context, ownerID int64) ([]*model.Project, error)
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, id int64) error
}

// ProjectMemberRepository provides access to project member storage.
type ProjectMemberRepository interface {
	Create(ctx context.Context, member *model.ProjectMember) error
	GetByProjectAndUserID(
		ctx context.Context,
		projectID int64,
		userID int64,
	) (*model.ProjectMember, error)
	ListByProjectID(
		ctx context.Context,
		projectID int64,
	) ([]*model.ProjectMember, error)
	Update(ctx context.Context, member *model.ProjectMember) error
	Delete(ctx context.Context, projectID int64, userID int64) error
}

// TaskRepository provides access to task storage.
type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id int64) (*model.Task, error)
	ListByProjectID(ctx context.Context, projectID int64) ([]*model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id int64) error
}

// CommentRepository provides access to comment storage.
type CommentRepository interface {
	Create(ctx context.Context, comment *model.Comment) error
	GetByID(ctx context.Context, id int64) (*model.Comment, error)
	ListByTaskID(ctx context.Context, taskID int64) ([]*model.Comment, error)
	Update(ctx context.Context, comment *model.Comment) error
	Delete(ctx context.Context, id int64) error
}

// TaskHistoryRepository provides access to task history storage.
type TaskHistoryRepository interface {
	Create(ctx context.Context, history *model.TaskHistory) error
	GetByID(ctx context.Context, id int64) (*model.TaskHistory, error)
	ListByTaskID(ctx context.Context, taskID int64) ([]*model.TaskHistory, error)
	ListByProjectID(ctx context.Context, projectID int64) ([]*model.TaskHistory, error)
}
