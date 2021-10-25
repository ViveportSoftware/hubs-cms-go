package config

import (
	"fmt"
)

func GetMastodonVerifyCredentialsURI() string {
	return fmt.Sprintf("%s/api/v1/accounts/verify_credentials", EnvVariable.MastodonBaseURI)
}

func GetMastodonUpdateCredentialsURI() string {
	return fmt.Sprintf("%s/api/v1/accounts/update_credentials", EnvVariable.MastodonBaseURI)
}
