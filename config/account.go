package config

import (
	"fmt"
	"hubs-cms-go/logger"
	"net/url"
)

func GetDirectusGetAccountURI(mastodonAccount string) string {
	return fmt.Sprintf(`%s/items/account?fields=*,active_avatar.*,liked_rooms.id,liked_rooms.room_id,liked_events.id,liked_events.event_id&filter={"mastodon_account":{"_eq":"%s"}}`, EnvVariable.DirectusBaseURI, mastodonAccount)
}

func GetDirectusPatchAccountURI(accountID string) string {
	return fmt.Sprintf("%s/items/account/%s?fields=*,active_avatar.*,liked_rooms.id,liked_rooms.room_id,liked_events.id,liked_events.event_id", EnvVariable.DirectusBaseURI, accountID)
}

func GetDirectusCreateAccountURI() string {
	return fmt.Sprintf("%s/items/account", EnvVariable.DirectusBaseURI)
}

func genAccountUrl(accountID string) (*url.URL, error) {
	uri, err := url.Parse(EnvVariable.DirectusBaseURI)
	if err != nil {
		logger.Error.Printf("[genAccountUrl] Parse %s error: %v", EnvVariable.DirectusBaseURI, err)
		return nil, err
	}
	if accountID == "" {
		uri.Path = "/items/account"
	} else {
		uri.Path = "/items/account/" + accountID
	}

	return uri, err
}

//
// Type should be "room" or "event"
//
func GetDirectusGetAccountsLikedStuffURI(Type string, offset, limit int64) string {

	uri, err := genAccountUrl("")
	if err != nil {
		return ""
	}

	q := url.Values{}

	q.Add("meta", "filter_count")

	q.Add("fields", fmt.Sprintf("liked_%[1]ss.%[1]s_id", Type))

	q.Add(fmt.Sprintf("filter[liked_%[1]ss][%[1]s_id][_nnull]", Type), "true")

	q.Set("offset", fmt.Sprintf("%v", offset))
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%v", limit))
	}

	uri.RawQuery = q.Encode()

	return uri.String()
}
