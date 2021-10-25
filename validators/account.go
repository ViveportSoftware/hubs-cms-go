package validators

import "github.com/go-playground/validator/v10"

// AccountDisplayNameValidator checks if value is allowed
func AccountDisplayNameValidator(fl validator.FieldLevel) bool {
	// TODO: display name black list?
	return true
}
