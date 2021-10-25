package handler

import (
	"fmt"
	"hubs-cms-go/config"
	"hubs-cms-go/constant"
	"hubs-cms-go/dto"
	"hubs-cms-go/errors"
	"hubs-cms-go/logger"
	"hubs-cms-go/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func ErrorMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.Last()
		if err == nil {
			return
		}
		logger.Error.Printf("[ErrorMiddleware] error: %v\n", err)
	}
}

func MastodonTokenHandler(c *gin.Context) {
	// get token from header Authorization: Bearer
	bearerTokenMiddlewareRequest := dto.BearerTokenMiddlewareRequest{}
	if err := c.ShouldBindHeader(&bearerTokenMiddlewareRequest); err != nil {
		c.Set(constant.HeaderMastodonHandlerStatus, http.StatusUnauthorized)
		return
	}

	// check token by /api/v1/accounts/verify_credentials
	verifyCredentialsResponse, err := service.GetMastodonVerifyCredentials(bearerTokenMiddlewareRequest.Token)
	if err != nil {
		c.Set(constant.HeaderMastodonHandlerStatus, http.StatusForbidden)
		return
	}

	if !strings.Contains(verifyCredentialsResponse.MastodonAccount, "@") {
		verifyCredentialsResponse.MastodonAccount = fmt.Sprintf("%s%s", verifyCredentialsResponse.MastodonAccount, config.DefaultMastodonAccountDomain)
	}

	c.Set(constant.HeaderMastodonHandlerStatus, http.StatusOK)

	// insert X-Mastodon-* into gin context for further processing
	c.Set(constant.HeaderMastodonID, verifyCredentialsResponse.ID)
	c.Set(constant.HeaderMastodonAccount, verifyCredentialsResponse.MastodonAccount)
	c.Set(constant.HeaderMastodonUsername, verifyCredentialsResponse.UserName)
	c.Set(constant.HeaderMastodonDisplayName, verifyCredentialsResponse.DisplayName)
	c.Set(constant.HeaderMastodonAvatar, verifyCredentialsResponse.MastodonAvatar)
	c.Set(constant.HeaderMastodonToken, bearerTokenMiddlewareRequest.Token)
}

func MastodonTokenStatusHandler(c *gin.Context) {
	if mastodonStatus, exists := c.Get(constant.HeaderMastodonHandlerStatus); !exists {
		c.JSON(http.StatusBadRequest, errors.BadRequestError(http.StatusBadRequest, ""))
	} else {
		switch mastodonStatus {
		case http.StatusOK:
			return
		case http.StatusForbidden:
			c.JSON(http.StatusForbidden, errors.ForbiddenError)
		case http.StatusUnauthorized:
			c.JSON(http.StatusUnauthorized, errors.UnauthorizedError)
		default:
			c.JSON(http.StatusInternalServerError, errors.InternalError)
		}
	}
	c.Abort()
}

func GetMastodonAccountInfo(c *gin.Context) (dto.MastodonVerifyCredentialsResponse, error) {
	mastodonAccountInfo := dto.MastodonVerifyCredentialsResponse{}
	if id, exists := c.Get(constant.HeaderMastodonID); exists {
		mastodonAccountInfo.ID = fmt.Sprintf("%v", id)
	}
	if mastodonAccount, exists := c.Get(constant.HeaderMastodonAccount); exists {
		mastodonAccountInfo.MastodonAccount = fmt.Sprintf("%v", mastodonAccount)
	}
	if userName, exists := c.Get(constant.HeaderMastodonUsername); exists {
		mastodonAccountInfo.UserName = fmt.Sprintf("%v", userName)
	}
	if displayName, exists := c.Get(constant.HeaderMastodonDisplayName); exists {
		mastodonAccountInfo.DisplayName = fmt.Sprintf("%v", displayName)
	}
	if mastodonAvatar, exists := c.Get(constant.HeaderMastodonAvatar); exists {
		mastodonAccountInfo.MastodonAvatar = fmt.Sprintf("%v", mastodonAvatar)
	}
	if mastodonToken, exists := c.Get(constant.HeaderMastodonToken); exists {
		mastodonAccountInfo.MastodonToken = fmt.Sprintf("%v", mastodonToken)
	}
	if !mastodonAccountInfo.Validate() {
		return dto.MastodonVerifyCredentialsResponse{}, fmt.Errorf("[GetMastodonAccountInfo] cannot find valid dto.MastodonVerifyCredentialsResponse from gin.Context")
	}
	return mastodonAccountInfo, nil
}
