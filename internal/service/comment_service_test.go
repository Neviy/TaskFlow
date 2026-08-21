package service

import (
	"context"
	"errors"
	"taskflow/internal/model"
	"testing"
)

type MockCommentRepository struct {
	CreateFunc       func(ctx context.Context, comment *model.Comment) error
	GetByIDFunc      func(ctx context.Context, id int64) (*model.Comment, error)
	ListByTaskIDFunc func(ctx context.Context, taskID int64) ([]*model.Comment, error)
	DeleteFunc       func(ctx context.Context, id int64) error
}

func (m *MockCommentRepository) Create(
	ctx context.Context,
	comment *model.Comment,
) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, comment)
	}
	return nil
}

func (m *MockCommentRepository) GetByID(
	ctx context.Context,
	id int64,
) (*model.Comment, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockCommentRepository) ListByTaskID(
	ctx context.Context,
	taskID int64,
) ([]*model.Comment, error) {
	if m.ListByTaskIDFunc != nil {
		return m.ListByTaskIDFunc(ctx, taskID)
	}
	return nil, nil
}

func (m *MockCommentRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

type MockCommentTaskRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*model.Task, error)
}

func (m *MockCommentTaskRepository) Create(
	ctx context.Context,
	task *model.Task,
) error {
	return nil
}

func (m *MockCommentTaskRepository) GetByID(
	ctx context.Context,
	id int64,
) (*model.Task, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockCommentTaskRepository) ListByProjectID(
	ctx context.Context,
	projectID int64,
) ([]*model.Task, error) {
	return nil, nil
}

func (m *MockCommentTaskRepository) Update(
	ctx context.Context,
	task *model.Task,
) error {
	return nil
}

func (m *MockCommentTaskRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	return nil
}

type MockCommentUserRepository struct {
	GetByIDFunc func(ctx context.Context, id int64) (*model.User, error)
}

func (m *MockCommentUserRepository) Create(
	ctx context.Context,
	user *model.User,
) error {
	return nil
}

func (m *MockCommentUserRepository) GetByID(
	ctx context.Context,
	id int64,
) (*model.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockCommentUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	return nil, nil
}

func (m *MockCommentUserRepository) Update(
	ctx context.Context,
	user *model.User,
) error {
	return nil
}

func (m *MockCommentUserRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	return nil
}

func TestCommentService_Create(t *testing.T) {
	commentRepo := &MockCommentRepository{
		CreateFunc: func(
			ctx context.Context,
			comment *model.Comment,
		) error {
			comment.ID = 1
			return nil
		},
	}

	taskRepo := &MockCommentTaskRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Task, error) {
			return &model.Task{ID: id}, nil
		},
	}

	userRepo := &MockCommentUserRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.User, error) {
			return &model.User{ID: id}, nil
		},
	}

	cs := NewCommentService(
		commentRepo,
		taskRepo,
		userRepo,
	)

	comment, err := cs.Create(
		context.Background(),
		1,
		2,
		"Test comment",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if comment == nil {
		t.Fatal("expected comment, got nil")
	}

	if comment.ID != 1 {
		t.Fatalf("expected comment ID 1, got %d", comment.ID)
	}

	if comment.TaskID != 1 {
		t.Fatalf("expected task ID 1, got %d", comment.TaskID)
	}

	if comment.AuthorID != 2 {
		t.Fatalf("expected author ID 2, got %d", comment.AuthorID)
	}

	if comment.Text != "Test comment" {
		t.Fatalf("expected comment text 'Test comment', got %s", comment.Text)
	}
}

func TestCommentService_Create_TaskNotFound(t *testing.T) {
	taskRepo := &MockCommentTaskRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Task, error) {
			return nil, nil
		},
	}

	cs := NewCommentService(
		&MockCommentRepository{},
		taskRepo,
		&MockCommentUserRepository{},
	)

	_, err := cs.Create(
		context.Background(),
		1,
		2,
		"Test comment",
	)

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestCommentService_Create_UserNotFound(t *testing.T) {
	taskRepo := &MockCommentTaskRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Task, error) {
			return &model.Task{ID: id}, nil
		},
	}

	userRepo := &MockCommentUserRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.User, error) {
			return nil, nil
		},
	}

	cs := NewCommentService(
		&MockCommentRepository{},
		taskRepo,
		userRepo,
	)

	_, err := cs.Create(
		context.Background(),
		1,
		2,
		"Test comment",
	)

	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestCommentService_Create_InvalidText(t *testing.T) {
	cs := NewCommentService(
		&MockCommentRepository{},
		&MockCommentTaskRepository{},
		&MockCommentUserRepository{},
	)

	comment, err := cs.Create(
		context.Background(),
		1,
		2,
		"",
	)

	if comment != nil {
		t.Fatal("expected nil comment")
	}

	if !errors.Is(err, ErrInvalidCommentText) {
		t.Fatalf("expected ErrInvalidCommentText, got %v", err)
	}
}

func TestCommentService_GetByID(t *testing.T) {
	repo := &MockCommentRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Comment, error) {
			return &model.Comment{
				ID:       id,
				TaskID:   1,
				AuthorID: 2,
				Text:     "Test comment",
			}, nil
		},
	}

	cs := NewCommentService(
		repo,
		&MockCommentTaskRepository{},
		&MockCommentUserRepository{},
	)

	comment, err := cs.GetByID(
		context.Background(),
		1,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if comment == nil {
		t.Fatal("expected comment, got nil")
	}

	if comment.ID != 1 {
		t.Fatalf("expected ID 1, got %d", comment.ID)
	}
}

func TestCommentService_GetByID_NotFound(t *testing.T) {
	repo := &MockCommentRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Comment, error) {
			return nil, nil
		},
	}

	cs := NewCommentService(
		repo,
		&MockCommentTaskRepository{},
		&MockCommentUserRepository{},
	)

	_, err := cs.GetByID(
		context.Background(),
		1,
	)

	if !errors.Is(err, ErrCommentNotFound) {
		t.Fatalf("expected ErrCommentNotFound, got %v", err)
	}
}

func TestCommentService_GetByTaskID(t *testing.T) {
	commentRepo := &MockCommentRepository{
		ListByTaskIDFunc: func(
			ctx context.Context,
			taskID int64,
		) ([]*model.Comment, error) {
			return []*model.Comment{
				{
					ID:       1,
					TaskID:   taskID,
					AuthorID: 2,
					Text:     "First",
				},
				{
					ID:       2,
					TaskID:   taskID,
					AuthorID: 3,
					Text:     "Second",
				},
			}, nil
		},
	}

	taskRepo := &MockCommentTaskRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Task, error) {
			return &model.Task{ID: id}, nil
		},
	}

	cs := NewCommentService(
		commentRepo,
		taskRepo,
		&MockCommentUserRepository{},
	)

	comments, err := cs.GetByTaskID(
		context.Background(),
		1,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(comments) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(comments))
	}
}

func TestCommentService_GetByTaskID_TaskNotFound(t *testing.T) {
	taskRepo := &MockCommentTaskRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Task, error) {
			return nil, nil
		},
	}

	cs := NewCommentService(
		&MockCommentRepository{},
		taskRepo,
		&MockCommentUserRepository{},
	)

	_, err := cs.GetByTaskID(
		context.Background(),
		1,
	)

	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestCommentService_Delete(t *testing.T) {
	repo := &MockCommentRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Comment, error) {
			return &model.Comment{
				ID: id,
			}, nil
		},
		DeleteFunc: func(
			ctx context.Context,
			id int64,
		) error {
			return nil
		},
	}

	cs := NewCommentService(
		repo,
		&MockCommentTaskRepository{},
		&MockCommentUserRepository{},
	)

	err := cs.Delete(
		context.Background(),
		1,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCommentService_Delete_NotFound(t *testing.T) {
	repo := &MockCommentRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.Comment, error) {
			return nil, nil
		},
	}

	cs := NewCommentService(
		repo,
		&MockCommentTaskRepository{},
		&MockCommentUserRepository{},
	)

	err := cs.Delete(
		context.Background(),
		1,
	)

	if !errors.Is(err, ErrCommentNotFound) {
		t.Fatalf("expected ErrCommentNotFound, got %v", err)
	}
}
