package service

import (
	"fmt"
	"hubs-cms-go/client"
	"hubs-cms-go/config"
	"hubs-cms-go/constant"
	"hubs-cms-go/dto"
	"hubs-cms-go/logger"
	"net/http"

	"github.com/go-resty/resty/v2"
)

func GetDirectusAccountData(mastodonAccount string) (ret dto.DirectusAccountResponseData, err error) {
	var data []dto.DirectusAccountResponseData

	request := client.NewHTTPRequest().
		SetResult(&dto.DirectusGetResponse{Data: &data})
	request.Method = resty.MethodGet
	request.URL = config.GetDirectusGetAccountURI(mastodonAccount)

	if _, err = directusRequestHandler(&request); err != nil {
		logger.Error.Printf("[GetDirectusAccountData] get room data error: %v\n", err)
		return
	}

	if l := len(data); l > 0 {
		if l > 1 {
			logger.Warn.Printf("{GetDirectusAccountData} find %v accounts\n", l)
		}
		ret = data[0]
	}

	return
}

func GetDirectusAccount(mastodonAccount string, forceFetchDirectusAccessToken bool) (dto.DirectusAccount, error) {

	directusAccessToken, err := GetDirectusAccessToken(forceFetchDirectusAccessToken)
	if err != nil {
		logger.Error.Printf("[GetDirectusAccount] unable to get directus access token error: %v\n", err)
		return dto.DirectusAccount{}, err
	}

	directusGetAccountResponse := dto.DirectusGetAccountResponse{}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, directusAccessToken).
		SetResult(&directusGetAccountResponse).
		Get(config.GetDirectusGetAccountURI(mastodonAccount))
	if err != nil {
		logger.Error.Printf("[GetDirectusAccount] %s error: %v\n", config.GetDirectusGetAccountURI(mastodonAccount), err)
		return dto.DirectusAccount{}, err
	}

	if !response.IsSuccess() && !forceFetchDirectusAccessToken {
		// Try again with new directus access token
		return GetDirectusAccount(mastodonAccount, true)
	}

	if !response.IsSuccess() {
		logger.Error.Printf("[GetDirectusAccount] server response error status: %v\n", response.StatusCode())
		return dto.DirectusAccount{}, fmt.Errorf("[GetDirectusAccount] server response error status: %v", response.StatusCode())
	}

	if len(directusGetAccountResponse.Data) == 0 {
		// No record found, return empty result with nil error to hint the caller to create a new directus account
		logger.Error.Printf("[GetDirectusAccount] server response payload length incorrect\n")
		return dto.DirectusAccount{}, nil
	}

	if !directusGetAccountResponse.Validate() {
		logger.Error.Printf("[GetDirectusAccount] server response invalid payload: %v\n", directusGetAccountResponse)
		return dto.DirectusAccount{}, fmt.Errorf("[GetDirectusAccount] server response invalid payload: %v", directusGetAccountResponse)
	}

	directusAccount := dto.DirectusAccount{
		ID:              directusGetAccountResponse.Data[0].ID,
		MastodonAccount: directusGetAccountResponse.Data[0].MastodonAccount,
		MastodonAvatar:  directusGetAccountResponse.Data[0].MastodonAvatar,
		DisplayName:     directusGetAccountResponse.Data[0].DisplayName,
		IsAdmin:         directusGetAccountResponse.Data[0].IsAdmin,
	}

	if len(directusGetAccountResponse.Data[0].ActiveAvatar.ID) > 0 {
		directusAccount.ActiveAvatar = &dto.DirectusAvatar{
			ID:       directusGetAccountResponse.Data[0].ActiveAvatar.ID,
			Snapshot: config.GetDirectusGetAssetURI(directusGetAccountResponse.Data[0].ActiveAvatar.Snapshot),
			GLB:      config.GetDirectusGetAssetURI(directusGetAccountResponse.Data[0].ActiveAvatar.GLB),
			Owner:    directusGetAccountResponse.Data[0].ActiveAvatar.Owner,
			Source:   directusGetAccountResponse.Data[0].ActiveAvatar.Source,
			Title:    directusGetAccountResponse.Data[0].ActiveAvatar.Title,
			IsPublic: directusGetAccountResponse.Data[0].ActiveAvatar.IsPublic,
		}
	}

	directusAccount.LikedRooms = make([]string, len(directusGetAccountResponse.Data[0].LikedRooms))
	for i := range directusGetAccountResponse.Data[0].LikedRooms {
		directusAccount.LikedRooms[i] = directusGetAccountResponse.Data[0].LikedRooms[i].RoomID
	}

	directusAccount.LikedEvents = make([]string, len(directusGetAccountResponse.Data[0].LikedEvents))
	for i := range directusGetAccountResponse.Data[0].LikedEvents {
		directusAccount.LikedEvents[i] = directusGetAccountResponse.Data[0].LikedEvents[i].EventID
	}

	return directusAccount, nil
}

func CreateDirectusAccount(mastodonAccountInfo dto.MastodonVerifyCredentialsResponse, forceFetchDirectusAccessToken bool) (dto.DirectusAccount, error) {

	directusAccessToken, err := GetDirectusAccessToken(false)
	if err != nil {
		logger.Error.Printf("[CreateDirectusAccount] unable to get directus access token error: %v\n", err)
		return dto.DirectusAccount{}, err
	}

	directusUpsertAccountResponse := dto.DirectusUpsertAccountResponse{}
	directusErr := dto.DirectusErrorResponse{}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, directusAccessToken).
		SetBody(dto.DirectusCreateAccountRequest{
			MastodonAccount: mastodonAccountInfo.MastodonAccount,
			MastodonAvatar:  mastodonAccountInfo.MastodonAvatar,
			DisplayName:     mastodonAccountInfo.DisplayName,
			IsAdmin:         false,
		}).
		SetResult(&directusUpsertAccountResponse).
		SetError(&directusErr).
		Post(config.GetDirectusCreateAccountURI())
	if err != nil {
		logger.Error.Printf("[CreateDirectusAccount] request error: %v\n", err)
		return dto.DirectusAccount{}, err
	}

	if !response.IsSuccess() {
		switch response.StatusCode() {
		case http.StatusUnauthorized:
			{
				if !forceFetchDirectusAccessToken {
					// Try again with new directus access token
					return CreateDirectusAccount(mastodonAccountInfo, true)
				}
			}
		case http.StatusBadRequest:
			{
				directusErr.Status = http.StatusBadRequest
				return dto.DirectusAccount{}, &directusErr
			}
		}

		logger.Error.Printf("[CreateDirectusAccount] server response error status: %v\n", response.StatusCode())
		return dto.DirectusAccount{}, fmt.Errorf("[CreateDirectusAccount] server response error status: %v", response.StatusCode())
	}

	if !directusUpsertAccountResponse.Validate() {
		logger.Error.Printf("[CreateDirectusAccount] server response invalid payload: %v\n", directusUpsertAccountResponse)
		return dto.DirectusAccount{}, fmt.Errorf("[CreateDirectusAccount] server response invalid payload: %v", directusUpsertAccountResponse)
	}

	directusAccount := dto.DirectusAccount{
		ID:              directusUpsertAccountResponse.Data.ID,
		MastodonAccount: directusUpsertAccountResponse.Data.MastodonAccount,
		MastodonAvatar:  directusUpsertAccountResponse.Data.MastodonAvatar,
		DisplayName:     directusUpsertAccountResponse.Data.DisplayName,
		IsAdmin:         directusUpsertAccountResponse.Data.IsAdmin,
	}

	return directusAccount, nil
}

func PatchDirectusAccount(accountID string, patchBody interface{}, forceFetchDirectusAccessToken bool) (dto.DirectusAccount, error) {

	directusAccessToken, err := GetDirectusAccessToken(false)
	if err != nil {
		logger.Error.Printf("[PatchDirectusAccount] unable to get directus access token error: %v\n", err)
		return dto.DirectusAccount{}, err
	}

	directusUpsertAccountResponse := dto.DirectusUpsertAccountResponse{}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, directusAccessToken).
		SetHeader("Content-Type", "application/json").
		SetBody(patchBody).
		SetResult(&directusUpsertAccountResponse).
		Patch(config.GetDirectusPatchAccountURI(accountID))
	if err != nil {
		logger.Error.Printf("[PatchDirectusAccount] request error: %v\n", err)
		return dto.DirectusAccount{}, err
	}

	if !response.IsSuccess() && !forceFetchDirectusAccessToken {
		// Try again with new directus access token
		return PatchDirectusAccount(accountID, patchBody, true)
	}

	if !response.IsSuccess() {
		logger.Error.Printf("[PatchDirectusAccount] server response error status: %v\n", response.StatusCode())
		return dto.DirectusAccount{}, fmt.Errorf("[PatchDirectusAccount] server response error status: %v", response.StatusCode())
	}

	if !directusUpsertAccountResponse.Validate() {
		logger.Error.Printf("[PatchDirectusAccount] server response invalid payload: %v\n", directusUpsertAccountResponse)
		return dto.DirectusAccount{}, fmt.Errorf("[PatchDirectusAccount] server response invalid payload: %v", directusUpsertAccountResponse)
	}

	directusAccount := dto.DirectusAccount{
		ID:              directusUpsertAccountResponse.Data.ID,
		MastodonAccount: directusUpsertAccountResponse.Data.MastodonAccount,
		MastodonAvatar:  directusUpsertAccountResponse.Data.MastodonAvatar,
		DisplayName:     directusUpsertAccountResponse.Data.DisplayName,
		IsAdmin:         directusUpsertAccountResponse.Data.IsAdmin,
	}

	if len(directusUpsertAccountResponse.Data.ActiveAvatar.ID) > 0 {
		directusAccount.ActiveAvatar = &dto.DirectusAvatar{
			ID:       directusUpsertAccountResponse.Data.ActiveAvatar.ID,
			Snapshot: config.GetDirectusGetAssetURI(directusUpsertAccountResponse.Data.ActiveAvatar.Snapshot),
			GLB:      config.GetDirectusGetAssetURI(directusUpsertAccountResponse.Data.ActiveAvatar.GLB),
			Owner:    directusUpsertAccountResponse.Data.ActiveAvatar.Owner,
			Source:   directusUpsertAccountResponse.Data.ActiveAvatar.Source,
			Title:    directusUpsertAccountResponse.Data.ActiveAvatar.Title,
			IsPublic: directusUpsertAccountResponse.Data.ActiveAvatar.IsPublic,
		}
	}

	return directusAccount, nil
}

func GetAccountsLikedStuff(Type string, offset, limit int64) (ret []dto.DirectusGetAccountLikesResponseData, filterCount int64, err error) {

	logger.Debug.Printf("[GetAccountsLiked%vs] offset=%v, limit=%v\n", Type, offset, limit)

	directusResponse := dto.DirectusGetResponse{Data: &ret}

	request := client.NewHTTPRequest().SetResult(&directusResponse)
	request.Method = resty.MethodGet

	request.URL = config.GetDirectusGetAccountsLikedStuffURI(Type, offset, limit)

	if _, err = directusRequestHandler(&request); err != nil {
		logger.Error.Printf("[GetAccountsLiked%vs] %v\n", Type, err)
	}

	filterCount = directusResponse.Meta.FilterCount

	return
}
