package handler

import (
	"hubs-cms-go/dto"
	"hubs-cms-go/errors"
	"hubs-cms-go/service"
	"hubs-cms-go/validators"
	"net/http"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
)

// @Summary Get login user profile
// @Description Get login user profile
// @Tags accounts
// @Accept  json
// @Produce json
// @Success 200 {object} dto.DirectusAccount
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/me [get]
func GetProfileMe(c *gin.Context) {

	mastodonAccountInfo, err := GetMastodonAccountInfo(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	isFetchAgain := false

	for {
		directusAccount, err := service.GetDirectusAccount(mastodonAccountInfo.MastodonAccount, false)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.InternalError)
			return
		}

		if len(directusAccount.ID) > 0 {
			if directusAccount.DisplayName != mastodonAccountInfo.DisplayName ||
				directusAccount.MastodonAvatar != mastodonAccountInfo.MastodonAvatar {
				// Update directus account based on mastodon profile
				patchAccountRequestBody := dto.PatchAccountRequestBody{
					DisplayName:    mastodonAccountInfo.DisplayName,
					MastodonAvatar: mastodonAccountInfo.MastodonAvatar,
				}

				directusAccount, err = service.PatchDirectusAccount(directusAccount.ID, &patchAccountRequestBody, false)
				if err != nil {
					c.JSON(http.StatusInternalServerError, errors.InternalError)
					return
				}
			}

			c.JSON(http.StatusOK, directusAccount)
			return
		}

		if isFetchAgain {
			c.JSON(http.StatusInternalServerError, errors.InternalError)
			return
		}

		// Create a new directus account based on mastodon account
		directusAccount, err = service.CreateDirectusAccount(mastodonAccountInfo, false)
		if err != nil {
			if dsErr, ok := err.(*dto.DirectusErrorResponse); ok {
				if dsErr.Status == http.StatusBadRequest {
					isFetchAgain = true
					continue
				}
			} else {
				c.JSON(http.StatusInternalServerError, errors.InternalError)
				return
			}
		}

		c.JSON(http.StatusOK, directusAccount)
		break
	}
}

// @Summary Update user profile
// @Description Update user profile
// @Tags accounts
// @Accept  json
// @Produce json
// @Param id path string true "Account ID"
// @Param body body dto.PatchAccountRequestBody true "Update user profile"
// @Success 200 {object} dto.DirectusAccount
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/accounts/{id} [patch]
func PatchAccount(c *gin.Context) {

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

	accountIDRequestURI := dto.AccountIDRequestURI{}
	if err := c.ShouldBindUri(&accountIDRequestURI); err != nil {
		c.JSON(http.StatusBadRequest, errors.AccountsInvalidAccountID)
		return
	}

	if directusAccount.ID != accountIDRequestURI.AccountID {
		c.JSON(http.StatusForbidden, errors.ForbiddenError)
		return
	}

	patchAccountRequestBody := dto.PatchAccountRequestBody{}
	if err := c.ShouldBindBodyWith(&patchAccountRequestBody, binding.JSON); err != nil {
		if validators.IsInvalid("ActiveAvatarID", err) {
			c.JSON(http.StatusBadRequest, errors.AccountsInvalidActiveAvatarID)
			return
		}
		if validators.IsInvalid("DisplayName", err) {
			c.JSON(http.StatusBadRequest, errors.AccountsInvalidDisplayName)
			return
		}
		c.JSON(http.StatusBadRequest, errors.AccountsInvalidRequestFormat)
		return
	}
	if !patchAccountRequestBody.Validate() {
		c.JSON(http.StatusBadRequest, errors.AccountsInvalidRequestFormat)
		return
	}

	if len(patchAccountRequestBody.ActiveAvatarID) > 0 {
		if _, err := service.GetDirectusAvatar(patchAccountRequestBody.ActiveAvatarID, false); err != nil {
			c.JSON(http.StatusBadRequest, errors.AccountsInvalidActiveAvatarID)
			return
		}
	}

	if len(patchAccountRequestBody.DisplayName) > 0 {
		if _, err := service.PatchMastodonAccount(mastodonAccountInfo.MastodonToken, dto.MastodonPatchAccountRequestBody{DisplayName: patchAccountRequestBody.DisplayName}); err != nil {
			c.JSON(http.StatusInternalServerError, errors.InternalError)
			return
		}
	}

	directusAccount, err = service.PatchDirectusAccount(directusAccount.ID, &patchAccountRequestBody, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	c.JSON(http.StatusOK, directusAccount)
}
