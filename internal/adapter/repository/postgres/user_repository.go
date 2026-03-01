package postgres

import (
	"context"

	"github.com/EduCoelhoTs/nba-predict-api/internal/adapter/repository/postgres/sqlc"
	"github.com/EduCoelhoTs/nba-predict-api/internal/core/domain"
	"github.com/EduCoelhoTs/nba-predict-api/internal/core/port"
	"github.com/EduCoelhoTs/nba-predict-api/pkg/xdate"
	"github.com/EduCoelhoTs/nba-predict-api/pkg/xuuid"
	"github.com/google/uuid"
)

// QueriesPort define a porta para acesso ao banco de dados (Hexagonal Architecture)
type QueriesPort interface {
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) error
	GetAllUsers(ctx context.Context) ([]sqlc.AuthUser, error)
	GetUserById(ctx context.Context, id uuid.UUID) (sqlc.AuthUser, error)
	GetUserByEmail(ctx context.Context, email string) (sqlc.AuthUser, error)
	UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userRepository struct {
	queries QueriesPort
}

func NewUserRepository(queries QueriesPort) port.UserRepositoryInterface {
	return &userRepository{queries: queries}
}

func (r *userRepository) CreateUser(ctx context.Context, user domain.User) error {
	parsedId, err := xuuid.UUIDFromString(user.GetID())
	if err != nil {
		return err
	}

	parsedDate, err := xdate.ParseDate(user.GetBirthDate(), nil, nil)
	if err != nil {
		return err
	}

	userParams := sqlc.CreateUserParams{
		ID:        parsedId,
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Email:     user.GetEmail(),
		Password:  user.GetPassword(),
		BirthDate: parsedDate,
	}
	return r.queries.CreateUser(ctx, userParams)
}
func (r *userRepository) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	dbUsers, err := r.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var users []domain.User
	for _, dbUser := range dbUsers {
		users = append(users, dbUser.ToDomain())
	}

	return users, nil
}
func (r *userRepository) GetUserByID(ctx context.Context, id string) (domain.User, error) {
	parsedId, err := xuuid.UUIDFromString(id)
	if err != nil {
		return nil, err
	}

	dbUser, err := r.queries.GetUserById(ctx, parsedId)
	if err != nil {
		return nil, err
	}

	return dbUser.ToDomain(), nil
}
func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return dbUser.ToDomain(), nil
}
func (r *userRepository) UpdateUser(ctx context.Context, user domain.User) error {
	parsedId, err := xuuid.UUIDFromString(user.GetID())
	if err != nil {
		return err
	}

	parsedDate, err := xdate.ParseDate(user.GetBirthDate(), nil, nil)
	if err != nil {
		return err
	}

	userParams := sqlc.UpdateUserParams{
		ID:        parsedId,
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Email:     user.GetEmail(),
		Password:  user.GetPassword(),
		BirthDate: parsedDate,
	}
	return r.queries.UpdateUser(ctx, userParams)
}
func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	parsedId, err := xuuid.UUIDFromString(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteUser(ctx, parsedId)
}
