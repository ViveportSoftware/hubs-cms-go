package validators

import (
	"strconv"

	"github.com/go-playground/validator/v10"
)

// PageLimitValidator checks if page limit is valid
func PageLimitValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return false
	}
	if intValue < 0 || intValue > 50 {
		return false
	}
	return true
}

// PageStartValidator checks if page start is valid
func PageStartValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return false
	}
	if intValue < 0 {
		return false
	}
	return true
}
