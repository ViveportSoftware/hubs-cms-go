package validators

import (
	"hubs-cms-go/logger"
	"strings"

	"github.com/go-playground/validator/v10"
)

// IsInvalid finds specific invalid field
func IsInvalid(namespace string, err error) bool {
	ve, ok := err.(validator.ValidationErrors)
	if ok {
		for _, e := range ve {
			logger.Info.Printf("ns: %v\n", e.Namespace())
			if strings.HasSuffix(e.Namespace(), namespace) {
				return true
			}
		}
	}
	return false
}

// StringNotEmptyValidator checks if value is not an empty string
func StringNotEmptyValidator(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) > 0
}
