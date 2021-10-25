package service

import (
	"fmt"
	"hubs-cms-go/client"
	"hubs-cms-go/config"
	"hubs-cms-go/constant"
	"hubs-cms-go/dto"
	"hubs-cms-go/logger"
)

func GetMastodonVerifyCredentials(token string) (dto.MastodonVerifyCredentialsResponse, error) {

	verifyCredentialsResponse := dto.MastodonVerifyCredentialsResponse{}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, token).
		SetResult(&verifyCredentialsResponse).
		Get(config.GetMastodonVerifyCredentialsURI())
	if err != nil {
		return dto.MastodonVerifyCredentialsResponse{}, err
	}
	if !response.IsSuccess() {
		return dto.MastodonVerifyCredentialsResponse{}, fmt.Errorf("[GetMastodonVerifyCredentials] server response error code: %v", response.StatusCode())
	}

	if !verifyCredentialsResponse.Validate() {
		logger.Error.Printf("[GetMastodonVerifyCredentials] %v\n", err)
		return dto.MastodonVerifyCredentialsResponse{}, fmt.Errorf("[GetMastodonVerifyCredentials] verifyCredentialsResponse.Validate failed")
	}
	return verifyCredentialsResponse, nil
}

func PatchMastodonAccount(mastodonToken string, patchRequest dto.MastodonPatchAccountRequestBody) (dto.MastodonVerifyCredentialsResponse, error) {

	if len(mastodonToken) == 0 {
		logger.Error.Printf("[PatchMastodonAccount] unable to get mastodon token error\n")
		return dto.MastodonVerifyCredentialsResponse{}, fmt.Errorf("[PatchMastodonAccount] unable to get mastodon token error")
	}

	verifyCredentialsResponse := dto.MastodonVerifyCredentialsResponse{}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, mastodonToken).
		SetResult(&verifyCredentialsResponse).
		SetFormData(map[string]string{
			"display_name": patchRequest.DisplayName,
		}).
		Patch(config.GetMastodonUpdateCredentialsURI())

	if err != nil {
		logger.Error.Printf("[PatchMastodonAccount] request error: %v\n", err)
		return dto.MastodonVerifyCredentialsResponse{}, err
	}

	if !response.IsSuccess() {
		logger.Error.Printf("[PatchMastodonAccount] server response error status: %v\n", response.StatusCode())
		return dto.MastodonVerifyCredentialsResponse{}, fmt.Errorf("[PatchMastodonAccount] server response error status: %v", response.StatusCode())
	}

	if !verifyCredentialsResponse.Validate() {
		logger.Error.Printf("[PatchMastodonAccount] %v\n", err)
		return dto.MastodonVerifyCredentialsResponse{}, fmt.Errorf("[PatchMastodonAccount] verifyCredentialsResponse.Validate failed")
	}
	return verifyCredentialsResponse, nil
}
