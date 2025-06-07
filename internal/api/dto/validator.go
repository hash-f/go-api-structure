package dto

import (
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	// validatorInstance is the singleton validator instance
	validatorInstance *validator.Validate
	// once ensures the validator is initialized only once
	once sync.Once
)

// Validator returns a singleton validator instance with custom validations
func Validator() *validator.Validate {
	once.Do(func() {
		validatorInstance = validator.New()
		
		// Register custom validation for minimum length after trimming
		_ = validatorInstance.RegisterValidation("trimLenMin", func(fl validator.FieldLevel) bool {
			param, err := strconv.Atoi(fl.Param())
			if err != nil {
				return false
			}
			trimmed := strings.TrimSpace(fl.Field().String())
			return len(trimmed) >= param
		})
		
		// Register custom validation for maximum length after trimming
		_ = validatorInstance.RegisterValidation("trimLenMax", func(fl validator.FieldLevel) bool {
			param, err := strconv.Atoi(fl.Param())
			if err != nil {
				return false
			}
			trimmed := strings.TrimSpace(fl.Field().String())
			return len(trimmed) <= param
		})
	})
	
	return validatorInstance
}
