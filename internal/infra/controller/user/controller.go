package usercontroller

import (
	"fmt"
	"net/http"

	_http "github.com/EduCoelhoTs/base-hex-arq-api/internal/adapter/http"
	userusecase "github.com/EduCoelhoTs/base-hex-arq-api/internal/application/usecase/user"
	infraerrors "github.com/EduCoelhoTs/base-hex-arq-api/internal/infra/error"
	"github.com/EduCoelhoTs/base-hex-arq-api/pkg/xjson"
	"github.com/EduCoelhoTs/base-hex-arq-api/pkg/xvalidator"
)

type Controller interface {
	GetRoutes() _http.Routes
	Create(w http.ResponseWriter, r *http.Request)
}

type controller struct {
	CreateUserUseCase userusecase.CreateUserUseCase
}

func NewController(createUser userusecase.CreateUserUseCase) *controller {
	return &controller{CreateUserUseCase: createUser}
}

func (c *controller) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var body CreateUserRequestDTO
	if err := xjson.Decode(r.Body, &body); err != nil {
		fmt.Printf("usercontroller.create decode body error: %v", err)
		xjson.ResponseHttpError(w, http.StatusBadRequest, infraerrors.ERR_INTERNAL_SERVER_ERROR)
	}

	if err := xvalidator.ValidateStruct(body); err != nil {
		fmt.Printf("usercontroller.create validate body error: %v", err)
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

	output, err := c.CreateUserUseCase.Execute(r.Context(), input)
	if err != nil {
		fmt.Printf("usercontroller.create execute usecase error: %v", err)
		xjson.ResponseHttpError(w, http.StatusInternalServerError, infraerrors.ERR_INTERNAL_SERVER_ERROR)
		return
	}

	xjson.ReponseHttp(w, http.StatusCreated, CreateUserResponseDTO{ID: output.ID})
}
