package config

import (
	"fmt"
	"hubs-cms-go/logger"
	"net/url"
)

func GetHubsURL(hubsID string) (string, error) {
	uri, err := url.Parse(EnvVariable.HubsBaseURI)
	if err != nil {
		logger.Error.Printf("[GetHubsURL] Parse %s error: %v\n", EnvVariable.HubsBaseURI, err)
		return "", err
	}

	uri.Path = hubsID
	return uri.String(), nil
}

func GetDirectusGetMyRoomListURI(accountID, locale string, offset, limit int64) string {
	uri, err := url.Parse(EnvVariable.DirectusBaseURI)
	if err != nil {
		logger.Error.Printf("[GetDirectusGetMyRoomListURI] Parse %s error: %v\n", EnvVariable.DirectusBaseURI, err)
		return ""
	}

	uri.Path = "/items/room"
	q := &url.Values{}
	q.Set("filter[owner]", accountID) // filter by account id
	q.Set("fields", "*,gallery.id,events.event_id,nft_contract.*")
	attachPaging(attachTranslation(q, locale), offset, limit)
	uri.RawQuery = q.Encode() // sort queries by key
	return uri.String()
}

func GetDirectusGetRoomListURI(pHasNFT *bool, hubsID, locale string, offset, limit int64) string {
	uri, err := url.Parse(EnvVariable.DirectusBaseURI)
	if err != nil {
		logger.Error.Printf("[GetDirectusGetRoomListURI] Parse %s error: %v\n", EnvVariable.DirectusBaseURI, err)
		return ""
	}

	uri.Path = "/items/room"
	q := &url.Values{}
	q.Set("fields", "*,gallery.id,events.event_id,nft_contract.*")
	if len(hubsID) == 0 {
		q.Set("filter[is_public]", "true")
		if pHasNFT != nil {
			q.Set("filter[has_nft]", fmt.Sprintf("%v", *pHasNFT))
		}
	} else {
		q.Set("filter[hubs_id]", hubsID)
	}
	attachPaging(attachTranslation(q, locale), offset, limit)
	uri.RawQuery = q.Encode() // sort queries by key
	return uri.String()
}

func GetDirectusGetRoomURI(roomID, locale string) string {
	uri, err := url.Parse(EnvVariable.DirectusBaseURI)
	if err != nil {
		logger.Error.Printf("[GetDirectusGetRoomURI] Parse %s error: %v\n", EnvVariable.DirectusBaseURI, err)
		return ""
	}
	uri.Path = fmt.Sprintf("/items/room/%s", roomID)
	q := &url.Values{}
	q.Set("fields", "*,gallery.id,events.event_id,nft_contract.*")
	attachTranslation(q, locale)
	uri.RawQuery = q.Encode()
	return uri.String()
}

func GetDirectusGetRoomURISimple(roomID string) string {
	uri, err := url.Parse(EnvVariable.DirectusBaseURI)
	if err != nil {
		logger.Error.Printf("[GetDirectusGetRoomURISimple] Parse %s error: %v\n", EnvVariable.DirectusBaseURI, err)
		return ""
	}
	uri.Path = fmt.Sprintf("/items/room/%s", roomID)
	q := &url.Values{}
	q.Set("fields", "id")
	uri.RawQuery = q.Encode()
	return uri.String()
}

func attachPaging(q *url.Values, offset, limit int64) *url.Values {
	if q == nil {
		return q
	}

	setQuery(q, "meta", "filter_count")
	setQuery(q, "offset", offset)

	if limit > 0 {
		setQuery(q, "limit", limit)
	}
	return q
}

func attachTranslation(q *url.Values, locale string) *url.Values {
	if q == nil {
		return q
	}

	if !q.Has("fields") {
		q.Set("fields", "")
	}

	if len(locale) > 0 {
		q.Add("fields", "translations.*") // add field for translation
		setQuery(q, "deep[translations][_filter][languages_code]", locale)
	} else {
		q.Add("fields", "translations.id") // reduce payload
	}

	return q
}

func setQuery(q *url.Values, k string, v interface{}) {
	if q.Has(k) {
		logger.Warn.Printf("[setQuery] Override %v\n", k)
	}
	q.Set(k, fmt.Sprintf("%v", v))
}
