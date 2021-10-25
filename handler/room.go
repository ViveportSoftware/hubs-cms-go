package handler

import (
	"encoding/json"
	"fmt"
	"hubs-cms-go/cache"
	"hubs-cms-go/config"
	"hubs-cms-go/dto"
	"hubs-cms-go/errors"
	"hubs-cms-go/logger"
	"hubs-cms-go/service"
	"hubs-cms-go/validators"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	goCache "github.com/patrickmn/go-cache"
)

// @Summary check passcode by hubs ID
// @Description check passcode by hubs ID
// @Tags rooms
// @Accept  json
// @Produce json
// @Param hubsid path string true "Hubs ID"
// @Success 200 {string} string "ok"
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/passcode/{hubsid} [post]
func CheckHubsPasscode(c *gin.Context) {
	param := struct {
		dto.HubsIDRequest
		dto.HubsPasscodeRequest
	}{}
	if err := c.ShouldBindUri(&param.HubsIDRequest); err != nil {
		c.JSON(http.StatusBadRequest, errors.RoomInvalidHubsID)
		return
	}
	if err := c.ShouldBindJSON(&param.HubsPasscodeRequest); err != nil {
		c.JSON(http.StatusBadRequest, errors.RoomInvalidPasscode)
		return
	}

	hubsID := param.HubsID
	directusRoomList, total, err := service.GetDirectusRoomList(nil, hubsID, "", 0, 0)
	if err != nil {
		if dsErr, ok := err.(*dto.DirectusErrorResponse); ok {
			ee := directusErrorHandler(dsErr)
			c.JSON(ee.HttpStatus, ee)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalError)
		}
		return
	}

	if total > 0 {
		if total > 1 {
			rooms := make([]string, total)
			for i := int64(0); i < total; i++ {
				rooms[i] = directusRoomList[i].ID
			}
			logger.Error.Printf("[CheckHubsPasscode] Same HubsID(%v) for %v rooms: [%v]\n", hubsID, total, strings.Join(rooms, ","))
			c.JSON(http.StatusInternalServerError, errors.InternalError)
			return
		}

		if passcode := directusRoomList[0].Passcode; len(passcode) > 0 && passcode != param.Passcode {
			c.JSON(http.StatusForbidden, errors.ForbiddenError)
			return
		}
	}

	c.Status(http.StatusOK)
}

// @Summary To mark room as like
// @Description To mark the room as like, and to increase the like count by one
// @Tags rooms
// @Accept  json
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {string} string "ok"
// @Failure 400 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/rooms/{id}/liked [post]
func PostLikeRoom(c *gin.Context) {
	prepareToggleRoomLike(c, true)
}

// @Summary To mark the room as unlike
// @Description To mark the room as unlike, and to reduce the like count by one
// @Tags rooms
// @Accept  json
// @Produce json
// @Param id path string true "Room ID"
// @Success 200 {string} string "ok"
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/rooms/{id}/unliked [post]
func PostUnlikeRoom(c *gin.Context) {
	prepareToggleRoomLike(c, false)
}

func prepareToggleRoomLike(c *gin.Context, isDoLike bool) {
	param := dto.RoomIDRequest{}
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, errors.RoomInvalidID)
		return
	}

	pDirectusAccount := getDirectusAccountDataByHeaderInfo(c)
	if pDirectusAccount == nil {
		logger.Warn.Println("[prepareToggleRoomLike] Cannot find account")
		c.JSON(http.StatusUnauthorized, errors.UnauthorizedError)
		return
	}

	directusRoom, err := service.GetDirectusRoom(param.ID, "")
	if err != nil {
		if dsErr, ok := err.(*dto.DirectusErrorResponse); ok {
			ee := directusErrorHandler(dsErr)
			c.JSON(ee.HttpStatus, ee)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalError)
		}
		return
	}

	if !directusRoom.IsPublic {
		// if pDirectusAccount == nil {
		// 	c.JSON(http.StatusUnauthorized, errors.UnauthorizedError)
		// 	return
		// }

		logger.Debug.Println("[prepareToggleRoomLike] owner:", directusRoom.Owner, "account:", pDirectusAccount.ID)
		if directusRoom.Owner != pDirectusAccount.ID {
			c.JSON(http.StatusUnauthorized, errors.UnauthorizedError)
			return
		}
	}

	var likeCount int64
	likeCount, err = toggleRoomLike(pDirectusAccount, &directusRoom, isDoLike)
	if err != nil {
		if dsErr, ok := err.(*dto.DirectusErrorResponse); ok {
			ee := directusErrorHandler(dsErr)
			c.JSON(ee.HttpStatus, ee)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalError)
		}
		return
	}

	c.JSON(http.StatusOK, &dto.RoomLikeCountResponse{LikeCount: json.Number(fmt.Sprintf("%v", likeCount))})
}

func toggleRoomLike(pDirectusAccount *dto.DirectusAccountResponseData, pDirectusRoom *dto.DierctusRoomData, isDoLike bool) (likeCount int64, err error) {
	if pDirectusAccount == nil || pDirectusRoom == nil {
		return
	}

	indexOfRoom := -1
	for i := range pDirectusAccount.LikedRooms {
		if pDirectusRoom.ID == pDirectusAccount.LikedRooms[i].RoomID {
			indexOfRoom = i
			break
		}
	}

	// already liked / not like
	alreadyLiked := indexOfRoom >= 0

	likeCount = processRoomLikes(pDirectusRoom.ID, isDoLike, alreadyLiked)

	if isDoLike == alreadyLiked {
		logger.Debug.Printf("[toggleRoomLike] Ignore [%+v] action\n", isDoLike)
		return
	}

	m2mPatchBody := struct {
		LikedRooms struct {
			*dto.DirectusM2MPatchRequest
		} `json:"liked_rooms,omitempty"`
	}{}
	if isDoLike {

		createBody := struct {
			RoomID    string `json:"room_id"`
			AccountID string `json:"account_id"`
		}{
			RoomID:    pDirectusRoom.ID,
			AccountID: pDirectusAccount.ID,
		}
		m2mPatchBody.LikedRooms.DirectusM2MPatchRequest = &dto.DirectusM2MPatchRequest{
			Create: []interface{}{createBody},
		}
	} else {

		var recordID int64
		if recordID, err = pDirectusAccount.LikedRooms[indexOfRoom].ID.Int64(); err != nil {
			logger.Warn.Printf("[toggleRoomLike] Parse id error: %v\n", err)
			err = nil // ignore the err
		} else {
			m2mPatchBody.LikedRooms.DirectusM2MPatchRequest = &dto.DirectusM2MPatchRequest{
				Delete: []int64{recordID},
			}
		}
		logger.Debug.Printf("[toggleRoomLike] like count: %v, index: %v, id: %v\n", likeCount, indexOfRoom, recordID)
	}

	if _, err = service.PatchDirectusAccount(pDirectusAccount.ID, &m2mPatchBody, false); err != nil {
		return
	}

	return
}

func processRoomLikes(id string, isDoLike, alreadyLiked bool) (likes int64) {

	logger.Debug.Printf("[processRoomLikes] id=%v, isDoLike=%v, alreadyLiked=%v\n", id, isDoLike, alreadyLiked)
	if isDoLike {
		//like +1
		if alreadyLiked {
			_likes, found := cache.RoomLikes.Get(id)
			if !found {
				//id not found, it should be new event, so create it
				logger.Debug.Printf("[processRoomLikes] likes+1: id not found, so create it.")
				err := cache.RoomLikes.Add(id, int64(1), goCache.NoExpiration)
				if err != nil {
					logger.Error.Printf("[processRoomLikes] likes+1 [%+v]\n", err)
				}
			}
			likes = _likes.(int64)
		} else {
			_likes, err := cache.RoomLikes.IncrementInt64(id, 1)
			if err != nil {
				//new event, so create it
				logger.Debug.Printf("[processRoomLikes] likes+1 [%+v]\n", err)
				err = cache.RoomLikes.Add(id, int64(1), goCache.NoExpiration)
				if err != nil {
					logger.Debug.Printf("[processRoomLikes] likes+1 [%+v]\n", err)
					_likes, err = cache.RoomLikes.IncrementInt64(id, 1)
					if err != nil {
						logger.Error.Printf("[processRoomLikes] likes+1 [%+v]\n", err)
					}
				} else {
					_likes = 1
				}
			}
			likes = _likes
		}
	} else {
		//like -1
		if alreadyLiked {
			_likes, err := cache.RoomLikes.DecrementInt64(id, 1)
			if err != nil {
				logger.Debug.Printf("[processRoomLikes] likes-1 [%+v]\n", err)
			}
			likes = _likes
		} else {
			_likes, found := cache.RoomLikes.Get(id)
			if found {
				likes = _likes.(int64)
			}
		}
	}

	logger.Debug.Printf("[processRoomLikes] likes=%v", likes)
	return
}

// @Summary Get all rooms of login user
// @Description Get all rooms of login user
// @Tags rooms
// @Accept  json
// @Produce json
// @param locale query string false "en-US"
// @param start query int false "0" Format(int64)
// @param limit query int false "10" Format(int64)
// @Success 200 {object} dto.GetRoomListResponse
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/my-rooms [get]
func GetMyRooms(c *gin.Context) {
	param := dto.GetRoomListRequest{}
	if err := c.ShouldBindQuery(&param); err != nil {
		if validators.IsInvalid("GetRoomListRequest.Limit", err) {
			c.JSON(http.StatusBadRequest, errors.RoomInvalidLimit)
			return
		}
		if validators.IsInvalid("GetRoomListRequest.Start", err) {
			c.JSON(http.StatusBadRequest, errors.RoomInvalidStart)
			return
		}
		c.JSON(http.StatusBadRequest, errors.RoomInvalidRequestFormat)
		return
	}

	locale := param.Locale //c.DefaultQuery("locale", "en-US")
	var start, limit int64 = 0, 0
	if _, exist := c.GetQuery("start"); exist {
		start, _ = param.Start.Int64()
	}
	if _, exist := c.GetQuery("limit"); exist {
		limit, _ = param.Limit.Int64()
	}

	pDirectusAccount := getDirectusAccountDataByHeaderInfo(c)
	if pDirectusAccount == nil {
		logger.Warn.Println("[GetMyRooms] cannot find account")
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	directusRoomList, total, err := service.GetDirectusMyRoomList(pDirectusAccount.ID, locale, start, limit)
	if err != nil {
		if dsErr, ok := err.(*dto.DirectusErrorResponse); ok {
			ee := directusErrorHandler(dsErr)
			c.JSON(ee.HttpStatus, ee)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalError)
		}
		return
	}

	if total == 0 {
		c.JSON(http.StatusOK, dto.GetRoomListResponse{
			Results: []dto.GetRoomResponseWrap{}, // not return null
		})
		return
	}
	if start >= total { // filter count is normal, but data will be epmty
		c.JSON(http.StatusBadRequest, errors.RoomInvalidStart)
		return
	}

	results := make([]dto.GetRoomResponseWrap, len(directusRoomList))
	for i := range results {
		room := generateResponse(&directusRoomList[i], pDirectusAccount)
		results[i] = dto.GetRoomResponseWrap{
			GetRoomResponse: room,
		}
	}

	paging := generatePagingResponse(c.Request.RequestURI, start, limit, total)

	c.JSON(http.StatusOK, dto.GetRoomListResponse{
		Results: results,
		Pages:   *paging,
	})
}

// @Summary Retrieve room detail by ID
// @Description Retrieve given numbers of room detail data.
// @Tags rooms
// @Accept  json
// @Produce json
// @Param hubs_id query string false "Hubs ID"
// @param locale query string false "en-US"
// @param start query int false "0" Format(int64)
// @param limit query int false "10" Format(int64)
// @Success 200 {object} dto.GetRoomListResponse
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/rooms [get]
func GetRoomList(c *gin.Context) {
	param := dto.GetRoomListRequest{}
	if err := c.ShouldBindQuery(&param); err != nil {
		if validators.IsInvalid("GetRoomListRequest.Limit", err) {
			c.JSON(http.StatusBadRequest, errors.RoomInvalidLimit)
			return
		}
		if validators.IsInvalid("GetRoomListRequest.Start", err) {
			c.JSON(http.StatusBadRequest, errors.RoomInvalidStart)
			return
		}
		c.JSON(http.StatusBadRequest, errors.RoomInvalidRequestFormat)
		return
	}

	hubsID := param.HubsID
	locale := param.Locale //c.DefaultQuery("locale", "en-US")
	var start, limit int64 = 0, 0
	if _, exist := c.GetQuery("start"); exist {
		start, _ = param.Start.Int64()
	}
	if _, exist := c.GetQuery("limit"); exist {
		limit, _ = param.Limit.Int64()
	}
	var pHasNFT *bool
	if _, exist := c.GetQuery("has_nft"); exist {
		pHasNFT = &param.HasNFT
	}

	directusRoomList, total, err := service.GetDirectusRoomList(pHasNFT, hubsID, locale, start, limit)
	if err != nil {
		if dsErr, ok := err.(*dto.DirectusErrorResponse); ok {
			ee := directusErrorHandler(dsErr)
			c.JSON(ee.HttpStatus, ee)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalError)
		}
		return
	}

	if total == 0 {
		c.JSON(http.StatusOK, dto.GetRoomListResponse{
			Results: []dto.GetRoomResponseWrap{}, // not return null
		})
		return
	}
	if start >= total { // filter count is normal, but data will be epmty
		c.JSON(http.StatusBadRequest, errors.RoomInvalidStart)
		return
	}

	pDirectusAccount := getDirectusAccountDataByHeaderInfo(c)

	results := make([]dto.GetRoomResponseWrap, len(directusRoomList))
	for i := range results {
		room := generateResponse(&directusRoomList[i], pDirectusAccount)
		results[i] = dto.GetRoomResponseWrap{
			GetRoomResponse: room,
		}
	}

	paging := generatePagingResponse(c.Request.RequestURI, start, limit, total)

	c.JSON(http.StatusOK, dto.GetRoomListResponse{
		Results: results,
		Pages:   *paging,
	})
}

// @Summary Retrieve specific room detail
// @Description Retrieve specific room detail
// @Tags rooms
// @Accept  json
// @Produce json
// @Param id path string true "Room ID"
// @param locale path string false "en-US"
// @Success 200 {object} dto.GetRoomResponse
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/rooms/{id} [get]
func GetRoom(c *gin.Context) {
	param := dto.GetRoomRequest{}
	if err := c.ShouldBindUri(&param); err != nil {
		//c.Abort()
		c.JSON(http.StatusBadRequest, errors.RoomInvalidID)
		return
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		//c.Abort()
		c.JSON(http.StatusBadRequest, errors.RoomInvalidRequestFormat)
		return
	}

	roomId := param.ID
	locale := param.Locale

	directusRoom, err := service.GetDirectusRoom(roomId, locale)
	if err != nil {
		if dsErr, ok := err.(*dto.DirectusErrorResponse); ok {
			ee := directusErrorHandler(dsErr)
			c.JSON(ee.HttpStatus, ee)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalError)
		}
		return
	}

	pDirectusAccount := getDirectusAccountDataByHeaderInfo(c)

	if !directusRoom.IsPublic {
		if pDirectusAccount == nil {
			c.JSON(http.StatusUnauthorized, errors.UnauthorizedError)
			return
		}

		logger.Debug.Println("[GetRoom] owner: ", directusRoom.Owner, "account: ", pDirectusAccount.ID)
		if directusRoom.Owner != pDirectusAccount.ID {
			c.JSON(http.StatusUnauthorized, errors.UnauthorizedError)
			return
		}
	}

	c.JSON(http.StatusOK, generateResponse(&directusRoom, pDirectusAccount))
}

func generatePagingResponse(uri string, start, limit, total int64) *dto.Page {
	paging := dto.Page{}
	if limit > 0 && total > limit {
		if requestURI, err := url.Parse(uri); err == nil {
			q := requestURI.Query()
			if start > 0 {
				ss := start - limit
				if ss < 0 {
					ss = 0
				}
				q.Set("start", fmt.Sprintf("%v", ss))
				requestURI.RawQuery = q.Encode()
				paging.Prev = requestURI.String()
			}
			if start+limit < total {
				q.Set("start", fmt.Sprintf("%v", start+limit))
				requestURI.RawQuery = q.Encode()
				paging.Next = requestURI.String()
			}
		} else {
			logger.Error.Printf("[generatePagingResponse] Parse url error: %v\n", err)
		}
	}
	return &paging
}

func generateResponse(pDirectusRoom *dto.DierctusRoomData, pDirectusAccount *dto.DirectusAccountResponseData) *dto.GetRoomResponse {
	ret := dto.GetRoomResponse{
		ID:          pDirectusRoom.ID,
		Title:       pDirectusRoom.Title,
		Description: pDirectusRoom.Description,
		IsPublic:    pDirectusRoom.IsPublic,
		ViewCount:   pDirectusRoom.ViewCount,
		//LikeCount:   pDirectusRoom.LikeCount,
		Owner:       pDirectusRoom.Owner,
		IsProtected: len(pDirectusRoom.Passcode) > 0,
		HasNFT:      pDirectusRoom.HasNFT,
		//Is_liked: ,
	}

	if item, ok := cache.RoomLikes.Get(pDirectusRoom.ID); ok {
		if likes, ok := item.(int64); ok {
			ret.LikeCount = json.Number(fmt.Sprintf("%v", likes))
		}
	}

	if pDirectusRoom.Gallery.Validate() {
		ret.ImageURL = config.GetDirectusGetAssetURI(pDirectusRoom.Gallery.ID)
	}

	if len(pDirectusRoom.HubsID) > 0 {
		if hubsURL, err := config.GetHubsURL(pDirectusRoom.HubsID); err != nil {
			logger.Error.Printf("[generateResponse] Get hubs url error: %v\n", err)
		} else {
			ret.HubsURL = hubsURL
		}
	} else {
		logger.Error.Printf("[generateResponse] Missing hubs id for: %v\n", pDirectusRoom.ID)
	}

	if pDirectusRoom.NFTContract != nil {
		ret.NFT = &dto.RoomNFTResponse{
			Address:    pDirectusRoom.NFTContract.Address,
			Blockchain: pDirectusRoom.NFTContract.Blockchain,
			Standard:   pDirectusRoom.NFTContract.Standard,
		}
	}

	joinedEventCount := len(pDirectusRoom.JoinedEvents)
	if joinedEventCount > 0 {
		ret.Events = make([]string, joinedEventCount)
		for i := range pDirectusRoom.JoinedEvents {
			ret.Events[i] = string(pDirectusRoom.JoinedEvents[i].EventID)
		}
	}

	if pDirectusAccount != nil {
		for i := range pDirectusAccount.LikedRooms {
			if pDirectusRoom.ID == pDirectusAccount.LikedRooms[i].RoomID {
				ret.IsLiked = true
				break
			}
		}
	}

	return &ret
}

func directusErrorHandler(err *dto.DirectusErrorResponse) (ret *errors.ErrorInfo) {
	if err == nil {
		return
	}

	if len(err.Errors) > 1 {
		logger.Warn.Printf("[directusErrorHandler] Handle 1 of %v errors\n", len(err.Errors))
	}

	target := err.Errors[0]
	code := target.Extensions.Code
	switch code {
	case "FORBIDDEN":
		// not allowed to do the current action OR
		// non-existing items
		err := errors.ForbiddenError
		err.ErrorBody.Message = target.Message
		ret = &err
	//case "INVALID_CREDENTIALS":
	// Username / password or access token is wrong
	default:
		err := errors.InternalError
		err.ErrorBody.Message = code
		ret = &err
	}
	return
}

// @Summary Increase room viewed count whenever this API gets called.
// @Description increase overall view count by 1 for the given room
// @Tags rooms
// @Accept  json
// @Produce json
// @Param id path string true "Room ID"
// @param locale path string false "en-US"
// @Success 200 {object} dto.GetRoomResponse
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/rooms/{id}/viewed [post]
func RoomViewCountHandler(c *gin.Context) {
	param := dto.GetRoomRequest{}

	if err := c.ShouldBindUri(&param); err != nil {
		//c.Abort()
		c.JSON(http.StatusBadRequest, errors.RoomInvalidID)
		return
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		//c.Abort()
		c.JSON(http.StatusBadRequest, errors.RoomInvalidRequestFormat)
		return
	}

	roomId := param.ID
	locale := param.Locale

	directusRoom, err := service.GetDirectusRoom(roomId, locale)
	if err != nil {
		if dsErr, ok := err.(*dto.DirectusErrorResponse); ok {
			ee := directusErrorHandler(dsErr)
			c.JSON(ee.HttpStatus, ee)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalError)
		}
		return
	}

	pDirectusAccount := getDirectusAccountDataByHeaderInfo(c)

	if !directusRoom.IsPublic {
		if pDirectusAccount == nil {
			c.JSON(http.StatusUnauthorized, errors.UnauthorizedError)
			return
		}

		logger.Debug.Println("[RoomViewCountHandler] owner: ", directusRoom.Owner, "account: ", pDirectusAccount.ID)
		if directusRoom.Owner != pDirectusAccount.ID {
			c.JSON(http.StatusUnauthorized, errors.UnauthorizedError)
			return
		}
	}

	addViewCount, errInfo := service.PostRoomViewCount(directusRoom, locale, false)
	if errInfo != (errors.ErrorInfo{}) {
		c.JSON(errInfo.HttpStatus, errInfo)
		return
	}

	c.JSON(http.StatusOK, generateResponse(&addViewCount, pDirectusAccount))
}
