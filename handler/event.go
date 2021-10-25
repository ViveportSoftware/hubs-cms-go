package handler

import (
	"encoding/json"
	"fmt"
	"hubs-cms-go/cache"
	"hubs-cms-go/dto"
	"hubs-cms-go/errors"
	"hubs-cms-go/logger"
	"hubs-cms-go/service"
	"hubs-cms-go/validators"
	"net/http"

	"github.com/gin-gonic/gin"
	goCache "github.com/patrickmn/go-cache"
)

// @Summary Display all events we currently have.
// @Description Retrieve given numbers of event detail data.
// @Tags events
// @Accept  json
// @Produce json
// @param start path int false "0" Format(int64)
// @param limit path int false "10" Format(int64)
// @param locale path string false "en-US"
// @Success 200 {object} dto.GetEventsResponse
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/events [get]
func GetEvents(c *gin.Context) {

	param := dto.GetEventsRequestParam{}
	if err := c.ShouldBindQuery(&param); err != nil {
		fmt.Printf("err %v\n", err)
		if validators.IsInvalid("GetEventsRequestParam.Limit", err) {
			c.JSON(http.StatusBadRequest, errors.EventInvalidLimit)
			return
		}
		if validators.IsInvalid("GetEventsRequestParam.Start", err) {
			c.JSON(http.StatusBadRequest, errors.EventInvalidStart)
			return
		}
		c.JSON(http.StatusBadRequest, errors.EventInvalidRequestFormat)
		return
	}

	locale := param.Locale
	start, _ := param.Start.Int64()
	limit, _ := param.Limit.Int64()
	status := param.Status

	directusEvents, total, err := service.GetDirectusEvents(locale, status, start, limit)
	if err != nil {
		if dsErr, ok := err.(*dto.DirectusErrorResponse); ok {
			ee := directusErrorHandler(dsErr)
			c.JSON(ee.HttpStatus, ee)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalError)
		}
		return
	}
	if start >= total && total > 0 { // filter count is normal, but data will be empty
		c.JSON(http.StatusBadRequest, errors.EventInvalidStart)
		return
	}

	res := dto.NewDirectusEvents(directusEvents, start, limit, locale, total)

	c.JSON(http.StatusOK, res)
}

// @Summary Retrieve event detail by ID
// @Description Retrieve detail data for specific event.
// @Tags events
// @Accept  json
// @Produce json
// @Param id path string true "Event ID"
// @param locale path string false "en-US"
// @Success 200 {object} dto.GetEventResponse
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/events/{id} [get]
func GetEvent(c *gin.Context) {
	param := dto.GetEventRequest{}
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, errors.EventInvalidID)
		return
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, errors.EventInvalidRequestFormat)
		return
	}

	directusEvent, err := service.GetDirectusEvent(param.ID, param.Locale)
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

	res := dto.NewEventResponse(directusEvent, pDirectusAccount)

	c.JSON(http.StatusOK, res)
}

// @Summary To mark event as like
// @Description To mark the event as like, and to increase the like count by one
// @Tags events
// @Accept  json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {string} string "ok"
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/events/{id}/liked [post]
func PostLikeEvent(c *gin.Context) {
	prepareToggleEventLike(c, true)
}

// @Summary To mark the event as unlike
// @Description To mark the event as unlike, and to reduce the like count by one
// @Tags events
// @Accept  json
// @Produce json
// @Param id path string true "Event ID"
// @Success 200 {string} string "ok"
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/events/{id}/unliked [post]
func PostUnlikeEvent(c *gin.Context) {
	prepareToggleEventLike(c, false)
}

func prepareToggleEventLike(c *gin.Context, isDoLike bool) {

	param := dto.EventIDRequest{}
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, errors.EventInvalidID)
		return
	}

	pDirectusAccount := getDirectusAccountDataByHeaderInfo(c)
	if pDirectusAccount == nil {
		logger.Warn.Println("[prepareToggleEventLike] Cannot find account")
		c.JSON(http.StatusUnauthorized, errors.UnauthorizedError)
		return
	}

	directusEvent, err := service.GetDirectusEvent(param.ID, "")
	if err != nil {
		if dsErr, ok := err.(*dto.DirectusErrorResponse); ok {
			ee := directusErrorHandler(dsErr)
			c.JSON(ee.HttpStatus, ee)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalError)
		}
		return
	}

	likeCount, err := toggleEventLike(pDirectusAccount, &directusEvent, isDoLike)
	if err != nil {
		if dsErr, ok := err.(*dto.DirectusErrorResponse); ok {
			ee := directusErrorHandler(dsErr)
			c.JSON(ee.HttpStatus, ee)
		} else {
			c.JSON(http.StatusInternalServerError, errors.InternalError)
		}
		return
	}

	c.JSON(http.StatusOK, &dto.EventLikeCountResponse{LikeCount: json.Number(fmt.Sprintf("%v", likeCount))})
}

func toggleEventLike(pDirectusAccount *dto.DirectusAccountResponseData, pDirectusEvent *dto.DirectusEventResponseData, isDoLike bool) (likeCount int64, err error) {
	if pDirectusAccount == nil || pDirectusEvent == nil {
		return
	}

	indexOfEvent := -1
	for i := range pDirectusAccount.LikedEvents {
		if pDirectusEvent.ID == pDirectusAccount.LikedEvents[i].EventID {
			indexOfEvent = i
			break
		}
	}

	// already liked / not like
	alreadyLiked := indexOfEvent >= 0

	likeCount = processEventLikes(pDirectusEvent.ID, isDoLike, alreadyLiked)

	if isDoLike == alreadyLiked {
		logger.Debug.Printf("[toggleEventLike] Ignore [%+v] action\n", isDoLike)
		return
	}

	m2mPatchBody := struct {
		LikedEvents struct {
			*dto.DirectusM2MPatchRequest
		} `json:"liked_events,omitempty"`
	}{}
	if isDoLike {

		createBody := struct {
			EventID   string `json:"event_id"`
			AccountID string `json:"account_id"`
		}{
			EventID:   pDirectusEvent.ID,
			AccountID: pDirectusAccount.ID,
		}
		m2mPatchBody.LikedEvents.DirectusM2MPatchRequest = &dto.DirectusM2MPatchRequest{
			Create: []interface{}{createBody},
		}
	} else {

		var recordID int64
		if recordID, err = pDirectusAccount.LikedEvents[indexOfEvent].ID.Int64(); err != nil {
			logger.Warn.Printf("[toggleEventLike] Parse id error: %v\n", err)
			err = nil // ignore the err
		} else {
			m2mPatchBody.LikedEvents.DirectusM2MPatchRequest = &dto.DirectusM2MPatchRequest{
				Delete: []int64{recordID},
			}
		}
		logger.Debug.Printf("[toggleEventLike] like count: %v, index: %v, id: %v\n", likeCount, indexOfEvent, recordID)
	}

	if _, err = service.PatchDirectusAccount(pDirectusAccount.ID, &m2mPatchBody, false); err != nil {
		return
	}

	return
}

func processEventLikes(id string, isDoLike, alreadyLiked bool) (likes int64) {

	logger.Debug.Printf("[processEventLikes] id=%v, isDoLike=%v, alreadyLiked=%v\n", id, isDoLike, alreadyLiked)
	if isDoLike {
		//like +1
		if alreadyLiked {
			_likes, found := cache.EventLikes.Get(id)
			if !found {
				//id not found, it should be new event, so create it
				logger.Debug.Printf("[processEventLikes] likes+1: id not found, so create it.")
				err := cache.EventLikes.Add(id, int64(1), goCache.NoExpiration)
				if err != nil {
					logger.Error.Printf("[processEventLikes] likes+1 [%+v]\n", err)
				}
			}
			likes = _likes.(int64)
		} else {
			_likes, err := cache.EventLikes.IncrementInt64(id, 1)
			if err != nil {
				//new event, so create it
				logger.Debug.Printf("[processEventLikes] likes+1 [%+v]\n", err)
				err = cache.EventLikes.Add(id, int64(1), goCache.NoExpiration)
				if err != nil {
					logger.Debug.Printf("[processEventLikes] likes+1 [%+v]\n", err)
					_likes, err = cache.EventLikes.IncrementInt64(id, 1)
					if err != nil {
						logger.Error.Printf("[processEventLikes] likes+1 [%+v]\n", err)
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
			_likes, err := cache.EventLikes.DecrementInt64(id, 1)
			if err != nil {
				logger.Debug.Printf("[processEventLikes] likes-1 [%+v]\n", err)
			}
			likes = _likes
		} else {
			_likes, found := cache.EventLikes.Get(id)
			if found {
				likes = _likes.(int64)
			}
		}
	}

	logger.Debug.Printf("[processEventLikes] likes=%v", likes)
	return
}

// @Summary Increase event viewed count whenever this API gets called.
// @Description increase view count by 1 for the given event
// @Tags events
// @Accept  json
// @Produce json
// @Param id path string true "Event ID"
// @param locale path string false "en-US"
// @Success 200 {object} dto.GetEventResponse
// @Failure 400 {object} errors.ErrorInfo
// @Failure 500 {object} errors.ErrorInfo
// @Router /api/hubs-cms/v1/events/{id}/viewed [post]
func EventViewCountHandler(c *gin.Context) {
	param := dto.GetEventRequest{}
	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, errors.EventInvalidID)
		return
	}
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, errors.EventInvalidRequestFormat)
		return
	}

	getDirectusEvent, err := service.GetDirectusEvent(param.ID, param.Locale)

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

	addEventViewCount, addViewCountErrorInfo := service.PostDirectusEventViewCount(getDirectusEvent, param.Locale, false)
	if addViewCountErrorInfo != (errors.ErrorInfo{}) {
		c.JSON(http.StatusInternalServerError, errors.InternalError)
		return
	}

	if pDirectusAccount != nil {
		for i := range pDirectusAccount.LikedEvents {
			if addEventViewCount.ID == pDirectusAccount.LikedEvents[i].EventID {
				addEventViewCount.IsLiked = true
				break
			}
		}
	}

	res := dto.NewEventResponse(addEventViewCount, pDirectusAccount)

	c.JSON(http.StatusOK, res)

}
