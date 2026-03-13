package userusecase

import (
	"context"

	"github.com/EduCoelhoTs/base-hex-arq-api/internal/core/domain"
	"github.com/EduCoelhoTs/base-hex-arq-api/internal/core/port"
	"github.com/EduCoelhoTs/base-hex-arq-api/pkg/xuuid"
)

type createUserUseCase struct {
	repository port.UserRepositoryInterface
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
	Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error)
}

func NewCreateUserUseCase(repository port.UserRepositoryInterface) CreateUserUseCase {
	return &createUserUseCase{repository: repository}
}

func (uc *createUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	id := xuuid.NewV7()
	user := domain.NewUser(
		id,
		input.FirstName,
		input.LastName,
		input.Email,
		input.BirthDate,
		input.Password,
	)

	if err := user.IsValid(); err != nil {
		return nil, err
	}

	if err := uc.repository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &CreateUserOutput{ID: user.GetID()}, nil
}
