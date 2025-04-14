package usecases

type RegisterInputDTO struct {
	Email    string
	Password string
	Role     string
}

type RegisterOutputDTO struct {
	ID    string
	Email string
	Role  string
}

type LoginInputDTO struct {
	Email    string
	Password string
}

type LoginOutputDTO struct {
	Token string
}
