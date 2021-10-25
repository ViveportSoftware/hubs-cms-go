package handler

import (
	"hubs-cms-go/constant"
	"hubs-cms-go/dto"
	"hubs-cms-go/logger"
	"hubs-cms-go/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getDirectusAccountDataByHeaderInfo(c *gin.Context) (pDirectusAccountData *dto.DirectusAccountResponseData) {
	if mastodonStatus, exists := c.Get(constant.HeaderMastodonHandlerStatus); !exists || mastodonStatus != http.StatusOK {
		return
	}

	if mastodonAccountInfo, err := GetMastodonAccountInfo(c); err == nil {
		directusAccountData, err := service.GetDirectusAccountData(mastodonAccountInfo.MastodonAccount)
		if err == nil && len(directusAccountData.ID) > 0 {
			pDirectusAccountData = &directusAccountData
		}

		logger.Debug.Println("[getDirectusAccountDataByHeaderInfo] MastodonAccount: ", mastodonAccountInfo.MastodonAccount, " id: ", directusAccountData.ID, " err: ", err)
	}

	return
}
