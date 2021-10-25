package dto

import "mime/multipart"

type DirectusGetAvatarResponse struct {
	Data DirectusAvatarResponseData `json:"data"`
}

func (r DirectusGetAvatarResponse) Validate() bool {
	return r.Data.Validate()
}

type DirectusAvatarResponseData struct {
	ID       string `json:"id"`
	Snapshot string `json:"snapshot"`
	GLB      string `json:"glb"`
	Owner    string `json:"owner"`
	Source   string `json:"source"`
	Title    string `json:"title"`
	IsPublic bool   `json:"is_public"`
}

func (r DirectusAvatarResponseData) Validate() bool {
	return len(r.ID) > 0 &&
		len(r.Snapshot) > 0 &&
		len(r.GLB) > 0
}

type DirectusGetAvatarsResponse struct {
	Meta DirectusMeta                 `json:"meta"`
	Data []DirectusAvatarResponseData `json:"data"`
}

func (r DirectusGetAvatarsResponse) Validate() bool {
	for _, d := range r.Data {
		if !d.Validate() {
			return false
		}
	}
	return true
}

type GetAvatarsResponse struct {
	Results []DirectusAvatar `json:"results"`
	Pages   Page             `json:"pages"`
}

type DirectusAvatar struct {
	ID       string `json:"id"`
	Snapshot string `json:"snapshot_url"`
	GLB      string `json:"glb_url"`
	Owner    string `json:"owner"`
	Source   string `json:"source"`
	Title    string `json:"title"`
	IsPublic bool   `json:"is_public"`
}

type UploadAvatarRequest struct {
	Title    string                `form:"title" binding:"omitempty"`
	GLB      string                `form:"glb" binding:"required,url"`
	Source   string                `form:"source" binding:"required"`
	IsPublic *bool                 `form:"is_public" binding:"required"`
	Snapshot *multipart.FileHeader `form:"snapshot" binding:"required"`
}

type DirectusUploadAssetResponse struct {
	Data DirectusUploadAssetResponseData `json:"data"`
}

type DirectusUploadAssetResponseData struct {
	ID string `json:"id"`
}

func (r DirectusUploadAssetResponse) Validate() bool {
	return len(r.Data.ID) > 0
}

type DirectusImportAssetRequest struct {
	URL string `json:"url"`
}

type DirectusCreateAvatarRequest struct {
	Snapshot string `json:"snapshot"`
	GLB      string `json:"glb"`
	Owner    string `json:"owner"`
	Source   string `json:"source"`
	Title    string `json:"title"`
	IsPublic bool   `json:"is_public"`
}

type DirectusDeleteAvatarRequest struct {
	ID string `uri:"id" binding:"required,uuid"`
}
