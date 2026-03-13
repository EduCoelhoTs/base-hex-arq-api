package port

import (
	"context"

	"github.com/EduCoelhoTs/base-hex-arq-api/internal/core/domain"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetAllUsers(ctx context.Context) ([]domain.User, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) error
	DeleteUser(ctx context.Context, id string) error
}
