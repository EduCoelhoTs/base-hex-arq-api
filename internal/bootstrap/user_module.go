package bootstrap

import (
	"context"

	"github.com/EduCoelhoTs/base-hex-arq-api/internal/adapter/repository/postgres"
	"github.com/EduCoelhoTs/base-hex-arq-api/internal/application/service"
	userusecase "github.com/EduCoelhoTs/base-hex-arq-api/internal/application/usecase/user"
	usercontroller "github.com/EduCoelhoTs/base-hex-arq-api/internal/infra/controller/user"
)

func NewUserModule(ctx context.Context, db postgres.QueriesPort) usercontroller.Controller {
	repository := postgres.NewUserRepository(db)
	service := service.NewUserService(ctx, repository)
	usecase := userusecase.NewCreateUserUseCase(service)

	return usercontroller.NewController(usecase)
}
