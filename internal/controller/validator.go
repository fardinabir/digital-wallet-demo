package controller

import (
	"github.com/fardinabir/digital-wallet-demo/internal/model"
	"github.com/go-playground/validator/v10"
)

// CustomValidator is a custom validator for the echo framework
type CustomValidator struct {
	validator *validator.Validate
}

// Validate validates the input struct
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

// NewCustomValidator return a custom validator struct registering custom validator functions
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// Register the custom validation for Priority
	_ = v.RegisterValidation("validPriority", model.IsValidPriority)
	_ = v.RegisterValidation("validStatus", model.IsValidStatus)

	return &CustomValidator{validator: v}
}
