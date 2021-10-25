package dto

import "encoding/json"

type DirectusGetResponse struct {
	Data interface{}  `json:"data"`
	Meta DirectusMeta `json:"meta"`
}

type DirectusFile struct {
	ID string `json:"id"`
}

func (f *DirectusFile) Validate() bool {
	return f != nil && f.ID != ""
}

type DierctusRoomData struct {
	ID           string                `json:"id"`
	Title        string                `json:"title"`
	Description  string                `json:"description"`
	LikeCount    json.Number           `json:"like_count"`
	ViewCount    json.Number           `json:"view_count"`
	HasNFT       bool                  `json:"has_nft"`
	IsPublic     bool                  `json:"is_public"`
	Passcode     string                `json:"passcode"`
	Translations []DirectusRoomL10N    `json:"translations"`
	Gallery      DirectusFile          `json:"gallery"`
	Owner        string                `json:"owner"`
	HubsID       string                `json:"hubs_id"`
	JoinedEvents []DirectusJoinedEvent `json:"events"`
	NFTContract  *DirectusNFTContract  `json:"nft_contract"`
}

type DirectusRoomL10N struct {
	ID          json.Number `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
}

type DirectusJoinedEvent struct {
	EventID string `json:"event_id"`
}

type DirectusNFTContract struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Source     string `json:"source"`
	Platform   string `json:"platform"`
	Standard   string `json:"standard"`
	Blockchain string `json:"blockchain"`
	Address    string `json:"address"`
}

func (data *DierctusRoomData) UpdateTranslation() {
	if data == nil {
		return
	}

	if len(data.Translations) > 0 {
		target := &data.Translations[0]
		UpdateTranslationIfNeed(&data.Title, &target.Title)
		UpdateTranslationIfNeed(&data.Description, &target.Description)
	}
}

//=============================
type omit *struct{}

type GetRoomListResponse struct {
	Results []GetRoomResponseWrap `json:"results"`
	Pages   Page                  `json:"pages"`
}

type GetRoomResponseWrap struct {
	*GetRoomResponse
	NFT    omit `json:"nft,omitempty"`
	Events omit `json:"events,omitempty"`
}

type GetRoomResponse struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	ViewCount   json.Number      `json:"view_count"`
	LikeCount   json.Number      `json:"like_count"`
	ImageURL    string           `json:"image_url"`
	HasNFT      bool             `json:"has_nft"`
	IsLiked     bool             `json:"is_liked"`
	IsPublic    bool             `json:"is_public"`
	IsProtected bool             `json:"is_protected"`
	Owner       string           `json:"owner"`
	HubsURL     string           `json:"hubs_url"`
	NFT         *RoomNFTResponse `json:"nft"`
	Events      []string         `json:"events"`
}

type RoomNFTResponse struct {
	Standard   string `json:"standard"`
	Blockchain string `json:"blockchain"`
	Address    string `json:"address"`
}

type RoomLikeCountResponse struct {
	LikeCount json.Number `json:"like_count"`
}

//=============================
type GetRoomListRequest struct {
	Limit  json.Number `form:"limit" binding:"omitempty,PageLimitValidator"`
	Start  json.Number `form:"start" binding:"omitempty,PageStartValidator"`
	Locale string      `form:"locale" binding:"omitempty,bcp47_language_tag"`
	HubsID string      `form:"hubs_id" binding:"omitempty"`
	HasNFT bool        `form:"has_nft" binding:"omitempty"`
}

type GetRoomRequest struct {
	ID     string `uri:"id" binding:"required,uuid"`
	Locale string `form:"locale" binding:"omitempty,bcp47_language_tag"`
}

type RoomDataIncreaseViewCountRequest struct {
	ViewCount int64 `json:"view_count"`
}

type RoomIDRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type HubsIDRequest struct {
	HubsID string `uri:"hubsid" binding:"required,alphanum"`
}

type HubsPasscodeRequest struct {
	Passcode string `json:"passcode" binding:"required"`
}
