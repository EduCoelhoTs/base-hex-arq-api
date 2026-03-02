package usercontroller

import (
	"fmt"
	"net/http"

	userusecase "github.com/EduCoelhoTs/nba-predict-api/internal/application/usecase/user"
	infraerrors "github.com/EduCoelhoTs/nba-predict-api/internal/infra/error"
	"github.com/EduCoelhoTs/nba-predict-api/pkg/xjson"
	"github.com/EduCoelhoTs/nba-predict-api/pkg/xvalidator"
)

type UseCase struct {
	CreateUserUseCase userusecase.CreateUserUseCase
}

type controller struct {
	useCase UseCase
}

func NewController(useCase UseCase) *controller {
	return &controller{useCase: useCase}
}

func (c *controller) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var body CreateUserRequestDTO
	if err := xjson.Decode(r.Body, &body); err != nil {
		fmt.Printf("usercontroller.create decode body error: %w", err)
		xjson.ResponseHttpError(w, http.StatusBadRequest, infraerrors.ERR_INTERNAL_SERVER_ERROR)
	}

	if err := xvalidator.ValidateStruct(body); err != nil {
		fmt.Printf("usercontroller.create validate body error: %w", err)
		xjson.ResponseHttpError(w, http.StatusBadRequest, err.Error())
		return
	}

	input := userusecase.CreateUserInput{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		BirthDate: body.BirthDate,
		Password:  body.Password,
	}

	output, err := c.useCase.CreateUserUseCase.Execute(r.Context(), input)
	if err != nil {
		fmt.Printf("usercontroller.create execute usecase error: %w", err)
		xjson.ResponseHttpError(w, http.StatusInternalServerError, infraerrors.ERR_INTERNAL_SERVER_ERROR)
		return
	}

	xjson.ReponseHttp(w, http.StatusCreated, CreateUserResponseDTO{ID: output.ID})
}
