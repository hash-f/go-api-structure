package dto

import (
	"github.com/go-playground/validator/v10"
)

// LoginUserRequest defines the structure for a user login request.
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,trimLenMin=8,trimLenMax=72,min=8,max=72"`
}

// Valid checks if the LoginUserRequest fields are valid.
func (r *LoginUserRequest) Valid() map[string]string {
	err := Validator().Struct(r)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "Email":
			if err.Tag() == "required" {
				errors["email"] = "email must be provided"
			} else {
				errors["email"] = "email must be a valid email address"
			}
		case "Password":
			switch err.Tag() {
			case "required":
				errors["password"] = "password must be provided"
			case "trimLenMin":
				errors["password"] = "password must be at least 8 characters long"
			case "trimLenMax":
				errors["password"] = "password must not be more than 72 characters long"
			}
		}
	}

	return errors
}
