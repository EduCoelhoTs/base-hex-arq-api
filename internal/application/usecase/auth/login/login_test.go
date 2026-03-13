package loginusecase_test

import (
	"context"
	"testing"

	loginusecase "github.com/EduCoelhoTs/base-hex-arq-api/internal/application/usecase/auth/login"
	"github.com/EduCoelhoTs/base-hex-arq-api/internal/core/domain"
	"github.com/EduCoelhoTs/base-hex-arq-api/pkg/xcrypto"
	"github.com/stretchr/testify/require"
)

type MockUserRepository struct{}

func createMockPassword() string {
	hashPass, _ := xcrypto.HashPassword("password")
	return hashPass
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {

	// Mock implementation to return a user with the given email
	return domain.NewUser(
		"123",
		"John",
		"Doe",
		email,
		"1990-01-01",
		createMockPassword(),
	), nil
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user domain.User) error {
	return nil
}

func (m *MockUserRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	return nil, nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	return nil, nil
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user domain.User) error {
	return nil
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id string) error {
	return nil
}

type MockTokenService struct{}

func (m *MockTokenService) Generate(userID string) (string, error) {
	// Mock implementation to return a token for the given user ID
	return "mocked.jwt.token", nil
}

func (m *MockTokenService) Validate(token string) (string, error) {
	return "123", nil
}

func TestLoginUseCase(t *testing.T) {

	repository := &MockUserRepository{}
	tokenService := &MockTokenService{}

	testCases := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			"shoud login user",

			func(t *testing.T) {
				useCase := loginusecase.NewLoginUseCase(tokenService, repository)
				token, err := useCase.Execute(context.Background(), "teste@email.com", "password")

				require.Nil(t, err)
				require.NotEmpty(t, token)

			},
		},
		{
			"shoud return error if password is wrong",
			func(t *testing.T) {
				useCase := loginusecase.NewLoginUseCase(tokenService, repository)
				token, err := useCase.Execute(context.Background(), "teste@email.com", "wrongpassword")

				require.Error(t, err)
				require.Empty(t, token)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.fn)
	}
}
