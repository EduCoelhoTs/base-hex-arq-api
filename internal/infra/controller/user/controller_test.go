package usercontroller_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	userusecase "github.com/EduCoelhoTs/base-hex-arq-api/internal/application/usecase/user"
	usercontroller "github.com/EduCoelhoTs/base-hex-arq-api/internal/infra/controller/user"
	"github.com/stretchr/testify/require"
)

type MockCreateUserUseCase struct{}

func (m *MockCreateUserUseCase) Execute(ctx context.Context, input userusecase.CreateUserInput) (*userusecase.CreateUserOutput, error) {
	return &userusecase.CreateUserOutput{ID: "12345"}, nil
}

func TestController(t *testing.T) {

	controller := usercontroller.NewController(&MockCreateUserUseCase{})

	testCases := []struct {
		name       string
		body       usercontroller.CreateUserRequestDTO
		statusCode int
	}{
		{
			name: "shoud create user",
			body: usercontroller.CreateUserRequestDTO{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "johndoe@email.com",
				BirthDate: "1990-01-01",
				Password:  "password123",
			},
			statusCode: http.StatusCreated,
		},
		{
			name: "shoud return validation error and status 400",
			body: usercontroller.CreateUserRequestDTO{
				LastName:  "Doe",
				Email:     "johndoe@email.com",
				BirthDate: "1990-01-01",
				Password:  "password123",
			},
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			marshalledBody, _ := json.Marshal(tc.body)
			parsedBody := bytes.NewReader(marshalledBody)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/user", parsedBody)
			request.Header.Set("Content-Type", "application/json")

			handler := http.HandlerFunc(controller.Create)
			handler.ServeHTTP(recorder, request)

			require.Equal(t, recorder.Result().StatusCode, tc.statusCode)
		})
	}
}
