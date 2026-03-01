package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EduCoelhoTs/nba-predict-api/internal/adapter/repository/postgres/sqlc"
	"github.com/EduCoelhoTs/nba-predict-api/internal/core/domain"
	"github.com/google/uuid"
)

// UserQueriesPort define a interface para queries do usuário (Port pattern)
type UserQueriesPort interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) error
	GetAllUsers(ctx context.Context) ([]sqlc.AuthUser, error)
	GetUserById(ctx context.Context, id uuid.UUID) (sqlc.AuthUser, error)
	GetUserByEmail(ctx context.Context, email string) (sqlc.AuthUser, error)
	UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

// MockQueries implementa um mock para sqlc.Queries
type MockQueries struct {
	CreateUserFn     func(ctx context.Context, arg sqlc.CreateUserParams) error
	GetAllUsersFn    func(ctx context.Context) ([]sqlc.AuthUser, error)
	GetUserByIdFn    func(ctx context.Context, id uuid.UUID) (sqlc.AuthUser, error)
	GetUserByEmailFn func(ctx context.Context, email string) (sqlc.AuthUser, error)
	UpdateUserFn     func(ctx context.Context, arg sqlc.UpdateUserParams) error
	DeleteUserFn     func(ctx context.Context, id uuid.UUID) error
}

func (m *MockQueries) CreateUser(ctx context.Context, arg sqlc.CreateUserParams) error {
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, arg)
	}
	return nil
}

func (m *MockQueries) GetAllUsers(ctx context.Context) ([]sqlc.AuthUser, error) {
	if m.GetAllUsersFn != nil {
		return m.GetAllUsersFn(ctx)
	}
	return []sqlc.AuthUser{}, nil
}

func (m *MockQueries) GetUserById(ctx context.Context, id uuid.UUID) (sqlc.AuthUser, error) {
	if m.GetUserByIdFn != nil {
		return m.GetUserByIdFn(ctx, id)
	}
	return sqlc.AuthUser{}, nil
}

func (m *MockQueries) GetUserByEmail(ctx context.Context, email string) (sqlc.AuthUser, error) {
	if m.GetUserByEmailFn != nil {
		return m.GetUserByEmailFn(ctx, email)
	}
	return sqlc.AuthUser{}, nil
}

func (m *MockQueries) UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) error {
	if m.UpdateUserFn != nil {
		return m.UpdateUserFn(ctx, arg)
	}
	return nil
}

func (m *MockQueries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if m.DeleteUserFn != nil {
		return m.DeleteUserFn(ctx, id)
	}
	return nil
}

// Helper para criar um usuário de testes válido
func createValidTestUser() domain.User {
	return domain.NewUser(
		"123e4567-e89b-12d3-a456-426614174000",
		"João",
		"Silva",
		"joao@example.com",
		"1990-01-15",
		"password1234",
	)
}

// Helper para criar um usuário SQLC de testes válido
func createValidTestSQLCUser() sqlc.AuthUser {
	id := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
	return sqlc.AuthUser{
		ID:        id,
		FirstName: "João",
		LastName:  "Silva",
		Email:     "joao@example.com",
		Password:  "password1234",
		BirthDate: time.Date(1990, 1, 15, 0, 0, 0, 0, time.UTC),
	}
}

func TestNewUserRepository(t *testing.T) {
	queries := &MockQueries{}
	repo := NewUserRepository(queries)

	if repo == nil {
		t.Error("expected repository to be not nil")
	}
}

func TestUserRepository_CreateUser(t *testing.T) {
	tests := []struct {
		name      string
		user      domain.User
		mockError error
		wantErr   bool
	}{
		{
			name:      "should create user successfully",
			user:      createValidTestUser(),
			mockError: nil,
			wantErr:   false,
		},
		{
			name:      "should fail with database error",
			user:      createValidTestUser(),
			mockError: errors.New("database error"),
			wantErr:   true,
		},
		{
			name: "should fail with invalid UUID",
			user: domain.NewUser(
				"invalid-uuid",
				"João",
				"Silva",
				"joao@example.com",
				"1990-01-15",
				"password1234",
			),
			mockError: nil,
			wantErr:   true,
		},
		{
			name: "should fail with invalid date format",
			user: domain.NewUser(
				"123e4567-e89b-12d3-a456-426614174000",
				"João",
				"Silva",
				"joao@example.com",
				"invalid-date",
				"password1234",
			),
			mockError: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQueries := &MockQueries{
				CreateUserFn: func(ctx context.Context, arg sqlc.CreateUserParams) error {
					return tt.mockError
				},
			}

			repo := NewUserRepository(mockQueries)
			err := repo.CreateUser(context.Background(), tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_GetAllUsers(t *testing.T) {
	tests := []struct {
		name              string
		mockUsers         []sqlc.AuthUser
		mockError         error
		expectedUserCount int
		wantErr           bool
	}{
		{
			name: "should get all users successfully",
			mockUsers: []sqlc.AuthUser{
				createValidTestSQLCUser(),
				{
					ID:        uuid.MustParse("223e4567-e89b-12d3-a456-426614174000"),
					FirstName: "Maria",
					LastName:  "Santos",
					Email:     "maria@example.com",
					Password:  "password1234",
					BirthDate: time.Date(1992, 3, 20, 0, 0, 0, 0, time.UTC),
				},
			},
			mockError:         nil,
			expectedUserCount: 2,
			wantErr:           false,
		},
		{
			name:              "should handle empty users list",
			mockUsers:         []sqlc.AuthUser{},
			mockError:         nil,
			expectedUserCount: 0,
			wantErr:           false,
		},
		{
			name:              "should fail with database error",
			mockUsers:         nil,
			mockError:         errors.New("database error"),
			expectedUserCount: 0,
			wantErr:           true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQueries := &MockQueries{
				GetAllUsersFn: func(ctx context.Context) ([]sqlc.AuthUser, error) {
					return tt.mockUsers, tt.mockError
				},
			}

			repo := NewUserRepository(mockQueries)
			users, err := repo.GetAllUsers(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllUsers() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !tt.wantErr && len(users) != tt.expectedUserCount {
				t.Errorf("GetAllUsers() returned %d users, expected %d", len(users), tt.expectedUserCount)
			}
		})
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	validID := "123e4567-e89b-12d3-a456-426614174000"

	tests := []struct {
		name      string
		id        string
		mockUser  sqlc.AuthUser
		mockError error
		wantErr   bool
	}{
		{
			name:      "should get user by ID successfully",
			id:        validID,
			mockUser:  createValidTestSQLCUser(),
			mockError: nil,
			wantErr:   false,
		},
		{
			name:      "should fail with invalid UUID",
			id:        "invalid-uuid",
			mockUser:  sqlc.AuthUser{},
			mockError: nil,
			wantErr:   true,
		},
		{
			name:      "should fail with database error",
			id:        validID,
			mockUser:  sqlc.AuthUser{},
			mockError: errors.New("user not found"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQueries := &MockQueries{
				GetUserByIdFn: func(ctx context.Context, id uuid.UUID) (sqlc.AuthUser, error) {
					return tt.mockUser, tt.mockError
				},
			}

			repo := NewUserRepository(mockQueries)
			user, err := repo.GetUserByID(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByID() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !tt.wantErr && user == nil {
				t.Error("GetUserByID() returned nil user")
			}
		})
	}
}

func TestUserRepository_GetUserByEmail(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		mockUser  sqlc.AuthUser
		mockError error
		wantErr   bool
	}{
		{
			name:      "should get user by email successfully",
			email:     "joao@example.com",
			mockUser:  createValidTestSQLCUser(),
			mockError: nil,
			wantErr:   false,
		},
		{
			name:      "should fail with user not found",
			email:     "notfound@example.com",
			mockUser:  sqlc.AuthUser{},
			mockError: errors.New("no rows in result set"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQueries := &MockQueries{
				GetUserByEmailFn: func(ctx context.Context, email string) (sqlc.AuthUser, error) {
					return tt.mockUser, tt.mockError
				},
			}

			repo := NewUserRepository(mockQueries)
			user, err := repo.GetUserByEmail(context.Background(), tt.email)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByEmail() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !tt.wantErr && user == nil {
				t.Error("GetUserByEmail() returned nil user")
			}
		})
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	tests := []struct {
		name      string
		user      domain.User
		mockError error
		wantErr   bool
	}{
		{
			name:      "should update user successfully",
			user:      createValidTestUser(),
			mockError: nil,
			wantErr:   false,
		},
		{
			name:      "should fail with database error",
			user:      createValidTestUser(),
			mockError: errors.New("update failed"),
			wantErr:   true,
		},
		{
			name: "should fail with invalid UUID",
			user: domain.NewUser(
				"invalid-uuid",
				"João",
				"Silva",
				"joao@example.com",
				"1990-01-15",
				"password1234",
			),
			mockError: nil,
			wantErr:   true,
		},
		{
			name: "should fail with invalid date format",
			user: domain.NewUser(
				"123e4567-e89b-12d3-a456-426614174000",
				"João",
				"Silva",
				"joao@example.com",
				"invalid-date",
				"password1234",
			),
			mockError: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQueries := &MockQueries{
				UpdateUserFn: func(ctx context.Context, arg sqlc.UpdateUserParams) error {
					return tt.mockError
				},
			}

			repo := NewUserRepository(mockQueries)
			err := repo.UpdateUser(context.Background(), tt.user)

			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_DeleteUser(t *testing.T) {
	validID := "123e4567-e89b-12d3-a456-426614174000"

	tests := []struct {
		name      string
		id        string
		mockError error
		wantErr   bool
	}{
		{
			name:      "should delete user successfully",
			id:        validID,
			mockError: nil,
			wantErr:   false,
		},
		{
			name:      "should fail with invalid UUID",
			id:        "invalid-uuid",
			mockError: nil,
			wantErr:   true,
		},
		{
			name:      "should fail with database error",
			id:        validID,
			mockError: errors.New("delete failed"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockQueries := &MockQueries{
				DeleteUserFn: func(ctx context.Context, id uuid.UUID) error {
					return tt.mockError
				},
			}

			repo := NewUserRepository(mockQueries)
			err := repo.DeleteUser(context.Background(), tt.id)

			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteUser() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserRepository_ParamsCorrectlyFormatted(t *testing.T) {
	mockQueries := &MockQueries{
		CreateUserFn: func(ctx context.Context, arg sqlc.CreateUserParams) error {
			// Validar que os parâmetros foram formatados corretamente
			if arg.FirstName != "João" {
				t.Errorf("expected FirstName 'João', got %s", arg.FirstName)
			}
			if arg.Email != "joao@example.com" {
				t.Errorf("expected Email 'joao@example.com', got %s", arg.Email)
			}
			return nil
		},
	}

	repo := NewUserRepository(mockQueries)
	user := createValidTestUser()
	err := repo.CreateUser(context.Background(), user)

	if err != nil {
		t.Fatalf("CreateUser() error = %v", err)
	}
}
