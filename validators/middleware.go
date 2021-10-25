package validators

import (
	"hubs-cms-go/utils"

	"github.com/go-playground/validator/v10"
)

// BearerTokenValidator checks if the Authorization header starts with "Bearer "
func BearerTokenValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	_, ok := utils.GetBearerToken(value)
	return ok
}
