package models

import "github.com/go-playground/validator"

var (
	validation = validator.New()
)

type Register struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=32"`
}

func (r *Register) Validate() error {
	return validation.Struct(r)
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=32"`
}

func (r *Login) Validate() error {
	return validation.Struct(r)
}

type SetPassword struct {
	OldPassword string `json:"old" validate:"required,min=3,max=32"`
	NewPassword string `json:"new" validate:"required,min=3,max=32"`
}

func (s *SetPassword) Validate() error {
	return validation.Struct(s)
}
