package dto

type MastodonVerifyCredentialsResponse struct {
	ID              string `json:"id"`
	UserName        string `json:"username"`
	MastodonAccount string `json:"acct"`
	DisplayName     string `json:"display_name"`
	MastodonAvatar  string `json:"avatar_static"`
	MastodonToken   string `json:"token"`
}

// Validate validates MastodonVerifyCredentialsResponse field
func (r MastodonVerifyCredentialsResponse) Validate() bool {
	return len(r.UserName) > 0 && len(r.MastodonAccount) > 0
}

type MastodonPatchAccountRequestBody struct {
	DisplayName string `json:"display_name"`
}
