package service

import (
	"hubs-cms-go/client"
	"hubs-cms-go/config"
	"hubs-cms-go/constant"
	"hubs-cms-go/dto"
	"hubs-cms-go/errors"
	"hubs-cms-go/logger"

	"github.com/go-resty/resty/v2"
	goCache "github.com/patrickmn/go-cache"
)

func GetDirectusEvents(locale, status string, start, limit int64) (ret []dto.DirectusEventResponseData, total int64, err error) {
	directusResponse := dto.DirectusGetResponse{Data: &ret}
	request := client.NewHTTPRequest().SetResult(&directusResponse)
	request.Method = resty.MethodGet
	request.URL = config.GetDirectusGetEventsURI(locale, status, start, limit)

	if _, err = directusRequestHandler(&request); err != nil {
		return
	}
	total = directusResponse.Meta.FilterCount
	return
}

func GetDirectusEvent(eventID, locale string) (ret dto.DirectusEventResponseData, err error) {
	request := client.NewHTTPRequest().SetResult(&dto.DirectusGetResponse{Data: &ret})
	request.Method = resty.MethodGet
	request.URL = config.GetDirectusGetEventURI(eventID, locale)

	if _, err = directusRequestHandler(&request); err != nil {
		return
	}
	return
}

func PatchDirectusEvent(eventID string, patchBody interface{}) (err error) {
	request := client.NewHTTPRequest().
		SetHeader("Content-Type", "application/json").
		SetBody(patchBody)
	request.Method = resty.MethodPatch
	request.URL = config.GetDirectusGetEventURISimple(eventID)

	_, err = directusRequestHandler(&request)
	return
}

func PostDirectusEventViewCount(eventInfo dto.DirectusEventResponseData, locale string, isRetried bool) (dto.DirectusEventResponseData, errors.ErrorInfo) {

	ret := dto.DirectusEventResponseData{}

	directusAccessToken, tokenErr := GetDirectusAccessToken(isRetried)
	if tokenErr != nil {
		logger.Error.Printf("[PostDirectusEventViewCount] fail to get token. Turn on retry flag: %v\n", isRetried)
		if !isRetried {
			return PostDirectusEventViewCount(eventInfo, locale, true)
		}
		logger.Error.Printf("[PostDirectusEventViewCount] unable to get directus access token error: %v\n", tokenErr)
		return dto.DirectusEventResponseData{}, errors.InternalError
	}

	addNumber, convertNumberErr := eventInfo.ViewCount.Int64()
	newViewCount := addNumber + 1

	if convertNumberErr != nil {
		logger.Error.Printf("[PostDirectusEventViewCount] unable convert to int64 error: %v\n", convertNumberErr)
		return dto.DirectusEventResponseData{}, errors.InternalError
	}

	response, err := client.NewHTTPRequest().
		SetHeader(constant.HeaderAuthorization, directusAccessToken).
		SetBody(dto.RoomDataIncreaseViewCountRequest{
			ViewCount: newViewCount,
		}).
		SetResult(&dto.DirectusGetResponse{Data: &ret}).
		Patch(config.GetDirectusGetEventURI(eventInfo.ID, locale))

	if err != nil {
		logger.Error.Printf("[PostDirectusEventViewCount] request update room count error: %v\n", err)
		return dto.DirectusEventResponseData{}, errors.InternalError
	}

	if !response.IsSuccess() {
		logger.Error.Printf("[PostDirectusEventViewCount] fail to update view count. Turn on retry flag: %v\n", isRetried)
		if !isRetried {
			return PostDirectusEventViewCount(eventInfo, locale, true)
		}
		logger.Error.Printf("[PostDirectusEventViewCount] error: %v\n", err)
		return dto.DirectusEventResponseData{}, errors.InternalError
	}

	return ret, errors.ErrorInfo{}
}

func BackupLikeCount(Type string, items map[string]goCache.Item, Type2 string, items2 map[string]goCache.Item) (eventCount, totalLikes int64) {

	_, cmd := getGraphQLCmd(Type, items, Type2, items2)
	result, _ := SendDirectusGraphQLCmd(cmd)

	//
	// check server response string
	// example string:
	//  {"data":{"e0":{"like_count":4},"e1":{"like_count":1},"e2":{"like_count":4},"e3":{"like_count":2}}}
	// fail case:
	//  {"data":{"e0":null,"e1":null,"e2":null,"e3":null}}
	//
	for _, v := range result {
		//
		// conver value type to map[string]interface{}
		// v should be like this {"like_count":4}
		// but nil means err
		//
		switch n := v.(type) {
		case map[string]interface{}:
			// treat non-nil as scuccess
			if n == nil {
				continue
			}
			if val, ok := n["like_count"]; ok {
				eventCount++
				totalLikes += int64(val.(float64))
			}
		}
	}
	return
}
