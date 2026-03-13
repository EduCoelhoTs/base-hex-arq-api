package port_test

import (
	"context"
	"testing"

	"github.com/EduCoelhoTs/base-hex-arq-api/internal/core/domain"
	mock_port "github.com/EduCoelhoTs/base-hex-arq-api/internal/core/port/mock"
	"github.com/EduCoelhoTs/base-hex-arq-api/pkg/xuuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserRepository(t *testing.T) {
	testCtlr := gomock.NewController(t)
	defer testCtlr.Finish()
	ctx := context.Background()

	user := domain.NewUser(
		xuuid.NewV7(),
		"John",
		"Silva",
		"johnSilva@gmail.com",
		"2000-10-22",
		"12345678",
	)

	mockUserRepository := mock_port.NewMockUserRepositoryInterface(testCtlr)
	mockUserRepository.EXPECT().CreateUser(
		ctx,
		user,
	).Return(nil).Times(1)

	mockUserRepository.EXPECT().GetUserByID(ctx, user.GetID()).Return(user, nil).Times(1)
	mockUserRepository.EXPECT().GetAllUsers(ctx).Return([]domain.User{user}, nil).Times(1)
	mockUserRepository.EXPECT().GetUserByEmail(ctx, user.GetEmail()).Return(user, nil).Times(1)
	mockUserRepository.EXPECT().UpdateUser(
		ctx,
		user,
	).Return(nil).Times(1)
	mockUserRepository.EXPECT().DeleteUser(ctx, user.GetID()).Return(nil).Times(1)

	testCases := []struct {
		name string
		fn   func(t *testing.T)
	}{
		{
			name: "should create user",
			fn: func(t *testing.T) {
				err := mockUserRepository.CreateUser(ctx, user)

				require.Nil(t, err)
			},
		},
		{
			name: "should get user by id",
			fn: func(t *testing.T) {
				foundUser, err := mockUserRepository.GetUserByID(ctx, user.GetID())

				require.Nil(t, err)
				require.Equal(t, user.GetID(), foundUser.GetID())
			},
		},
		{
			name: "should get all users",
			fn: func(t *testing.T) {
				users, err := mockUserRepository.GetAllUsers(ctx)

				require.Nil(t, err)
				require.Len(t, users, 1)
				require.Equal(t, user.GetID(), users[0].GetID())
			},
		},
		{
			name: "should get user by email",
			fn: func(t *testing.T) {
				foundUser, err := mockUserRepository.GetUserByEmail(ctx, user.GetEmail())

				require.Nil(t, err)
				require.Equal(t, user.GetID(), foundUser.GetID())
			},
		},
		{
			name: "should update user",
			fn: func(t *testing.T) {
				err := mockUserRepository.UpdateUser(
					ctx,
					user,
				)

				require.Nil(t, err)
			},
		},
		{
			name: "should delete user",
			fn: func(t *testing.T) {
				err := mockUserRepository.DeleteUser(ctx, user.GetID())

				require.Nil(t, err)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.fn)
	}
}
