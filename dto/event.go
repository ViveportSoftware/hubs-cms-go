package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"hubs-cms-go/cache"
	"hubs-cms-go/config"
	"hubs-cms-go/logger"
	"hubs-cms-go/utils"
	"time"
)

type DirectusGetEventsResponse struct {
	Meta DirectusMeta                `json:"meta"`
	Data []DirectusEventResponseData `json:"data"`
}
type DirectusGetEventByIDResponse struct {
	Data DirectusEventResponseData `json:"data"`
}
type ParticipateIDTranslations struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type ParticipateID struct {
	ID           string                      `json:"id"`
	Name         string                      `json:"name"`
	Description  string                      `json:"description"`
	Translations []ParticipateIDTranslations `json:"translations"`
}
type DirectusHost struct {
	ParticipateID ParticipateID `json:"event_participate_id"`
}
type SpeakersParticipateID struct {
	ID           string                      `json:"id"`
	Name         string                      `json:"name"`
	Description  string                      `json:"description"`
	Image        string                      `json:"image"`
	Translations []ParticipateIDTranslations `json:"translations"`
}
type DirectusSpeaker struct {
	ParticipateID SpeakersParticipateID `json:"event_participate_id"`
}
type RoomIDTranslations struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
type RoomID struct {
	ID           string               `json:"id"`
	Title        string               `json:"title"`
	Gallery      string               `json:"gallery"`
	Description  string               `json:"description"`
	HubsID       string               `json:"hubs_id"`
	Translations []RoomIDTranslations `json:"translations"`
}
type DirectusRoom struct {
	RoomID RoomID `json:"room_id"`
}
type TypeTranslations struct {
	Name string `json:"name"`
}
type DirectusType struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Translations []TypeTranslations `json:"translations"`
}
type Translations struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Agenda      string `json:"agenda"`
}
type DirectusFilesID struct {
	ID string `json:"directus_files_id"`
}

func (f *DirectusFilesID) Validate() bool {
	return f != nil && f.ID != ""
}

type AccountID struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	IsAdmin     bool   `json:"is_admin"`
}
type HostedAccount struct {
	AccountID AccountID `json:"account_id"`
}
type HashtagID struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type DirectusHashtag struct {
	HashtagID HashtagID `json:"event_hashtag_id"`
}
type CategoryTranslations struct {
	Name string `json:"name"`
}
type DirectusCategory struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Translations []CategoryTranslations `json:"translations"`
}
type DirectusVideo struct {
	DirectusVideo DirectusVideoID `json:"video_id"`
}
type DirectusVideoID struct {
	CoverImage string `json:"cover_image"`
	Mp4        string `json:"mp4"`
	Webm       string `json:"webm"`
}

func (e DirectusVideo) NewVideo() (Video, error) {
	video := Video{}
	if len(e.DirectusVideo.CoverImage) > 0 {
		video.CoverImage = config.GetDirectusGetAssetURI(e.DirectusVideo.CoverImage)
	}
	if len(e.DirectusVideo.Mp4) > 0 {
		video.Mp4 = config.GetDirectusGetAssetURI(e.DirectusVideo.Mp4)
	}
	if len(e.DirectusVideo.Webm) > 0 {
		video.Webm = config.GetDirectusGetAssetURI(e.DirectusVideo.Webm)
	}
	if len(video.Mp4) == 0 && len(video.Webm) == 0 {
		return video, errors.New("video's url invalid")
	}
	return video, nil
}

type DirectusEventResponseData struct {
	ID             string            `json:"id"`
	Gallery        string            `json:"gallery"`
	Title          string            `json:"title"`
	Description    string            `json:"description"`
	Agenda         string            `json:"agenda"`
	IsPromoted     bool              `json:"is_promoted"`
	StartTime      time.Time         `json:"start_time"`
	EndTime        time.Time         `json:"end_time"`
	IsLiked        bool              `json:"is_liked"`
	LikeCount      json.Number       `json:"like_count"`
	ViewCount      json.Number       `json:"view_count"`
	Hosts          []DirectusHost    `json:"hosts"`
	Speakers       []DirectusSpeaker `json:"speakers"`
	Rooms          []DirectusRoom    `json:"rooms"`
	Translations   []Translations    `json:"translations"`
	Images         []DirectusFilesID `json:"images"`
	Videos         []DirectusVideo   `json:"videos"`
	HostedAccounts []HostedAccount   `json:"hosted_accounts"`
	Hashtags       []DirectusHashtag `json:"hashtags"`
	Category       DirectusCategory  `json:"category"`
	// Type           DirectusType            `json:"type"`
}
type DirectusEventResponseData2 struct {
	ID        string      `json:"id"`
	LikeCount json.Number `json:"like_count"`
}
type GetEventsResponse struct {
	Results []GetEventResponse `json:"results"`
	Pages   Page               `json:"pages"`
}
type Host struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}
type Speaker struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}
type Room struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Gallery     string `json:"gallery"`
	Description string `json:"description"`
	HubsURL     string `json:"hubs_url"`
}
type Type struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
type Hashtag struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
type Category struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
type Video struct {
	CoverImage string `json:"cover_image"`
	Mp4        string `json:"mp4"`
	Webm       string `json:"webm"`
}
type GetEventResponse struct {
	ID          string      `json:"id"`
	Gallery     string      `json:"gallery"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Agenda      string      `json:"agenda"`
	StartTime   time.Time   `json:"start_time"`
	EndTime     time.Time   `json:"end_time"`
	IsLiked     bool        `json:"is_liked"`
	IsPromoted  bool        `json:"is_promoted"`
	LikeCount   json.Number `json:"like_count"`
	ViewCount   json.Number `json:"view_count"`
	Hosts       []Host      `json:"hosts"`
	Speakers    []Speaker   `json:"speakers"`
	Rooms       []Room      `json:"rooms"`
	Images      []string    `json:"images"`
	Videos      []Video     `json:"videos"`
	Hashtags    []Hashtag   `json:"hashtags"`
	Category    Category    `json:"category"`
	// Type        Type      `json:"type"`
}

func (d *GetEventResponse) parse(rsp DirectusEventResponseData, account *DirectusAccountResponseData) {

	var data = rsp
	d.ID = data.ID
	d.Title = data.Title
	d.Agenda = data.Agenda
	d.Description = data.Description
	d.StartTime = data.StartTime
	d.EndTime = data.EndTime
	d.ViewCount = data.ViewCount
	//d.LikeCount = data.LikeCount
	d.IsPromoted = data.IsPromoted

	if item, ok := cache.EventLikes.Get(data.ID); ok {
		if likes, ok := item.(int64); ok {
			d.LikeCount = json.Number(fmt.Sprintf("%v", likes))
		}
	}

	if len(data.Translations) > 0 {
		if data.Translations[0].Title != "" {
			d.Title = data.Translations[0].Title
		}
		if data.Translations[0].Agenda != "" {
			d.Agenda = data.Translations[0].Agenda
		}
		if data.Translations[0].Description != "" {
			d.Description = data.Translations[0].Description
		}
	}

	d.Category = Category{ID: data.Category.ID, Value: data.Category.Name}
	if len(data.Category.Translations) > 0 && data.Category.Translations[0].Name != "" {
		d.Category.Value = data.Category.Translations[0].Name
	}

	if len(data.Gallery) > 0 {
		d.Gallery = config.GetDirectusGetAssetURI(data.Gallery)
	}
	// d.Type.ID = data.Type.ID
	// if len(data.Type.Translations) > 0 {
	// 	d.Type.Value = data.Type.Translations[0].Name
	// } else {
	// 	d.Type.Value = data.Type.Name
	// }

	d.Speakers = []Speaker{}
	for _, s := range data.Speakers {
		speaker := Speaker{ID: s.ParticipateID.ID}
		if len(s.ParticipateID.Translations) > 0 {
			speaker.DisplayName = s.ParticipateID.Translations[0].Name
		} else {
			speaker.DisplayName = s.ParticipateID.Name
		}
		if len(s.ParticipateID.Image) > 0 {
			speaker.ImageURL = config.GetDirectusGetAssetURI(s.ParticipateID.Image)
		}
		d.Speakers = append(d.Speakers, speaker)
	}

	d.Hosts = []Host{}
	for _, h := range data.Hosts {
		host := Host{ID: h.ParticipateID.ID}
		if len(h.ParticipateID.Translations) > 0 {
			host.DisplayName = h.ParticipateID.Translations[0].Name
		} else {
			host.DisplayName = h.ParticipateID.Name
		}
		d.Hosts = append(d.Hosts, host)
	}

	d.Rooms = []Room{}
	for _, r := range data.Rooms {
		room := Room{
			ID:          r.RoomID.ID,
			Title:       r.RoomID.Title,
			Description: r.RoomID.Description,
			Gallery:     "",
		}

		if len(r.RoomID.Gallery) > 0 {
			room.Gallery = config.GetDirectusGetAssetURI(r.RoomID.Gallery)
		}

		if len(r.RoomID.Translations) > 0 {
			if title := r.RoomID.Translations[0].Title; len(title) > 0 {
				room.Title = title
			}
			if desc := r.RoomID.Translations[0].Description; len(desc) > 0 {
				room.Description = desc
			}
		}

		if len(r.RoomID.HubsID) > 0 {
			if hubsURL, err := config.GetHubsURL(r.RoomID.HubsID); err != nil {
				logger.Error.Printf("[parse] Get hubs url error: %v\n", err)
			} else {
				room.HubsURL = hubsURL
			}
		} else {
			logger.Error.Printf("[parse] Missing hubs id for: %v\n", r.RoomID.ID)
		}

		d.Rooms = append(d.Rooms, room)
	}

	d.Hashtags = []Hashtag{}
	for _, h := range data.Hashtags {
		hashTag := Hashtag{ID: h.HashtagID.ID, Value: h.HashtagID.Name}
		d.Hashtags = append(d.Hashtags, hashTag)
	}

	d.Images = []string{}
	for _, image := range data.Images {
		if image.Validate() {
			d.Images = append(d.Images, config.GetDirectusGetAssetURI(image.ID))
		}
	}

	d.Videos = []Video{}
	for _, video := range data.Videos {
		if v, err := video.NewVideo(); err == nil {
			d.Videos = append(d.Videos, v)
		}
	}

	//it is for single event(get event api)
	if account != nil {
		for i := range account.LikedEvents {
			if d.ID == account.LikedEvents[i].EventID {
				d.IsLiked = true
				break
			}
		}
	}
}

func (d *GetEventsResponse) parse(data []DirectusEventResponseData, start int64, limit int64, locale string, account *DirectusAccountResponseData, total int64) {

	d.Results = []GetEventResponse{}
	for _, event := range data {
		directusEvent := GetEventResponse{}
		directusEvent.parse(event, account)
		d.Results = append(d.Results, directusEvent)
	}

	if limit > 0 && start > 0 {
		d.Pages.Prev = fmt.Sprintf("/api/hubs-cms/v1/events?start=%d&limit=%d&locale=%s", utils.Max(0, start-limit), limit, locale)
	}

	if limit > 0 && total > start+limit {
		d.Pages.Next = fmt.Sprintf("/api/hubs-cms/v1/events?start=%d&limit=%d&locale=%s", start+limit, limit, locale)
	}
}

func NewEventResponse(data DirectusEventResponseData, accountData *DirectusAccountResponseData) *GetEventResponse {
	result := GetEventResponse{}
	result.parse(data, accountData)
	return &result
}

func NewDirectusEvents(data []DirectusEventResponseData, start int64, limit int64, locale string, total int64) *GetEventsResponse {
	result := GetEventsResponse{}
	result.parse(data, start, limit, locale, nil, total)
	return &result
}

type GetEventsRequestParam struct {
	Locale string      `form:"locale" binding:"omitempty,bcp47_language_tag"`
	Limit  json.Number `form:"limit" binding:"omitempty,PageLimitValidator"`
	Start  json.Number `form:"start" binding:"omitempty,PageStartValidator"`
	Status string      `form:"status" binding:"omitempty"`
}
type GetEventRequest struct {
	ID     string `uri:"id" binding:"required,uuid"`
	Locale string `form:"locale" binding:"omitempty,bcp47_language_tag"`
}
type GetEventIDParam struct {
	Locale string `form:"locale" binding:"omitempty,bcp47_language_tag"`
}
type EventIDRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}
type EventLikeCountResponse struct {
	LikeCount json.Number `json:"like_count"`
}
