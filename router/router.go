package router

import (
	"hubs-cms-go/config"
	_ "hubs-cms-go/docs"
	"hubs-cms-go/handler"
	"hubs-cms-go/validators"

	"github.com/gin-gonic/gin/binding"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Router represents the router of gin, use Router.run(:${PORT}) to start service
var Router *gin.Engine

// SetupRouter setup gin router
func SetupRouter() *gin.Engine {

	if config.IsDevEnv() {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// setup middlewares
	router.Use(gin.LoggerWithWriter(gin.DefaultWriter, "/health", "/version"))
	router.Use(gin.Recovery())
	router.Use(handler.ErrorMiddleware())
	router.Use(handler.CORSMiddleware())

	// setup validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("StringNotEmptyValidator", validators.StringNotEmptyValidator)
		_ = v.RegisterValidation("BearerTokenValidator", validators.BearerTokenValidator)
		_ = v.RegisterValidation("AccountDisplayNameValidator", validators.AccountDisplayNameValidator)
		_ = v.RegisterValidation("PageStartValidator", validators.PageStartValidator)
		_ = v.RegisterValidation("PageLimitValidator", validators.PageLimitValidator)
	}

	// checker api
	router.GET("/version", handler.VersionHandler)
	router.GET("/health", handler.HealthHandler)

	// account api
	router.GET("/api/hubs-cms/v1/me", handler.MastodonTokenHandler, handler.MastodonTokenStatusHandler, handler.GetProfileMe)
	router.PATCH("/api/hubs-cms/v1/accounts/:accountId", handler.MastodonTokenHandler, handler.MastodonTokenStatusHandler, handler.PatchAccount)

	// avatar api
	router.GET("/api/hubs-cms/v1/avatars", handler.GetPublicAvatars)
	router.GET("/api/hubs-cms/v1/my-avatars", handler.MastodonTokenHandler, handler.MastodonTokenStatusHandler, handler.GetMyAvatars)
	router.POST("/api/hubs-cms/v1/avatars", handler.MastodonTokenHandler, handler.MastodonTokenStatusHandler, handler.CreateAvatar)
	router.DELETE("/api/hubs-cms/v1/avatars/:id", handler.MastodonTokenHandler, handler.MastodonTokenStatusHandler, handler.DeleteAvatar)

	// room api
	router.GET("/api/hubs-cms/v1/rooms/:id", handler.MastodonTokenHandler, handler.GetRoom)
	router.GET("/api/hubs-cms/v1/rooms", handler.MastodonTokenHandler, handler.GetRoomList)
	router.GET("/api/hubs-cms/v1/my-rooms", handler.MastodonTokenHandler, handler.MastodonTokenStatusHandler, handler.GetMyRooms)
	router.POST("/api/hubs-cms/v1/rooms/:id/viewed", handler.MastodonTokenHandler, handler.RoomViewCountHandler)
	router.POST("/api/hubs-cms/v1/rooms/:id/liked", handler.MastodonTokenHandler, handler.MastodonTokenStatusHandler, handler.PostLikeRoom)
	router.POST("/api/hubs-cms/v1/rooms/:id/unliked", handler.MastodonTokenHandler, handler.MastodonTokenStatusHandler, handler.PostUnlikeRoom)
	router.POST("/api/hubs-cms/v1/passcode/:hubsid", handler.CheckHubsPasscode)

	// event api
	router.GET("/api/hubs-cms/v1/events", handler.GetEvents)
	router.GET("/api/hubs-cms/v1/events/:id", handler.MastodonTokenHandler, handler.GetEvent)
	router.POST("/api/hubs-cms/v1/events/:id/liked", handler.MastodonTokenHandler, handler.MastodonTokenStatusHandler, handler.PostLikeEvent)
	router.POST("/api/hubs-cms/v1/events/:id/unliked", handler.MastodonTokenHandler, handler.MastodonTokenStatusHandler, handler.PostUnlikeEvent)
	router.POST("/api/hubs-cms/v1/events/:id/viewed", handler.MastodonTokenHandler, handler.EventViewCountHandler)
	if mode := gin.Mode(); mode == gin.DebugMode {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router
}

// Setup setup gin router
func Setup() {
	Router = SetupRouter()
}
