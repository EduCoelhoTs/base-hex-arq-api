package userusecase_test

import (
	"context"
	"testing"

	userusecase "github.com/EduCoelhoTs/nba-predict-api/internal/application/usecase/user"
	mock_userusecase "github.com/EduCoelhoTs/nba-predict-api/internal/application/usecase/user/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreate(t *testing.T) {
	testController := gomock.NewController(t)
	defer testController.Finish()

	ctx := context.Background()
	userUseCase := mock_userusecase.NewMockCreateUserUseCase(testController)

	testsCases := []struct {
		name   string
		input  userusecase.CreateUserInput
		output userusecase.CreateUserOutput
		fn     func(t *testing.T)
	}{
		{
			name: "should create a user successfully",
			fn: func(t *testing.T) {
				input := userusecase.CreateUserInput{
					FirstName: "João",
					LastName:  "Silva",
					Email:     "joao.silva@example.com",
					BirthDate: "1990-01-15",
					Password:  "password1234",
				}
				output := userusecase.CreateUserOutput{
					ID: "12345",
				}
				userUseCase.EXPECT().Execute(ctx, input).Return(output, nil)

				output, err := userUseCase.Execute(ctx, userusecase.CreateUserInput{
					FirstName: "João",
					LastName:  "Silva",
					Email:     "joao.silva@example.com",
					BirthDate: "1990-01-15",
					Password:  "password1234",
				})
				require.Nil(t, err)
				assert.Equal(t, output.ID, "12345")
			},
		},
		{
			name: "should fail when email already exists",
			fn: func(t *testing.T) {
				input := userusecase.CreateUserInput{
					FirstName: "João",
					LastName:  "Silva",
					Email:     "joao.silva@example.com",
					BirthDate: "1990-01-15",
					Password:  "password1234",
				}
				userUseCase.EXPECT().Execute(ctx, input).Return(userusecase.CreateUserOutput{}, assert.AnError)

				_, err := userUseCase.Execute(ctx, input)
				require.NotNil(t, err)
			},
		},
		{
			name: "should fail when email is invalid",
			fn: func(t *testing.T) {
				input := userusecase.CreateUserInput{
					FirstName: "João",
					LastName:  "Silva",
					Email:     "invalid-email",
					BirthDate: "1990-01-15",
					Password:  "password1234",
				}
				userUseCase.EXPECT().Execute(ctx, input).Return(userusecase.CreateUserOutput{}, assert.AnError)

				_, err := userUseCase.Execute(ctx, input)
				require.NotNil(t, err)
			},
		},
		{
			name: "should fail when password is too short",
			fn: func(t *testing.T) {
				input := userusecase.CreateUserInput{
					FirstName: "João",
					LastName:  "Silva",
					Email:     "joao.silva@example.com",
					BirthDate: "1990-01-15",
					Password:  "123",
				}
				userUseCase.EXPECT().Execute(ctx, input).Return(userusecase.CreateUserOutput{}, assert.AnError)

				_, err := userUseCase.Execute(ctx, input)
				require.NotNil(t, err)
			},
		},
		{
			name: "should fail when required fields are missing",
			fn: func(t *testing.T) {
				input := userusecase.CreateUserInput{
					FirstName: "",
					LastName:  "Silva",
					Email:     "joao.silva@example.com",
					BirthDate: "1990-01-15",
					Password:  "password1234",
				}
				userUseCase.EXPECT().Execute(ctx, input).Return(userusecase.CreateUserOutput{}, assert.AnError)

				_, err := userUseCase.Execute(ctx, input)
				require.NotNil(t, err)
			},
		},
	}

	for _, tc := range testsCases {
		t.Run(tc.name, tc.fn)
	}
}
