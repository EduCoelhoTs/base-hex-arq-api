package bootstrap

import (
	"context"

	"github.com/EduCoelhoTs/base-hex-arq-api/internal/adapter/repository/postgres"
	usercontroller "github.com/EduCoelhoTs/base-hex-arq-api/internal/infra/controller/user"
)

type App struct {
	UserController usercontroller.Controller
}

func NewApp(ctx context.Context, db postgres.QueriesPort) *App {
	return &App{
		UserController: NewUserModule(ctx, db),
	}
}
