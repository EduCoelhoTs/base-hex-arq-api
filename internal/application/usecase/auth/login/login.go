package loginusecase

import (
	"context"

	authport "github.com/EduCoelhoTs/base-hex-arq-api/internal/core/port/auth"
	userport "github.com/EduCoelhoTs/base-hex-arq-api/internal/core/port/user"
	"github.com/EduCoelhoTs/base-hex-arq-api/pkg/xcrypto"
)

type LoginUseCase struct {
	tokenService   authport.TokenService
	userRepository userport.UserRepositoryInterface
}

func NewLoginUseCase(tokenService authport.TokenService, userRepository userport.UserRepositoryInterface) *LoginUseCase {
	return &LoginUseCase{
		tokenService:   tokenService,
		userRepository: userRepository,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, email, password string) (string, error) {
	user, err := uc.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if err := xcrypto.ComparePassword(user.GetPassword(), password); err != nil {
		return "", err
	}

	token, err := uc.tokenService.Generate(user.GetID())
	if err != nil {
		return "", err
	}

	return token, nil
}
