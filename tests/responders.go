package tests

import (
	"hubs-cms-go/config"
	"hubs-cms-go/dto"

	"github.com/jarcoal/httpmock"
)

//
//Purpose: always response directus token for event test
//
func regDirTokenRes() {
	getDirectusAccessTokenJsonResponder, _ := httpmock.NewJsonResponder(200,
		dto.DirectusAuthLoginResponse{
			Data: dto.DirectusAuthLoginResponseData{
				AccessToken:  "xyz",
				Expires:      900000,
				RefreshToken: "abc"}})

	httpmock.RegisterResponder(
		"POST",
		config.GetDirectusAccessTokenURI(),
		getDirectusAccessTokenJsonResponder)
}

func setUpResponder(statusCode int, body interface{}, httpMethod string, url string) {
	jsonResponder, _ := httpmock.NewJsonResponder(statusCode,
		body)

	httpmock.RegisterResponder(
		httpMethod,
		url,
		jsonResponder)
}

func setUpErrorResponder(body error, httpMethod string, url string) {
	jsonResponder := httpmock.NewErrorResponder(body)

	httpmock.RegisterResponder(
		httpMethod,
		url,
		jsonResponder)
}
