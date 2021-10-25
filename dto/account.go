package dto

import "encoding/json"

type DirectusAccount struct {
	ID              string          `json:"id"`
	MastodonAccount string          `json:"mastodon_account"`
	MastodonAvatar  string          `json:"mastodon_avatar"`
	DisplayName     string          `json:"display_name"`
	IsAdmin         bool            `json:"is_admin"`
	ActiveAvatar    *DirectusAvatar `json:"active_avatar"`
	LikedRooms      []string        `json:"liked_rooms"`
	LikedEvents     []string        `json:"liked_events"`
}

type DirectusCreateAccountRequest struct {
	MastodonAccount string `json:"mastodon_account"`
	MastodonAvatar  string `json:"mastodon_avatar"`
	DisplayName     string `json:"display_name"`
	IsAdmin         bool   `json:"is_admin"`
}

type DirectusUpsertAccountResponse struct {
	Data DirectusAccountResponseData `json:"data"`
}

type DirectusAccountResponseData struct {
	ID              string                      `json:"id"`
	MastodonAccount string                      `json:"mastodon_account"`
	MastodonAvatar  string                      `json:"mastodon_avatar"`
	DisplayName     string                      `json:"display_name"`
	IsAdmin         bool                        `json:"is_admin"`
	ActiveAvatar    DirectusAvatarResponseData  `json:"active_avatar"`
	LikedRooms      []DirectusAccountLikedRoom  `json:"liked_rooms"`
	LikedEvents     []DirectusAccountLikedEvent `json:"liked_events"`
}

type DirectusAccountLikedRoom struct {
	ID     json.Number `json:"id"`
	RoomID string      `json:"room_id"`
}

type DirectusAccountLikedEvent struct {
	ID      json.Number `json:"id"`
	EventID string      `json:"event_id"`
}

func (r DirectusUpsertAccountResponse) Validate() bool {
	return len(r.Data.MastodonAccount) > 0
}

type DirectusGetAccountRequestParam struct {
	Email string `form:"mastodon_account" binding:"required,StringNotEmptyValidator"`
}

type DirectusGetAccountResponse struct {
	Data []DirectusAccountResponseData `json:"data"`
}

func (r DirectusGetAccountResponse) Validate() bool {
	return len(r.Data) == 1 && len(r.Data[0].MastodonAccount) > 0
}

type AccountIDRequestURI struct {
	AccountID string `uri:"accountId" binding:"required,uuid4"`
}

type PatchAccountRequestBody struct {
	DisplayName    string `json:"display_name,omitempty" binding:"omitempty,AccountDisplayNameValidator" example:"User X"`
	ActiveAvatarID string `json:"active_avatar,omitempty" binding:"omitempty,uuid4" example:"xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"`
	MastodonAvatar string `json:"mastodon_avatar,omitempty" binding:"omitempty,url" example:"https://test.host.com/assets/my-avatar-url"`
}

func (r PatchAccountRequestBody) Validate() bool {
	return !(len(r.ActiveAvatarID) == 0 && len(r.DisplayName) == 0)
}

type DirectusGetAccountLikesResponse struct {
	Meta DirectusMeta                          `json:"meta"`
	Data []DirectusGetAccountLikesResponseData `json:"data"`
}

type LikedEvents struct {
	EventID string `json:"event_id"`
}

type LikedRooms struct {
	RoomID string `json:"room_id"`
}

type DirectusGetAccountLikesResponseData struct {
	LikedEvents []LikedEvents `json:"liked_events"`
	LikedRooms  []LikedRooms  `json:"liked_rooms"`
}
