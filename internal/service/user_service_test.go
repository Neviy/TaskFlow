package service

import (
	"context"
	"errors"
	"testing"

	"taskflow/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	CreateFunc     func(ctx context.Context, user *model.User) error
	GetByIDFunc    func(ctx context.Context, id int64) (*model.User, error)
	GetByEmailFunc func(ctx context.Context, email string) (*model.User, error)
	UpdateFunc     func(ctx context.Context, user *model.User) error
	DeleteFunc     func(ctx context.Context, id int64) error
}

func (m *MockUserRepository) Create(
	ctx context.Context,
	user *model.User,
) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) GetByID(
	ctx context.Context,
	id int64,
) (*model.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockUserRepository) GetByEmail(
	ctx context.Context,
	email string,
) (*model.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, nil
}

func (m *MockUserRepository) Update(
	ctx context.Context,
	user *model.User,
) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return nil
}

func (m *MockUserRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

// --------------------------------------------------
// Register
// --------------------------------------------------

func TestUserService_Register(t *testing.T) {
	repo := &MockUserRepository{
		GetByEmailFunc: func(
			ctx context.Context,
			email string,
		) (*model.User, error) {
			return nil, nil
		},
		CreateFunc: func(
			ctx context.Context,
			user *model.User,
		) error {
			user.ID = 1
			return nil
		},
	}
	us := NewUserService(repo)
	user, err := us.Register(
		context.Background(),
		"ivan",
		"ivan@example.com",
		"password123",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.ID != 1 {
		t.Fatalf("expected user ID 1, got %d", user.ID)
	}
	if user.PasswordHash == "password123" {
		t.Fatal("password was not hashed")
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte("password123"),
	); err != nil {
		t.Fatalf("invalid password hash: %v", err)
	}
}

func TestUserService_Register_UserAlreadyExists(t *testing.T) {
	repo := &MockUserRepository{
		GetByEmailFunc: func(
			ctx context.Context,
			email string,
		) (*model.User, error) {
			return &model.User{
				ID:    1,
				Email: email,
			}, nil
		},
	}
	us := NewUserService(repo)

	user, err := us.Register(
		context.Background(),
		"ivan",
		"ivan@example.com",
		"password123",
	)

	if user != nil {
		t.Fatal("expected nil user")
	}
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf(
			"expected ErrUserAlreadyExists, got %v",
			err,
		)
	}
}

// --------------------------------------------------
// Login
// --------------------------------------------------

func TestUserService_Login(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}
	repo := &MockUserRepository{
		GetByEmailFunc: func(
			ctx context.Context,
			email string,
		) (*model.User, error) {
			return &model.User{
				ID:           1,
				Username:     "ivan",
				Email:        email,
				PasswordHash: string(passwordHash),
			}, nil
		},
	}
	us := NewUserService(repo)

	token, err := us.Login(
		context.Background(),
		"ivan@example.com",
		"password123",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

func TestUserService_Login_UserNotFound(t *testing.T) {
	repo := &MockUserRepository{
		GetByEmailFunc: func(
			ctx context.Context,
			email string,
		) (*model.User, error) {
			return nil, nil
		},
	}
	us := NewUserService(repo)

	token, err := us.Login(
		context.Background(),
		"ivan@example.com",
		"password123",
	)

	if token != "" {
		t.Fatal("expected empty token")
	}
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf(
			"expected ErrUserNotFound, got %v",
			err,
		)
	}
}

func TestUserService_Login_InvalidPassword(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}
	repo := &MockUserRepository{
		GetByEmailFunc: func(
			ctx context.Context,
			email string,
		) (*model.User, error) {
			return &model.User{
				ID:           1,
				Username:     "ivan",
				Email:        email,
				PasswordHash: string(passwordHash),
			}, nil
		},
	}

	us := NewUserService(repo)
	token, err := us.Login(
		context.Background(),
		"ivan@example.com",
		"wrong-password",
	)

	if token != "" {
		t.Fatal("expected empty token")
	}
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf(
			"expected ErrInvalidCredentials, got %v",
			err,
		)
	}
}

// --------------------------------------------------
// GetByID
// --------------------------------------------------

func TestUserService_GetByID(t *testing.T) {
	repo := &MockUserRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.User, error) {
			return &model.User{
				ID:       id,
				Username: "ivan",
				Email:    "ivan@example.com",
			}, nil
		},
	}
	us := NewUserService(repo)
	user, err := us.GetByID(
		context.Background(),
		1,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user == nil {
		t.Fatal("expected user, got nil")
	}
	if user.ID != 1 {
		t.Fatalf("expected ID 1, got %d", user.ID)
	}
}

func TestUserService_GetByID_InvalidID(t *testing.T) {
	repo := &MockUserRepository{}
	us := NewUserService(repo)

	user, err := us.GetByID(
		context.Background(),
		0,
	)
	if user != nil {
		t.Fatal("expected nil user")
	}
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf(
			"expected ErrInvalidUserID, got %v",
			err,
		)
	}
}

func TestUserService_GetByID_UserNotFound(t *testing.T) {
	repo := &MockUserRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.User, error) {
			return nil, nil
		},
	}
	us := NewUserService(repo)
	user, err := us.GetByID(
		context.Background(),
		1,
	)
	if user != nil {
		t.Fatal("expected nil user")
	}
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf(
			"expected ErrUserNotFound, got %v",
			err,
		)
	}
}

// --------------------------------------------------
// Update
// --------------------------------------------------

func TestUserService_Update(t *testing.T) {
	repo := &MockUserRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.User, error) {
			return &model.User{
				ID:       id,
				Username: "old",
				Email:    "old@example.com",
			}, nil
		},
		GetByEmailFunc: func(
			ctx context.Context,
			email string,
		) (*model.User, error) {
			return nil, nil
		},
		UpdateFunc: func(
			ctx context.Context,
			user *model.User,
		) error {
			return nil
		},
	}
	us := NewUserService(repo)
	user := &model.User{
		ID:       1,
		Username: "new",
		Email:    "new@example.com",
	}
	err := us.Update(
		context.Background(),
		user,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestUserService_Update_UserNotFound(t *testing.T) {
	repo := &MockUserRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.User, error) {
			return nil, nil
		},
	}
	us := NewUserService(repo)
	user := &model.User{
		ID:       1,
		Username: "ivan",
		Email:    "ivan@example.com",
	}
	err := us.Update(
		context.Background(),
		user,
	)
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf(
			"expected ErrUserNotFound, got %v",
			err,
		)
	}
}

func TestUserService_Update_EmailAlreadyExists(t *testing.T) {
	repo := &MockUserRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.User, error) {
			return &model.User{
				ID: id,
			}, nil
		},
		GetByEmailFunc: func(
			ctx context.Context,
			email string,
		) (*model.User, error) {
			return &model.User{
				ID:    2,
				Email: email,
			}, nil
		},
	}
	us := NewUserService(repo)
	user := &model.User{
		ID:       1,
		Username: "ivan",
		Email:    "other@example.com",
	}
	err := us.Update(
		context.Background(),
		user,
	)
	if !errors.Is(err, ErrUserAlreadyExists) {
		t.Fatalf(
			"expected ErrUserAlreadyExists, got %v",
			err,
		)
	}
}

// --------------------------------------------------
// Delete
// --------------------------------------------------

func TestUserService_Delete(t *testing.T) {
	called := false
	repo := &MockUserRepository{
		DeleteFunc: func(
			ctx context.Context,
			id int64,
		) error {
			called = true
			return nil
		},
	}
	us := NewUserService(repo)
	err := us.Delete(
		context.Background(),
		1,
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !called {
		t.Fatal("expected Delete repository method to be called")
	}
}

func TestUserService_Delete_InvalidID(t *testing.T) {
	repo := &MockUserRepository{}
	us := NewUserService(repo)
	err := us.Delete(
		context.Background(),
		0,
	)
	if !errors.Is(err, ErrInvalidUserID) {
		t.Fatalf(
			"expected ErrInvalidUserID, got %v",
			err,
		)
	}
}

// --------------------------------------------------
// ChangePassword
// --------------------------------------------------

func TestUserService_ChangePassword(t *testing.T) {
	oldHash, err := bcrypt.GenerateFromPassword(
		[]byte("old-password"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf("generate old password hash: %v", err)
	}
	var updatedUser *model.User
	repo := &MockUserRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.User, error) {
			return &model.User{
				ID:           id,
				Username:     "ivan",
				Email:        "ivan@example.com",
				PasswordHash: string(oldHash),
			}, nil
		},
		UpdateFunc: func(
			ctx context.Context,
			user *model.User,
		) error {
			updatedUser = user
			return nil
		},
	}
	us := NewUserService(repo)
	err = us.ChangePassword(
		context.Background(),
		1,
		"old-password",
		"new-password",
	)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updatedUser == nil {
		t.Fatal("expected updated user")
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(updatedUser.PasswordHash),
		[]byte("new-password"),
	); err != nil {
		t.Fatalf("new password hash is invalid: %v", err)
	}
}

func TestUserService_ChangePassword_InvalidOldPassword(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte("old-password"),
		bcrypt.DefaultCost,
	)
	if err != nil {
		t.Fatalf("generate password hash: %v", err)
	}
	repo := &MockUserRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.User, error) {
			return &model.User{
				ID:           1,
				PasswordHash: string(passwordHash),
			}, nil
		},
	}
	us := NewUserService(repo)
	err = us.ChangePassword(
		context.Background(),
		1,
		"wrong-password",
		"new-password",
	)
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf(
			"expected ErrInvalidCredentials, got %v",
			err,
		)
	}
}

func TestUserService_ChangePassword_UserNotFound(t *testing.T) {
	repo := &MockUserRepository{
		GetByIDFunc: func(
			ctx context.Context,
			id int64,
		) (*model.User, error) {
			return nil, nil
		},
	}
	us := NewUserService(repo)
	err := us.ChangePassword(
		context.Background(),
		1,
		"old-password",
		"new-password",
	)
	if !errors.Is(err, ErrUserNotFound) {
		t.Fatalf(
			"expected ErrUserNotFound, got %v",
			err,
		)
	}
}
