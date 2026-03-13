package userusecase

import (
	"context"
	"log/slog"

	"github.com/EduCoelhoTs/base-hex-arq-api/internal/core/domain"
	port "github.com/EduCoelhoTs/base-hex-arq-api/internal/core/port/user"
	"github.com/EduCoelhoTs/base-hex-arq-api/pkg/xcrypto"
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
	hashedPassword, err := xcrypto.HashPassword(input.Password)
	if err != nil {
		slog.Error("usecase.create.execute", "hashpassword", err.Error())
		return nil, err
	}

	user := domain.NewUser(
		id,
		input.FirstName,
		input.LastName,
		input.Email,
		input.BirthDate,
		hashedPassword,
	)

	if err := user.IsValid(); err != nil {
		return nil, err
	}

	if err := uc.repository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return &CreateUserOutput{ID: user.GetID()}, nil
}
