package service

import (
	"fmt"
	"hubs-cms-go/client"
	"hubs-cms-go/config"
	"hubs-cms-go/constant"
	"hubs-cms-go/dto"
	"hubs-cms-go/errors"
	"hubs-cms-go/logger"
	"hubs-cms-go/utils"
	"mime/multipart"
)

func GetDirectusAvatar(avatarID string, forceFetchDirectusAccessToken bool) (dto.DirectusAvatar, error) {

	directusAccessToken, err := GetDirectusAccessToken(forceFetchDirectusAccessToken)
	if err != nil {
		logger.Error.Printf("[GetDirectusAvatar] unable to get directus access token error: %v\n", err)
		return dto.DirectusAvatar{}, err
	}

	directusGetAvatarResponse := dto.DirectusGetAvatarResponse{}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, directusAccessToken).
		SetResult(&directusGetAvatarResponse).
		Get(config.GetDirectusGetAvatarURI(avatarID))
	if err != nil {
		logger.Error.Printf("[GetDirectusAvatar] %s error: %v\n", config.GetDirectusGetAvatarURI(avatarID), err)
		return dto.DirectusAvatar{}, err
	}

	if !response.IsSuccess() && !forceFetchDirectusAccessToken {
		// Try again with new directus access token
		return GetDirectusAvatar(avatarID, true)
	}

	if !response.IsSuccess() {
		logger.Error.Printf("[GetDirectusAvatar] server response error status: %v", response.StatusCode())
		return dto.DirectusAvatar{}, fmt.Errorf("[GetDirectusAvatar] server response error status: %v", response.StatusCode())
	}

	if !directusGetAvatarResponse.Validate() {
		logger.Error.Printf("[GetDirectusAvatar] server response invalid payload: %v", directusGetAvatarResponse)
		return dto.DirectusAvatar{}, fmt.Errorf("[GetDirectusAvatar] server response invalid payload: %v", directusGetAvatarResponse)
	}

	directusAvatar := dto.DirectusAvatar{
		ID:       directusGetAvatarResponse.Data.ID,
		Snapshot: config.GetDirectusGetAssetURI(directusGetAvatarResponse.Data.Snapshot),
		GLB:      config.GetDirectusGetAssetURI(directusGetAvatarResponse.Data.GLB),
		Owner:    directusGetAvatarResponse.Data.Owner,
		Source:   directusGetAvatarResponse.Data.Source,
		Title:    directusGetAvatarResponse.Data.Title,
		IsPublic: directusGetAvatarResponse.Data.IsPublic,
	}

	return directusAvatar, nil
}

func GetPublicAvatars(start int64, limit int64, forceFetchDirectusAccessToken bool) (dto.GetAvatarsResponse, errors.ErrorInfo) {

	directusAccessToken, err := GetDirectusAccessToken(forceFetchDirectusAccessToken)
	if err != nil {
		logger.Error.Printf("[GetPublicAvatars] unable to get directus access token error: %v\n", err)
		return dto.GetAvatarsResponse{}, errors.InternalError
	}

	directusGetAvatarsResponse := dto.DirectusGetAvatarsResponse{}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, directusAccessToken).
		SetResult(&directusGetAvatarsResponse).
		Get(config.GetDirectusGetPublicAvatarURI(start, limit))
	if err != nil {
		logger.Error.Printf("[GetPublicAvatars] %s error: %v\n", config.GetDirectusGetPublicAvatarURI(start, limit), err)
		return dto.GetAvatarsResponse{}, errors.InternalError
	}

	if !response.IsSuccess() && !forceFetchDirectusAccessToken {
		// Try again with new directus access token
		return GetPublicAvatars(start, limit, true)
	}

	if !response.IsSuccess() {
		logger.Error.Printf("[GetPublicAvatars] server response error status: %v", response.StatusCode())
		return dto.GetAvatarsResponse{}, errors.InternalError
	}

	if !directusGetAvatarsResponse.Validate() {
		logger.Error.Printf("[GetPublicAvatars] server response invalid payload: %v", directusGetAvatarsResponse)
		return dto.GetAvatarsResponse{}, errors.InternalError
	}

	if directusGetAvatarsResponse.Meta.FilterCount == 0 { // no content
		return dto.GetAvatarsResponse{}, errors.ErrorInfo{}
	}

	if directusGetAvatarsResponse.Meta.FilterCount > 0 && start >= directusGetAvatarsResponse.Meta.FilterCount {
		return dto.GetAvatarsResponse{}, errors.AvatarsInvalidStart
	}

	return parseDirectusGetAvatarsResponse(directusGetAvatarsResponse, start, limit), errors.ErrorInfo{}
}

func GetMyAvatars(accountID string, start int64, limit int64, forceFetchDirectusAccessToken bool) (dto.GetAvatarsResponse, errors.ErrorInfo) {

	directusAccessToken, err := GetDirectusAccessToken(forceFetchDirectusAccessToken)
	if err != nil {
		logger.Error.Printf("[GetMyAvatars] unable to get directus access token error: %v\n", err)
		return dto.GetAvatarsResponse{}, errors.InternalError
	}

	directusGetAvatarsResponse := dto.DirectusGetAvatarsResponse{}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, directusAccessToken).
		SetResult(&directusGetAvatarsResponse).
		Get(config.GetDirectusGetMyAvatarURI(accountID, start, limit))
	if err != nil {
		logger.Error.Printf("[GetMyAvatars] %s error: %v\n", config.GetDirectusGetMyAvatarURI(accountID, start, limit), err)
		return dto.GetAvatarsResponse{}, errors.InternalError
	}

	if !response.IsSuccess() && !forceFetchDirectusAccessToken {
		// Try again with new directus access token
		return GetMyAvatars(accountID, start, limit, true)
	}

	if !response.IsSuccess() {
		logger.Error.Printf("[GetMyAvatars] server response error status: %v", response.StatusCode())
		return dto.GetAvatarsResponse{}, errors.InternalError
	}

	if !directusGetAvatarsResponse.Validate() {
		logger.Error.Printf("[GetMyAvatars] server response invalid payload: %v", directusGetAvatarsResponse)
		return dto.GetAvatarsResponse{}, errors.InternalError
	}

	if directusGetAvatarsResponse.Meta.FilterCount > 0 && start >= directusGetAvatarsResponse.Meta.FilterCount {
		return dto.GetAvatarsResponse{}, errors.AvatarsInvalidStart
	}

	return parseDirectusGetAvatarsResponse(directusGetAvatarsResponse, start, limit), errors.ErrorInfo{}
}

func UploadAsset(asset *multipart.FileHeader, forceFetchDirectusAccessToken bool) (string, error) {

	directusAccessToken, err := GetDirectusAccessToken(forceFetchDirectusAccessToken)
	if err != nil {
		logger.Error.Printf("[UploadAsset] unable to get directus access token error: %v\n", err)
		return "", fmt.Errorf("[UploadAsset] unable to get directus access token error: %v", err)
	}

	multipartFile, err := asset.Open()
	if err != nil {
		return "", err
	}

	directusUploadAssetResponse := dto.DirectusUploadAssetResponse{}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, directusAccessToken).
		SetFileReader("unused", asset.Filename, multipartFile).
		SetResult(&directusUploadAssetResponse).
		Post(config.GetDirectusUploadAssetURI())
	if err != nil {
		logger.Error.Printf("[UploadAsset] %s error: %v\n", config.GetDirectusUploadAssetURI(), err)
		return "", fmt.Errorf("[UploadAsset] %s error: %v", config.GetDirectusUploadAssetURI(), err)
	}

	if !response.IsSuccess() && !forceFetchDirectusAccessToken {
		// Try again with new directus access token
		return UploadAsset(asset, true)
	}

	if !response.IsSuccess() {
		logger.Error.Printf("[UploadAsset] server response error status: %v\n", response.StatusCode())
		return "", fmt.Errorf("[UploadAsset] server response error status: %v", response.StatusCode())
	}

	if !directusUploadAssetResponse.Validate() {
		logger.Error.Printf("[UploadAsset] server response invalid payload: %v", directusUploadAssetResponse)
		return "", fmt.Errorf("[UploadAsset] server response invalid payload: %v\n", directusUploadAssetResponse)
	}

	return directusUploadAssetResponse.Data.ID, nil
}

func ImportAsset(assetURL string, forceFetchDirectusAccessToken bool) (string, error) {

	directusAccessToken, err := GetDirectusAccessToken(forceFetchDirectusAccessToken)
	if err != nil {
		logger.Error.Printf("[ImportAsset] unable to get directus access token error: %v\n", err)
		return "", fmt.Errorf("[ImportAsset] unable to get directus access token error: %v", err)
	}

	directusUploadAssetResponse := dto.DirectusUploadAssetResponse{}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, directusAccessToken).
		SetBody(dto.DirectusImportAssetRequest{URL: assetURL}).
		SetResult(&directusUploadAssetResponse).
		Post(config.GetDirectusImportAssetURI())
	if err != nil {
		logger.Error.Printf("[ImportAsset] %s error: %v\n", config.GetDirectusImportAssetURI(), err)
		return "", fmt.Errorf("[ImportAsset] %s error: %v", config.GetDirectusImportAssetURI(), err)
	}

	if !response.IsSuccess() && !forceFetchDirectusAccessToken {
		// Try again with new directus access token
		return ImportAsset(assetURL, true)
	}

	if !response.IsSuccess() {
		logger.Error.Printf("[ImportAsset] server response error status: %v\n", response.StatusCode())
		return "", fmt.Errorf("[ImportAsset] server response error status: %v", response.StatusCode())
	}

	if !directusUploadAssetResponse.Validate() {
		logger.Error.Printf("[ImportAsset] server response invalid payload: %v", directusUploadAssetResponse)
		return "", fmt.Errorf("[ImportAsset] server response invalid payload: %v\n", directusUploadAssetResponse)
	}

	return directusUploadAssetResponse.Data.ID, nil
}

func CreateAvatar(createAvatarRequest dto.DirectusCreateAvatarRequest, forceFetchDirectusAccessToken bool) (dto.DirectusAvatar, error) {

	directusAccessToken, err := GetDirectusAccessToken(forceFetchDirectusAccessToken)
	if err != nil {
		logger.Error.Printf("[CreateAvatar] unable to get directus access token error: %v\n", err)
		return dto.DirectusAvatar{}, fmt.Errorf("[CreateAvatar] unable to get directus access token error: %v", err)
	}

	directusGetAvatarResponse := dto.DirectusGetAvatarResponse{}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, directusAccessToken).
		SetBody(createAvatarRequest).
		SetResult(&directusGetAvatarResponse).
		Post(config.GetDirectusCreateAvatarURI())
	if err != nil {
		logger.Error.Printf("[CreateAvatar] %s error: %v\n", config.GetDirectusImportAssetURI(), err)
		return dto.DirectusAvatar{}, fmt.Errorf("[CreateAvatar] %s error: %v", config.GetDirectusImportAssetURI(), err)
	}

	if !response.IsSuccess() && !forceFetchDirectusAccessToken {
		// Try again with new directus access token
		return CreateAvatar(createAvatarRequest, true)
	}

	if !directusGetAvatarResponse.Validate() {
		logger.Error.Printf("[CreateAvatar] server response invalid payload: %v", directusGetAvatarResponse)
		return dto.DirectusAvatar{}, fmt.Errorf("[CreateAvatar] server response invalid payload: %v", directusGetAvatarResponse)
	}

	directusAvatar := dto.DirectusAvatar{
		ID:       directusGetAvatarResponse.Data.ID,
		Snapshot: config.GetDirectusGetAssetURI(directusGetAvatarResponse.Data.Snapshot),
		GLB:      config.GetDirectusGetAssetURI(directusGetAvatarResponse.Data.GLB),
		Owner:    directusGetAvatarResponse.Data.Owner,
		Source:   directusGetAvatarResponse.Data.Source,
		Title:    directusGetAvatarResponse.Data.Title,
		IsPublic: directusGetAvatarResponse.Data.IsPublic,
	}

	return directusAvatar, nil
}

func GetAvatar(avatarID string, forceFetchDirectusAccessToken bool) (dto.DirectusGetAvatarResponse, error) {

	directusAccessToken, err := GetDirectusAccessToken(forceFetchDirectusAccessToken)
	if err != nil {
		logger.Error.Printf("[GetAvatar] unable to get directus access token error: %v\n", err)
		return dto.DirectusGetAvatarResponse{}, fmt.Errorf("[GetAvatar] unable to get directus access token error: %v", err)
	}

	avatarResponse := dto.DirectusGetAvatarResponse{}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, directusAccessToken).
		SetResult(&avatarResponse).
		Get(config.GetDirectusSingleAvatarURI(avatarID))

	if err != nil {
		logger.Error.Printf("[GetAvatar] %s error: %v\n", config.GetDirectusSingleAvatarURI(avatarID), err)
		return dto.DirectusGetAvatarResponse{}, fmt.Errorf("[GetAvatar] %s error: %v", config.GetDirectusSingleAvatarURI(avatarID), err)
	}

	if !response.IsSuccess() && !forceFetchDirectusAccessToken {
		// Try again with new directus access token
		return GetAvatar(avatarID, true)
	}

	if !avatarResponse.Validate() {
		return dto.DirectusGetAvatarResponse{}, fmt.Errorf("[GetAvatar] server response invalid payload: %v", avatarResponse)
	}

	return avatarResponse, nil
}

func DeleteAvatar(avatarID string, forceFetchDirectusAccessToken bool) error {

	directusAccessToken, err := GetDirectusAccessToken(forceFetchDirectusAccessToken)
	if err != nil {
		logger.Error.Printf("[DeleteAvatar] unable to get directus access token error: %v\n", err)
		return fmt.Errorf("[DeleteAvatar] unable to get directus access token error: %v", err)
	}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, directusAccessToken).
		Delete(config.GetDirectusSingleAvatarURI(avatarID))

	if err != nil {
		logger.Error.Printf("[DeleteAvatar] %s error: %v\n", config.GetDirectusSingleAvatarURI(avatarID), err)
		return fmt.Errorf("[DeleteAvatar] %s error: %v", config.GetDirectusSingleAvatarURI(avatarID), err)
	}

	if !response.IsSuccess() && !forceFetchDirectusAccessToken {
		// Try again with new directus access token
		return DeleteAvatar(avatarID, true)
	}

	return nil
}

func parseDirectusGetAvatarsResponse(origin dto.DirectusGetAvatarsResponse, start int64, limit int64) dto.GetAvatarsResponse {

	result := dto.GetAvatarsResponse{}
	result.Results = []dto.DirectusAvatar{}

	for _, d := range origin.Data {
		result.Results = append(result.Results, dto.DirectusAvatar{
			ID:       d.ID,
			Snapshot: config.GetDirectusGetAssetURI(d.Snapshot),
			GLB:      config.GetDirectusGetAssetURI(d.GLB),
			Owner:    d.Owner,
			Source:   d.Source,
			Title:    d.Title,
			IsPublic: d.IsPublic,
		})
	}

	if limit > 0 && start > 0 {
		result.Pages.Prev = fmt.Sprintf("/api/hubs-cms/v1/my-avatars?start=%d&limit=%d", utils.Max(0, start-limit), limit)
	}

	if limit > 0 && origin.Meta.FilterCount > start+limit {
		result.Pages.Next = fmt.Sprintf("/api/hubs-cms/v1/my-avatars?start=%d&limit=%d", start+limit, limit)
	}

	return result
}
