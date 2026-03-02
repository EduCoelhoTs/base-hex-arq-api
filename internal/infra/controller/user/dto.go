package usercontroller

type CreateUserRequestDTO struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	BirthDate string `json:"birth_date" validate:"required,datetime=2006-01-02"`
	Password  string `json:"password" validate:"required,min=8"`
}

type CreateUserResponseDTO struct {
	ID string `json:"id"`
}
