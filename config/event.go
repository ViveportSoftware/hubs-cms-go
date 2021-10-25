package config

import (
	"fmt"
	"hubs-cms-go/logger"
	"net/url"
	"strings"
)

const (
	ev_CLOSED byte = 1 << iota
	ev_OPENED
	ev_SOON
)

func genUrl(eventId string) (*url.URL, error) {
	uri, err := url.Parse(EnvVariable.DirectusBaseURI)
	if err != nil {
		logger.Error.Printf("[genUrl] Parse %s error: %v", EnvVariable.DirectusBaseURI, err)
		return nil, err
	}
	if eventId == "" {
		uri.Path = "/items/event"
	} else {
		uri.Path = "/items/event/" + eventId
	}

	return uri, err
}
func genUrlValues(locale string) url.Values {
	q := url.Values{}
	q.Set("fields", "*,"+
		// "type.id,type.name,"+
		"category.id,category.name,"+
		"hashtags.event_hashtag_id.id,"+
		"hashtags.event_hashtag_id.name,"+
		"hosted_accounts.account_id.id,"+
		"hosted_accounts.account_id.display_name,"+
		"hosted_accounts.account_id.is_admin,"+
		"images.directus_files_id,"+
		"videos.video_id.cover_image,"+
		"videos.video_id.mp4,"+
		"videos.video_id.webm,"+
		"hosts.event_participate_id.id,"+
		"hosts.event_participate_id.name,"+
		"hosts.event_participate_id.description,"+
		"speakers.event_participate_id.id,"+
		"speakers.event_participate_id.name,"+
		"speakers.event_participate_id.description,"+
		"speakers.event_participate_id.image,"+
		"rooms.room_id.id,"+
		"rooms.room_id.title,"+
		"rooms.room_id.description,"+
		"rooms.room_id.gallery,"+
		"rooms.room_id.hubs_id")

	if locale != "" {
		q.Add("fields", "translations.*")
		q.Set("deep[translations][_filter][languages_code][_eq]", locale)

		// q.Add("fields", "type.translations.*")
		// q.Set("deep[type][translations][_filter][languages_code][_eq]", locale)

		q.Add("fields", "category.translations.*")
		q.Set("deep[category][translations][_filter][languages_code][_eq]", locale)

		q.Add("fields", "hosts.event_participate_id.translations.*")
		q.Set("deep[hosts][event_participate_id][translations][_filter][languages_code]", locale)

		q.Add("fields", "speakers.event_participate_id.translations.*")
		q.Set("deep[speakers][event_participate_id][translations][_filter][languages_code]", locale)

		q.Add("fields", "rooms.room_id.translations.*")
		q.Set("deep[rooms][room_id][translations][_filter][languages_code]", locale)
	} else {
		q.Add("fields", "translations.id") // reduce payload
	}

	q.Add("sort", "-is_promoted")
	return q
}

func genUrlValues2() url.Values {
	q := url.Values{}
	q.Set("fields", "id")
	q.Add("fields", "like_count")
	return q
}

func GetDirectusGetEventsURI(locale, status string, offset, limit int64) string {

	uri, err := genUrl("")
	if err != nil {
		return ""
	}
	q := genUrlValues(locale)
	attachEventStatusFilter(q, status)

	q.Set("meta", "filter_count")
	q.Set("offset", fmt.Sprintf("%v", offset))
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%v", limit))
	}
	uri.RawQuery = q.Encode()
	return uri.String()
}

func GetDirectusGetEventURI(eventId string, locale string) string {
	uri, err := genUrl(eventId)
	if err != nil {
		return ""
	}
	q := genUrlValues(locale)
	uri.RawQuery = q.Encode()
	return uri.String()
}

func GetDirectusGetEventURISimple(eventID string) string {
	uri, err := url.Parse(EnvVariable.DirectusBaseURI)
	if err != nil {
		logger.Error.Printf("[GetDirectusGetEventURISimple] Parse %s error: %v\n", EnvVariable.DirectusBaseURI, err)
		return ""
	}
	uri.Path = fmt.Sprintf("/items/event/%s", eventID)
	q := &url.Values{}
	q.Set("fields", "id")
	uri.RawQuery = q.Encode()
	fmt.Printf("GetDirectusGetEventURISimple %s\n", uri.String())
	return uri.String()
}

func GetDirectusGraphQLURI() string {
	return fmt.Sprintf("%s/graphql", EnvVariable.DirectusBaseURI)
}

func attachEventStatusFilter(q url.Values, status string) {
	// distinct status
	state := make(map[string]byte)
	for _, s := range strings.Split(status, "|") {
		if len(s) == 0 {
			continue
		}

		switch s {
		case "closed":
			{
				state[s] = ev_CLOSED
			}
		case "soon":
			{
				state[s] = ev_SOON
			}
		case "opened":
			{
				state[s] = ev_OPENED
			}
		default:
			{
				state["opened"] = ev_OPENED
			}
		}
	}

	// bitwise state
	var iStat byte
	for _, v := range state {
		iStat |= v
	}

	// default
	if iStat == 0 {
		//if not set, use opened
		iStat = ev_OPENED
	}

	switch iStat {
	case ev_OPENED: //opened
		{
			q.Set("filter[start_time][_lte]", "now")
			q.Set("filter[end_time][_gt]", "now")
		}
	case ev_SOON: //soon
		{
			q.Set("filter[start_time][_gt]", "now")
			q.Set("filter[end_time][_nnull]", "true")
		}
	case ev_CLOSED: //closed
		{
			q.Set("filter[start_time][_nnull]", "true")
			q.Set("filter[end_time][_lte]", "now")
		}
	case (ev_OPENED | ev_SOON): //opened + soon
		{
			q.Set("filter[start_time][_nnull]", "true")
			q.Set("filter[end_time][_gt]", "now")
		}
	case (ev_CLOSED | ev_OPENED): //closed + opened
		{
			q.Set("filter[start_time][_lte]", "now")
			q.Set("filter[end_time][_nnull]", "true")
		}
	case (ev_CLOSED | ev_SOON): //closed + soon
		{
			//q.Set("filter[start_time][_gt]", "now")
			//q.Set("filter[end_time][_lte]", "now")
			q.Set("filter", `{"_and":[{"start_time":{"_nnull":true}},{"end_time":{"_nnull":true}},{"_or":[{"start_time":{"_gt":"now"}},{"end_time":{"_lte":"now"}}]}]}`)
		}
	case (ev_CLOSED | ev_OPENED | ev_SOON): // all status
		{
			q.Set("filter[start_time][_nnull]", "true")
			q.Set("filter[end_time][_nnull]", "true")
		}
	}
}
