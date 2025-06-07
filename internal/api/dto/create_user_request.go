package dto

import (
	"github.com/go-playground/validator/v10"
)

// CreateUserRequest defines the expected structure for a new user registration request.
// It includes fields for username, email, and password.
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,trimLenMin=3,trimLenMax=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,trimLenMin=8,trimLenMax=72"`
}



// Valid checks the validity of the CreateUserRequest fields.
// It returns a map of validation errors if any are found, otherwise nil.
func (r *CreateUserRequest) Valid() map[string]string {
	err := Validator().Struct(r)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		switch err.Field() {
		case "Username":
			switch err.Tag() {
			case "required":
				errors["username"] = "username must be provided"
			case "trimLenMin":
				errors["username"] = "username must be at least 3 characters long"
			case "trimLenMax":
				errors["username"] = "username must not be more than 50 characters long"
			}
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
