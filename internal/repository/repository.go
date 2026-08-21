package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	Comments       *CommentRepository
	Projects       *ProjectRepository
	ProjectMembers *ProjectMemberRepository
	Tasks          *TaskRepository
	TaskHistory    *TaskHistoryRepository
	Users          *UserRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Comments:       NewCommentRepository(db),
		Projects:       NewProjectRepository(db),
		ProjectMembers: NewProjectMemberRepository(db),
		Tasks:          NewTaskRepository(db),
		TaskHistory:    NewTaskHistoryRepository(db),
		Users:          NewUserRepository(db),
	}
}
