package userusecase

import (
	"context"

	"github.com/EduCoelhoTs/nba-predict-api/internal/core/port"
)

type createUserUseCase struct {
	service port.UserServiceInterface
}

type CreateUserInput struct {
	FirstName string
	LastName  string
	Email     string
	BirthDate string
	Password  string
}

type CreateUserOutput struct {
	ID string
}

type CreateUserUseCase interface {
	Execute(ctx context.Context, input CreateUserInput) (CreateUserOutput, error)
}

func NewCreateUserUseCase(service port.UserServiceInterface) CreateUserUseCase {
	return &createUserUseCase{service: service}
}

func (uc *createUserUseCase) Execute(ctx context.Context, input CreateUserInput) (CreateUserOutput, error) {
	user, err := uc.service.CreateUser(
		ctx,
		input.FirstName,
		input.LastName,
		input.Email,
		input.BirthDate,
		input.Password,
	)

	if err != nil {
		return CreateUserOutput{}, err
	}

	return CreateUserOutput{ID: user.GetID()}, nil
}
