package handler

import (
	"hubs-cms-go/dto"
	"hubs-cms-go/errors"
	"hubs-cms-go/service"
	"hubs-cms-go/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// @Summary Get all avatars
// @Description Get all avatars
// @Tags avatars
// @Accept  json
// @Produce json
// @param start path int true "0" Format(int64)
// @param limit path int true "10" Format(int64)
// @Success 200 {object} dto.GetAvatarsResponse
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/avatars [get]
func GetPublicAvatars(c *gin.Context) {

	pageRequestParam := dto.PageRequestParam{}
	if err := c.ShouldBindWith(&pageRequestParam, binding.Form); err != nil {
		if validators.IsInvalid("PageRequestParam.Limit", err) {
			c.JSON(http.StatusBadRequest, errors.AvatarsInvalidLimit)
			return
		}
		if validators.IsInvalid("PageRequestParam.Start", err) {
			c.JSON(http.StatusBadRequest, errors.AvatarsInvalidStart)
			return
		}
		c.JSON(http.StatusBadRequest, errors.AvatarsInvalidRequestFormat)
		return
	}

	start, _ := pageRequestParam.Start.Int64()
	limit, _ := pageRequestParam.Limit.Int64()

	directusAvatars, errorInfo := service.GetPublicAvatars(start, limit, false)
	if errorInfo != (errors.ErrorInfo{}) {
		c.JSON(http.StatusInternalServerError, errorInfo)
		return
	}

	c.JSON(http.StatusOK, directusAvatars)
}

// @Summary Get all avatars of login user
// @Description Get all avatars of login user
// @Tags avatars
// @Accept  json
// @Produce json
// @param start path int true "0" Format(int64)
// @param limit path int true "10" Format(int64)
// @Success 200 {object} dto.GetAvatarsResponse
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/my-avatars [get]
func GetMyAvatars(c *gin.Context) {

	mastodonAccountInfo, err := GetMastodonAccountInfo(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	pageRequestParam := dto.PageRequestParam{}
	if err := c.ShouldBindWith(&pageRequestParam, binding.Form); err != nil {
		if validators.IsInvalid("PageRequestParam.Limit", err) {
			c.JSON(http.StatusBadRequest, errors.AvatarsInvalidLimit)
			return
		}
		if validators.IsInvalid("PageRequestParam.Start", err) {
			c.JSON(http.StatusBadRequest, errors.AvatarsInvalidStart)
			return
		}
		c.JSON(http.StatusBadRequest, errors.AvatarsInvalidRequestFormat)
		return
	}

	start, _ := pageRequestParam.Start.Int64()
	limit, _ := pageRequestParam.Limit.Int64()

	directusAccount, err := service.GetDirectusAccount(mastodonAccountInfo.MastodonAccount, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	if len(directusAccount.ID) == 0 {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	directusAvatars, errorInfo := service.GetMyAvatars(directusAccount.ID, start, limit, false)
	if errorInfo != (errors.ErrorInfo{}) {
		c.JSON(http.StatusInternalServerError, errorInfo)
		return
	}

	c.JSON(http.StatusOK, directusAvatars)
}

// @Summary Create an avatar after glb and snapshot uploaded
// @Description Create an avatar after glb and snapshot uploaded
// @Tags avatars
// @Accept  multipart/form-data
// @Produce json
// @Param id path string true "Account ID"
// @Param body body dto.UploadAvatarRequest.ID true "Create avatar"
// @Param file formData file true "snapshot"
// @Success 200 {object} dto.DirectusAvatar
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/avatars/{id} [post]
func CreateAvatar(c *gin.Context) {

	mastodonAccountInfo, err := GetMastodonAccountInfo(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	directusAccount, err := service.GetDirectusAccount(mastodonAccountInfo.MastodonAccount, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	uploadAvatarRequest := dto.UploadAvatarRequest{}
	if err := c.Bind(&uploadAvatarRequest); err != nil {
		c.JSON(http.StatusBadRequest, errors.AvatarsInvalidRequestFormat)
		return
	}

	snapshotID, err := service.UploadAsset(uploadAvatarRequest.Snapshot, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	glbID, err := service.ImportAsset(uploadAvatarRequest.GLB, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	directusCreateAvatarRequest := dto.DirectusCreateAvatarRequest{
		Title:    uploadAvatarRequest.Title,
		Source:   uploadAvatarRequest.Source,
		IsPublic: *uploadAvatarRequest.IsPublic,
		GLB:      glbID,
		Snapshot: snapshotID,
		Owner:    directusAccount.ID,
	}

	createdAvatar, err := service.CreateAvatar(directusCreateAvatarRequest, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	c.JSON(http.StatusOK, createdAvatar)
}

// @Summary delete an avatar by id
// @Description delete an avatar by id
// @Tags avatars
// @Param id path string true "Account ID"
// @Success 200 {string} ok
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/avatars/{id} [delete]
func DeleteAvatar(c *gin.Context) {

	mastodonAccountInfo, err := GetMastodonAccountInfo(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	directusAccount, err := service.GetDirectusAccount(mastodonAccountInfo.MastodonAccount, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	deleteAvatarRequest := dto.DirectusDeleteAvatarRequest{}
	if err := c.ShouldBindUri(&deleteAvatarRequest); err != nil {
		c.JSON(http.StatusBadRequest, errors.AvatarsInvalidID)
		return
	}

	avatarResponse, err := service.GetAvatar(deleteAvatarRequest.ID, false)
	if err != nil {
		c.JSON(http.StatusForbidden, errors.ForbiddenError)
		return
	}

	if avatarResponse.Data.Owner != directusAccount.ID || avatarResponse.Data.IsPublic {
		c.JSON(http.StatusForbidden, errors.ForbiddenError)
		return
	}

	if err := service.DeleteAvatar(deleteAvatarRequest.ID, false); err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	c.Status(http.StatusOK)
}
