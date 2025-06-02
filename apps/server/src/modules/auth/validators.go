package auth

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidators registers custom validation functions
func RegisterCustomValidators(v *validator.Validate) {
	v.RegisterValidation("password", validatePassword)
}

// validatePassword checks if the password meets all requirements:
// - minimum length of 8 characters
// - at least one uppercase letter
// - at least one lowercase letter
// - at least one number
// - at least one special character
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check minimum length
	if len(password) < 8 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case strings.ContainsRune(specialChars, char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}
