package dto

type DirectusAuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DirectusAuthLoginResponse struct {
	Data DirectusAuthLoginResponseData `json:"data"`
}

type DirectusAuthLoginResponseData struct {
	AccessToken  string `json:"access_token"`
	Expires      int64  `json:"expires"`
	RefreshToken string `json:"refresh_token"`
}

// Validate validates DirectusAuthLoginResponse field
func (r DirectusAuthLoginResponse) Validate() bool {
	return r.Data != DirectusAuthLoginResponseData{} &&
		len(r.Data.AccessToken) > 0 &&
		len(r.Data.RefreshToken) > 0 &&
		r.Data.Expires > 0
}
