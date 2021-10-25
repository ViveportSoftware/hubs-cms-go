package service

import (
	"bytes"
	"fmt"
	"hubs-cms-go/constant"
	"hubs-cms-go/dto"
	"hubs-cms-go/logger"
	"net/http"

	"github.com/go-resty/resty/v2"
	goCache "github.com/patrickmn/go-cache"
)

func directusRequestHandler(request **resty.Request) (response *resty.Response, err error) {
	var directusAccessToken string

	for forceFetch := false; true; {
		directusAccessToken, err = GetDirectusAccessToken(forceFetch)
		if err != nil {
			logger.Error.Printf("[requestHandler] unable to get directus access token error: %v\n", err)
			return
		}

		copyReq := **request
		directusErr := dto.DirectusErrorResponse{}
		response, err = copyReq.
			SetHeader(constant.HeaderAuthorization, directusAccessToken).
			SetError(&directusErr).
			Send()
		logger.Debug.Println("[requestHandler] response:", response)

		if err != nil {
			logger.Error.Printf("[requestHandler] get data error: %v\n", err)
			return
		}

		// 4XX~
		if response.IsError() {
			if response.StatusCode() == http.StatusUnauthorized && !forceFetch {
				forceFetch = true
				continue
			}
			directusErr.Status = response.StatusCode()
			logger.Error.Printf("[requestHandler] server response error status: %v\n", response.StatusCode())
			err = &directusErr
			return
		}

		// 3xx
		if !response.IsSuccess() {
			logger.Error.Printf("[requestHandler] server response error status: %v\n", response.StatusCode())
			err = dto.DirectusErrorResponseFromHttpStstus(response.StatusCode())
			return
		}

		*request = &copyReq
		break
	}

	return
}

func getGraphQLCmd(Type string, items map[string]goCache.Item, Type2 string, items2 map[string]goCache.Item) (count int, cmd string) {

	var b bytes.Buffer
	b.WriteString(`mutation {`)

	temp := `e%[1]d: update_%[2]s_item(id: "%[3]s", data: { like_count: %[4]d }) {like_count}`

	for key, value := range items {
		b.WriteString(fmt.Sprintf(temp, count, Type, key, value.Object.(int64)))
		count++
	}

	temp = `r%[1]d: update_%[2]s_item(id: "%[3]s", data: { like_count: %[4]d }) {like_count}`

	for key, value := range items2 {
		b.WriteString(fmt.Sprintf(temp, count, Type2, key, value.Object.(int64)))
		count++
	}

	b.WriteString(`}`)
	cmd = b.String()

	return
}
