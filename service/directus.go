package service

import (
	"fmt"
	"hubs-cms-go/cache"
	"hubs-cms-go/client"
	"hubs-cms-go/config"
	"hubs-cms-go/constant"
	"hubs-cms-go/dto"
	"hubs-cms-go/logger"
	"time"
)

func getDirectusAccessTokenFromCache() (string, error) {
	accessToken, found := cache.Store.Get(constant.CacheKeyDirectusAccessToken)
	if found {
		return accessToken.(string), nil
	}
	return "", fmt.Errorf("no cached access token found")
}

func setDirectusAccessTokenToCache(accessToken string, duration time.Duration) {
	cache.Store.Set(constant.CacheKeyDirectusAccessToken, accessToken, duration)
}

func GetDirectusAccessToken(forceFetchFromServer bool) (string, error) {

	if !forceFetchFromServer {
		cachedAccessToken, err := getDirectusAccessTokenFromCache()
		if err == nil {
			return cachedAccessToken, nil
		}
	}

	directusAuthLoginResponse := dto.DirectusAuthLoginResponse{}

	response, err := client.NewHTTPRequest().
		SetBody(dto.DirectusAuthLoginRequest{
			Email:    config.EnvVariable.DirectusAdminEmail,
			Password: config.EnvVariable.DirectusAdminPassword,
		}).
		SetResult(&directusAuthLoginResponse).
		Post(config.GetDirectusAccessTokenURI())
	if err != nil {
		return "", err
	}

	if !response.IsSuccess() {
		return "", fmt.Errorf("[GetDirectusAccessToken] server response error code: %v", response.StatusCode())
	}
	if !directusAuthLoginResponse.Validate() {
		logger.Error.Printf("[GetDirectusAccessToken] token invalid\n")
		return "", err
	}

	bearerToken := fmt.Sprintf("Bearer %s", directusAuthLoginResponse.Data.AccessToken)

	// Save directus access token back to cache store
	setDirectusAccessTokenToCache(bearerToken, time.Duration(directusAuthLoginResponse.Data.Expires)*time.Millisecond)

	return bearerToken, nil
}
