package config

import (
	"fmt"
)

func GetDirectusAccessTokenURI() string {
	return fmt.Sprintf("%s/auth/login", EnvVariable.DirectusBaseURI)
}
