package entity

import "github.com/go-playground/validator"

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type SignInInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,gte=6"`
}

type SignUpInput struct {
	Name     string `validate:"required,gte=2"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,gte=6"`
}

func (i SignInInput) Validate() error {
	return validate.Struct(i)
}

func (i SignUpInput) Validate() error {
	return validate.Struct(i)
}
