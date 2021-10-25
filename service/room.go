package service

import (
	"hubs-cms-go/client"
	"hubs-cms-go/config"
	"hubs-cms-go/constant"
	"hubs-cms-go/dto"
	"hubs-cms-go/errors"
	"hubs-cms-go/logger"

	"github.com/go-resty/resty/v2"
)

func PatchDirectusRoom(roomID string, patchBody interface{}) (err error) {
	request := client.NewHTTPRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(patchBody)
	request.Method = resty.MethodPatch
	request.URL = config.GetDirectusGetRoomURISimple(roomID)

	_, err = directusRequestHandler(&request)
	return
}

func GetDirectusMyRoomList(accountID, locale string, start, limit int64) (ret []dto.DierctusRoomData, total int64, err error) {
	directusResponse := dto.DirectusGetResponse{Data: &ret}
	request := client.NewHTTPRequest().SetResult(&directusResponse)
	request.Method = resty.MethodGet
	request.URL = config.GetDirectusGetMyRoomListURI(accountID, locale, start, limit)

	if _, err = directusRequestHandler(&request); err != nil {
		return
	}

	total = directusResponse.Meta.FilterCount
	if len(locale) > 0 {
		for i := range ret {
			ret[i].UpdateTranslation()
		}
	}
	return
}

func GetDirectusRoomList(pHasNFT *bool, hubsID, locale string, start, limit int64) (ret []dto.DierctusRoomData, total int64, err error) {
	directusResponse := dto.DirectusGetResponse{Data: &ret}
	request := client.NewHTTPRequest().SetResult(&directusResponse)
	request.Method = resty.MethodGet
	request.URL = config.GetDirectusGetRoomListURI(pHasNFT, hubsID, locale, start, limit)

	if _, err = directusRequestHandler(&request); err != nil {
		return
	}

	total = directusResponse.Meta.FilterCount
	if len(locale) > 0 {
		for i := range ret {
			ret[i].UpdateTranslation()
		}
	}
	return
}

func GetDirectusRoomWithCustomData(roomID, locale string, customData interface{}) (err error) {
	request := client.NewHTTPRequest().SetResult(&dto.DirectusGetResponse{Data: customData})
	request.Method = resty.MethodGet
	request.URL = config.GetDirectusGetRoomURI(roomID, locale)

	if _, err = directusRequestHandler(&request); err != nil {
		return
	}

	return
}

func GetDirectusRoom(roomID, locale string) (ret dto.DierctusRoomData, err error) {
	request := client.NewHTTPRequest().SetResult(&dto.DirectusGetResponse{Data: &ret})
	request.Method = resty.MethodGet
	request.URL = config.GetDirectusGetRoomURI(roomID, locale)

	if _, err = directusRequestHandler(&request); err != nil {
		return
	}

	// update translation
	if len(locale) > 0 {
		ret.UpdateTranslation()
	}
	return
}

func PostRoomViewCount(roomData dto.DierctusRoomData, locale string, isRetried bool) (dto.DierctusRoomData, errors.ErrorInfo) {
	ret := dto.DierctusRoomData{}

	directusAccessToken, tokenErr := GetDirectusAccessToken(isRetried)
	if tokenErr != nil {
		logger.Error.Printf("[PostRoomViewCount] fail to get token. Turn on retry flag: %v\n", isRetried)
		if !isRetried {
			return PostRoomViewCount(roomData, locale, true)
		}
		logger.Error.Printf("[PostRoomViewCount] unable to get directus access token error: %v\n", tokenErr)
		return dto.DierctusRoomData{}, errors.InternalError
	}

	addNumber, convertNumberErr := roomData.ViewCount.Int64()
	newViewCount := addNumber + 1
	if convertNumberErr != nil {
		logger.Error.Printf("[PostRoomViewCount] unable to convert number error: %v\n", convertNumberErr)
		return dto.DierctusRoomData{}, errors.InternalError
	}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, directusAccessToken).
		SetBody(dto.RoomDataIncreaseViewCountRequest{
			ViewCount: newViewCount,
		}).
		SetResult(&dto.DirectusGetResponse{Data: &ret}).
		Patch(config.GetDirectusGetRoomURI(roomData.ID, locale))

	if err != nil {
		logger.Error.Printf("[PostRoomViewCount] request update room count error: %v\n", err)
		return dto.DierctusRoomData{}, errors.InternalError
	}

	if !response.IsSuccess() {
		logger.Error.Printf("[PostRoomViewCount] fail to update view count. Turn on retry flag: %v\n", isRetried)
		if !isRetried {
			return PostRoomViewCount(roomData, locale, true)
		}
		logger.Error.Printf("[PostRoomViewCount] error: %v\n", err)
		return dto.DierctusRoomData{}, errors.InternalError
	}

	return ret, errors.ErrorInfo{}
}
