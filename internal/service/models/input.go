package models

type SignUpInput struct {
	Name              string `validate:"required,min=5,max=100" json:"name"`
	Email             string `validate:"required,max=100,email" json:"email"`
	Password          string `validate:"required,min=8,max=72" json:"password"`
	PasswordConfirmed string `validate:"required,min=8,max=72" json:"password_confirmed"`
}

type SignInInput struct {
	Email    string `validate:"required,max=100,email" json:"email"`
	Password string `validate:"required,min=8,max=72" json:"password"`
}
