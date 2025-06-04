package dto

import (
	"net/mail"
	"strings"
	"unicode/utf8"
)

// CreateUserRequest defines the expected structure for a new user registration request.
// It includes fields for username, email, and password.
// It implements the Validator interface (defined in the api package) for self-validation.
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Valid checks the validity of the CreateUserRequest fields.
// It returns a map of validation errors if any are found, otherwise nil.
func (r *CreateUserRequest) Valid() map[string]string {
	errors := make(map[string]string)

	// Username validation
	if strings.TrimSpace(r.Username) == "" {
		errors["username"] = "username must be provided"
	} else if utf8.RuneCountInString(r.Username) < 3 {
		errors["username"] = "username must be at least 3 characters long"
	} else if utf8.RuneCountInString(r.Username) > 50 {
		errors["username"] = "username must not be more than 50 characters long"
	}

	// Email validation
	if strings.TrimSpace(r.Email) == "" {
		errors["email"] = "email must be provided"
	} else {
		_, err := mail.ParseAddress(r.Email)
		if err != nil {
			errors["email"] = "email must be a valid email address"
		}
	}

	// Password validation
	if strings.TrimSpace(r.Password) == "" {
		errors["password"] = "password must be provided"
	} else if utf8.RuneCountInString(r.Password) < 8 {
		errors["password"] = "password must be at least 8 characters long"
	} else if utf8.RuneCountInString(r.Password) > 72 { // bcrypt has a 72-byte limit
		errors["password"] = "password must not be more than 72 characters long"
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}

// TODO: Add DTOs for VendorRequest, MerchantRequest, LoginRequest etc. as needed.
