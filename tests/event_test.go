package tests

import (
	"encoding/json"
	"fmt"
	"hubs-cms-go/cache"
	"hubs-cms-go/client"
	"hubs-cms-go/config"
	"hubs-cms-go/dto"
	"hubs-cms-go/errors"
	"hubs-cms-go/logger"
	"hubs-cms-go/router"
	"hubs-cms-go/service"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func SetupRouter() *gin.Engine {
	r := router.SetupRouter()
	gin.SetMode(gin.TestMode)
	return r
}

func Init() *gin.Engine {
	config.Setup()
	cache.Setup()
	logger.Setup(config.EnvVariable.LogLevel)
	return SetupRouter()
}

func TestGetEventsInvalidStart2(t *testing.T) {
	testRouter := Init()

	req, err := http.NewRequest("GET", "/api/hubs-cms/v1/events?start=-1", nil)

	if err != nil {
		t.Errorf("NewRequest，err:%v\n", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	b, _ := json.Marshal(errors.EventInvalidStart)
	s := resp.Body.Bytes()

	assert.Equal(t, string(b), string(s[:]))
}

func TestGetEventsInvalidLimit(t *testing.T) {
	testRouter := Init()

	req, err := http.NewRequest("GET", "/api/hubs-cms/v1/events?limit=-1", nil)

	if err != nil {
		t.Errorf("NewRequest，err:%v\n", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	b, _ := json.Marshal(errors.EventInvalidLimit)
	s := resp.Body.Bytes()

	assert.Equal(t, string(b), string(s[:]))
}

func TestGetEventsInvalidLimit2(t *testing.T) {
	testRouter := Init()

	req, err := http.NewRequest("GET", "/api/hubs-cms/v1/events?limit=100", nil)

	if err != nil {
		t.Errorf("NewRequest，err:%v\n", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	b, _ := json.Marshal(errors.EventInvalidLimit)
	s := resp.Body.Bytes()

	assert.Equal(t, string(b), string(s[:]))
}

func TestGetEventsWithWrongLocale(t *testing.T) {
	testRouter := Init()

	req, err := http.NewRequest("GET", "/api/hubs-cms/v1/events?locale=xxx-ooo", nil)

	if err != nil {
		t.Errorf("NewRequest，err:%v\n", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	//check response body
	b, _ := json.Marshal(errors.EventInvalidRequestFormat)
	s := resp.Body.Bytes()
	assert.Equal(t, string(b), string(s[:]))
}

func TestGetEventByIDWrongParam(t *testing.T) {
	testRouter := Init()

	testID := gofakeit.UUID()
	testApi := fmt.Sprintf("/api/hubs-cms/v1/events/%s?locale=xxxxxxx", testID)

	req, err := http.NewRequest("GET", testApi, nil)

	if err != nil {
		t.Errorf("NewRequest，err:%v\n", err)
	}

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestInvalidEventId(t *testing.T) {
	t.Run("Invalid event id", func(t *testing.T) {
		testRouter := Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		testID := gofakeit.UUID()
		testLocale := "zh-TW"

		getRoomResponder, _ := httpmock.NewJsonResponder(http.StatusForbidden,
			dto.DirectusErrorResponse{
				Status: http.StatusForbidden,
				Errors: []dto.DirectusError{
					{
						Message: "You don't have permission to access this.",
						Extensions: dto.DirectusErrorExtension{
							Code: "FORBIDDEN",
						},
					},
				},
			})

		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusGetEventURI(testID, testLocale),
			getRoomResponder)

		w := httptest.NewRecorder()
		testApi := fmt.Sprintf("/api/hubs-cms/v1/events/%s?locale=%s", testID, testLocale)

		req, countErr := http.NewRequest(http.MethodGet, testApi, nil)

		testRouter.ServeHTTP(w, req)
		assert.Nil(t, countErr)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestGetEvent(t *testing.T) {
	t.Run("Get event", func(t *testing.T) {
		testRouter := Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		now := time.Now()
		testID := gofakeit.UUID()
		testLocale := "zh-TW"
		testGallery := gofakeit.UUID()
		testTitle := "testTitle"
		testDescription := "testDescription"
		testAgenda := "testAgenda"
		testTitleTranslation := "testTitleTranslation"
		testDescriptionTranslation := "testDescriptionTranslation"
		testAgendaTranslation := "testAgendaTranslation"
		testStartTime := now
		testEndTime := now.Add(time.Hour * 24)
		testIsPromoted := true
		testLikeCount := json.Number("100")
		testViewCount := json.Number("100000")
		testCategory := dto.DirectusCategory{
			ID:   gofakeit.UUID(),
			Name: "category name",
			Translations: []dto.CategoryTranslations{
				{
					Name: "category name us",
				},
			},
		}
		testParticipate := dto.ParticipateID{
			ID:          gofakeit.UUID(),
			Name:        "host name",
			Description: "host description",
			Translations: []dto.ParticipateIDTranslations{
				{Name: "host name us", Description: "host description us"},
			},
		}
		testSpeaker := dto.SpeakersParticipateID{
			ID:          gofakeit.UUID(),
			Name:        "host name",
			Description: "host description",
			Translations: []dto.ParticipateIDTranslations{
				{
					Name: "host name us",
				},
			},
			Image: "1234567890",
		}
		testRoom := dto.RoomID{
			ID:          gofakeit.UUID(),
			Title:       "host name",
			Gallery:     "this-is-gallery",
			Description: "room description",
			HubsID:      "This-is-hubs-ID",
			Translations: []dto.RoomIDTranslations{
				{
					Title:       "this is gallery us",
					Description: "room description us",
				},
			},
		}
		testTranslation := dto.Translations{
			Title:       testTitleTranslation,
			Description: testDescriptionTranslation,
			Agenda:      testAgendaTranslation,
		}
		testVideoId := dto.DirectusVideo{
			dto.DirectusVideoID{
				CoverImage: "http://testlink/cover_image.png",
				Mp4:        "http://testlink/test.mp4",
				Webm:       "http://testlink/test.webm",
			},
		}
		testImageId := dto.DirectusFilesID{
			ID: "DirectusFilesID 1234 Image",
		}
		testHashtag := dto.HashtagID{
			ID:   gofakeit.UUID(),
			Name: "live share",
		}
		testAccountID := dto.AccountID{
			ID:          gofakeit.UUID(),
			DisplayName: "hello world",
			IsAdmin:     true,
		}
		mockDirectusEventResponse := dto.DirectusGetEventByIDResponse{
			Data: dto.DirectusEventResponseData{
				ID:          testID,
				Gallery:     testGallery,
				Title:       testTitle,
				Description: testDescription,
				Agenda:      testAgenda,
				IsPromoted:  testIsPromoted,
				StartTime:   testStartTime,
				EndTime:     testEndTime,
				LikeCount:   testLikeCount,
				ViewCount:   testViewCount,
				Category:    testCategory,
				Hosts: []dto.DirectusHost{
					{
						ParticipateID: testParticipate,
					},
				},
				Speakers: []dto.DirectusSpeaker{
					{
						ParticipateID: testSpeaker,
					},
				},
				Rooms: []dto.DirectusRoom{
					{
						RoomID: testRoom,
					},
				},
				Translations: []dto.Translations{
					testTranslation,
				},
				Images: []dto.DirectusFilesID{
					testImageId,
				},
				Videos: []dto.DirectusVideo{
					testVideoId,
				},
				HostedAccounts: []dto.HostedAccount{
					{
						AccountID: testAccountID,
					},
				},
				Hashtags: []dto.DirectusHashtag{
					{
						HashtagID: testHashtag,
					},
				},
			},
		}

		getRoomResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			&mockDirectusEventResponse)

		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusGetEventURI(testID, testLocale),
			getRoomResponder)

		res := httptest.NewRecorder()
		testApi := fmt.Sprintf("/api/hubs-cms/v1/events/%s?locale=%s", testID, testLocale)
		req, err := http.NewRequest(http.MethodGet, testApi, nil)

		testRouter.ServeHTTP(res, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Code)

		result := res.Result()
		defer result.Body.Close()
		body, _ := ioutil.ReadAll(result.Body)

		event := &dto.GetEventResponse{}
		if err := json.Unmarshal(body, event); err != nil {
			t.Errorf("Unmarshal，err:%v\n", err)
		}

		assert.Equal(t, testID, event.ID)
		assert.Equal(t, testIsPromoted, event.IsPromoted)
		assert.Equal(t, testTitleTranslation, event.Title)
		assert.Equal(t, testAgendaTranslation, event.Agenda)
		assert.Equal(t, testDescriptionTranslation, event.Description)
		assert.Equal(t, testEndTime.GoString(), event.EndTime.GoString())
		assert.Equal(t, fmt.Sprintf("%s/assets/%s", config.EnvVariable.DirectusBaseURI, testGallery), event.Gallery)
		assert.Equal(t, testCategory.ID, event.Category.ID)
		assert.Equal(t, testCategory.Translations[0].Name, event.Category.Value)
		assert.Equal(t, testParticipate.ID, event.Hosts[0].ID)
		assert.Equal(t, testParticipate.Translations[0].Name, event.Hosts[0].DisplayName)
		assert.Equal(t, testSpeaker.ID, event.Speakers[0].ID)
		assert.Equal(t, testSpeaker.Translations[0].Name, event.Speakers[0].DisplayName)
		assert.Equal(t, testRoom.ID, event.Rooms[0].ID)
		assert.Equal(t, fmt.Sprintf("%s/assets/%s", config.EnvVariable.DirectusBaseURI, testRoom.Gallery), event.Rooms[0].Gallery)
		assert.Equal(t, testRoom.Translations[0].Title, event.Rooms[0].Title)
		assert.Equal(t, testHashtag.ID, event.Hashtags[0].ID)
		assert.Equal(t, testHashtag.Name, event.Hashtags[0].Value)
		assert.Equal(t, os.Getenv("HUBS_BASE_URI")+"/"+testRoom.HubsID, event.Rooms[0].HubsURL)
	})
}

func genDirEvent(id string, viewCount string) *dto.DirectusEventResponseData {
	now := time.Now()
	testID := id
	testGallery := gofakeit.UUID()
	testTitle := "testTitle"
	testDescription := "testDescription"
	testAgenda := "testAgenda"
	testTitleTranslation := "testTitleTranslation"
	testDescriptionTranslation := "testDescriptionTranslation"
	testAgendaTranslation := "testAgendaTranslation"
	testStartTime := now
	testEndTime := now.Add(time.Hour * 24)
	testIsPromoted := true
	testLikeCount := json.Number("100")
	testViewCount := json.Number(viewCount)
	testCategory := dto.DirectusCategory{
		ID:   gofakeit.UUID(),
		Name: "category name",
		Translations: []dto.CategoryTranslations{
			{
				Name: "category name us",
			},
		},
	}
	testParticipate := dto.ParticipateID{
		ID:          gofakeit.UUID(),
		Name:        "host name",
		Description: "host description",
		Translations: []dto.ParticipateIDTranslations{
			{Name: "host name us", Description: "host description us"},
		},
	}
	testSpeaker := dto.SpeakersParticipateID{
		ID:          gofakeit.UUID(),
		Name:        "host name",
		Description: "host description",
		Translations: []dto.ParticipateIDTranslations{
			{
				Name: "host name us",
			},
		},
		Image: "1234567890",
	}
	testRoom := dto.RoomID{
		ID:          gofakeit.UUID(),
		Title:       "host name",
		Gallery:     "this-is-gallery",
		Description: "room description",
		HubsID:      "1234567890",
		Translations: []dto.RoomIDTranslations{
			{
				Title:       "this is gallery us",
				Description: "room description us",
			},
		},
	}
	testTranslation := dto.Translations{
		Title:       testTitleTranslation,
		Description: testDescriptionTranslation,
		Agenda:      testAgendaTranslation,
	}
	testVideoId := dto.DirectusVideo{
		dto.DirectusVideoID{
			CoverImage: "http://testlink/cover_image.png",
			Mp4:        "http://testlink/test.mp4",
			Webm:       "http://testlink/test.webm",
		},
	}
	testImageId := dto.DirectusFilesID{
		ID: "DirectusFilesID 1234 Image",
	}
	testHashtag := dto.HashtagID{
		ID:   gofakeit.UUID(),
		Name: "live share",
	}
	testAccountID := dto.AccountID{
		ID:          gofakeit.UUID(),
		DisplayName: "hello world",
		IsAdmin:     true,
	}
	mockDirectusEventDataResponse := dto.DirectusEventResponseData{
		ID:          testID,
		Gallery:     testGallery,
		Title:       testTitle,
		Description: testDescription,
		Agenda:      testAgenda,
		IsPromoted:  testIsPromoted,
		StartTime:   testStartTime,
		EndTime:     testEndTime,
		LikeCount:   testLikeCount,
		ViewCount:   testViewCount,
		Category:    testCategory,
		Hosts: []dto.DirectusHost{
			{
				ParticipateID: testParticipate,
			},
		},
		Speakers: []dto.DirectusSpeaker{
			{
				ParticipateID: testSpeaker,
			},
		},
		Rooms: []dto.DirectusRoom{
			{
				RoomID: testRoom,
			},
		},
		Translations: []dto.Translations{
			testTranslation,
		},
		Images: []dto.DirectusFilesID{
			testImageId,
		},
		Videos: []dto.DirectusVideo{
			testVideoId,
		},
		HostedAccounts: []dto.HostedAccount{
			{
				AccountID: testAccountID,
			},
		},
		Hashtags: []dto.DirectusHashtag{
			{
				HashtagID: testHashtag,
			},
		},
	}

	return &mockDirectusEventDataResponse
}

func TestGetMultipleEventsResults(t *testing.T) {
	t.Run("Get events", func(t *testing.T) {
		testRouter := Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		testLocale := "zh-TW"
		testStatus := "closed|opened|soon"
		testOffset := int64(0)
		testLimit := int64(2)
		testApi := fmt.Sprintf("/api/hubs-cms/v1/events?locale=%s&status=%s&offset=%d&limit=%d", testLocale, testStatus, testOffset, testLimit)

		mockDirectusResponse := dto.DirectusGetEventsResponse{
			Meta: dto.DirectusMeta{
				FilterCount: 4,
				TotalCount:  4,
			},
			Data: []dto.DirectusEventResponseData{
				*genDirEvent("1", "10000"),
				*genDirEvent("2", "20000"),
			},
		}

		getEventsResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			&mockDirectusResponse)

		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusGetEventsURI(testLocale, testStatus, testOffset, testLimit),
			getEventsResponder)

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, testApi, nil)

		testRouter.ServeHTTP(res, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Code)

		result := res.Result()
		defer result.Body.Close()
		body, _ := ioutil.ReadAll(result.Body)

		events := &dto.GetEventsResponse{}
		if err := json.Unmarshal(body, events); err != nil {
			t.Errorf("Unmarshal，err:%v\n", err)
		}

		assert.True(t, len(events.Pages.Next) > 0)
		assert.True(t, len(events.Pages.Prev) == 0)
		assert.Equal(t, "1", events.Results[0].ID)
		assert.Equal(t, "2", events.Results[1].ID)
		assert.Equal(t, os.Getenv("HUBS_BASE_URI")+"/"+"1234567890", events.Results[1].Rooms[0].HubsURL)
		//todo: add some other asserts
	})
}

func TestGetEventsEmpty(t *testing.T) {
	t.Run("Get empty event list", func(t *testing.T) {
		testRouter := Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		testLocale := "zh-TW"
		testStatus := "closed|opened|soon"
		testOffset := int64(0)
		testLimit := int64(2)
		testApi := fmt.Sprintf("/api/hubs-cms/v1/events?locale=%s&status=%s&offset=%d&limit=%d", testLocale, testStatus, testOffset, testLimit)

		mockDirectusResponse := dto.DirectusGetEventsResponse{
			Meta: dto.DirectusMeta{
				FilterCount: 0,
				TotalCount:  0,
			},
			Data: []dto.DirectusEventResponseData{},
		}

		getEventsResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			&mockDirectusResponse)

		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusGetEventsURI(testLocale, testStatus, testOffset, testLimit),
			getEventsResponder)

		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodGet, testApi, nil)

		testRouter.ServeHTTP(res, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Code)

		result := res.Result()
		defer result.Body.Close()
		body, _ := ioutil.ReadAll(result.Body)

		events := &dto.GetEventsResponse{}
		if err := json.Unmarshal(body, events); err != nil {
			t.Errorf("Unmarshal，err:%v\n", err)
		}

		assert.True(t, len(events.Pages.Next) == 0)
		assert.True(t, len(events.Pages.Prev) == 0)
		assert.NotNil(t, events.Results)
		assert.Empty(t, events.Results)
	})
}

func TestPostEventViewCount(t *testing.T) {
	t.Run("Post event", func(t *testing.T) {
		testRouter := Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		now := time.Now()
		testID := gofakeit.UUID()
		testLocale := "zh-TW"
		testGallery := gofakeit.UUID()
		testTitle := "testTitle"
		testDescription := "testDescription"
		testAgenda := "testAgenda"
		testTitleTranslation := "testTitleTranslation"
		testDescriptionTranslation := "testDescriptionTranslation"
		testAgendaTranslation := "testAgendaTranslation"
		testStartTime := now
		testEndTime := now.Add(time.Hour * 24)
		testIsPromoted := true
		testLikeCount := json.Number("100")
		testViewCount := json.Number("100000")
		testCategory := dto.DirectusCategory{
			ID:   gofakeit.UUID(),
			Name: "category name",
			Translations: []dto.CategoryTranslations{
				{
					Name: "category name us",
				},
			},
		}
		testParticipate := dto.ParticipateID{
			ID:          gofakeit.UUID(),
			Name:        "host name",
			Description: "host description",
			Translations: []dto.ParticipateIDTranslations{
				{Name: "host name us", Description: "host description us"},
			},
		}
		testSpeaker := dto.SpeakersParticipateID{
			ID:          gofakeit.UUID(),
			Name:        "host name",
			Description: "host description",
			Translations: []dto.ParticipateIDTranslations{
				{
					Name: "host name us",
				},
			},
			Image: "1234567890",
		}
		testRoom := dto.RoomID{
			ID:          gofakeit.UUID(),
			Title:       "host name",
			Gallery:     "this-is-gallery",
			Description: "room description",
			HubsID:      "This-is-hubs-ID",
			Translations: []dto.RoomIDTranslations{
				{
					Title:       "this is gallery us",
					Description: "room description us",
				},
			},
		}
		testTranslation := dto.Translations{
			Title:       testTitleTranslation,
			Description: testDescriptionTranslation,
			Agenda:      testAgendaTranslation,
		}
		testVideoId := dto.DirectusVideo{
			dto.DirectusVideoID{
				CoverImage: "http://testlink/cover_image.png",
				Mp4:        "http://testlink/test.mp4",
				Webm:       "http://testlink/test.webm",
			},
		}
		testImageId := dto.DirectusFilesID{
			ID: "DirectusFilesID 1234 Image",
		}
		testHashtag := dto.HashtagID{
			ID:   gofakeit.UUID(),
			Name: "live share",
		}
		testAccountID := dto.AccountID{
			ID:          gofakeit.UUID(),
			DisplayName: "hello world",
			IsAdmin:     true,
		}
		mockDirectusEventResponse := dto.DirectusGetEventByIDResponse{
			Data: dto.DirectusEventResponseData{
				ID:          testID,
				Gallery:     testGallery,
				Title:       testTitle,
				Description: testDescription,
				Agenda:      testAgenda,
				IsPromoted:  testIsPromoted,
				StartTime:   testStartTime,
				EndTime:     testEndTime,
				LikeCount:   testLikeCount,
				ViewCount:   testViewCount,
				Category:    testCategory,
				Hosts: []dto.DirectusHost{
					{
						ParticipateID: testParticipate,
					},
				},
				Speakers: []dto.DirectusSpeaker{
					{
						ParticipateID: testSpeaker,
					},
				},
				Rooms: []dto.DirectusRoom{
					{
						RoomID: testRoom,
					},
				},
				Translations: []dto.Translations{
					testTranslation,
				},
				Images: []dto.DirectusFilesID{
					testImageId,
				},
				Videos: []dto.DirectusVideo{
					testVideoId,
				},
				HostedAccounts: []dto.HostedAccount{
					{
						AccountID: testAccountID,
					},
				},
				Hashtags: []dto.DirectusHashtag{
					{
						HashtagID: testHashtag,
					},
				},
			},
		}

		getEventResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			&mockDirectusEventResponse)

		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusGetEventURI(testID, testLocale),
			getEventResponder)

		mockPostDirectusEventResponse := mockDirectusEventResponse
		mockPostDirectusEventResponse.Data.ViewCount = json.Number("100001")

		patchEventResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			&mockPostDirectusEventResponse)

		httpmock.RegisterResponder(
			"PATCH",
			config.GetDirectusGetEventURI(testID, testLocale),
			patchEventResponder)

		// 		r := router.SetupRouter()
		res := httptest.NewRecorder()
		req, err := http.NewRequest(http.MethodPost, "/api/hubs-cms/v1/events/"+testID+"/viewed?locale="+testLocale, nil)

		testRouter.ServeHTTP(res, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, res.Code)

		result := res.Result()
		defer result.Body.Close()
		body, _ := ioutil.ReadAll(result.Body)

		event := &dto.GetEventResponse{}
		if err := json.Unmarshal(body, event); err != nil {
			t.Errorf("Unmarshal，err:%v\n", err)
		}

		assert.Equal(t, mockPostDirectusEventResponse.Data.ViewCount, event.ViewCount)
		assert.Equal(t, testID, event.ID)
		assert.Equal(t, testIsPromoted, event.IsPromoted)
		assert.Equal(t, testTitleTranslation, event.Title)
		assert.Equal(t, testAgendaTranslation, event.Agenda)
		assert.Equal(t, testDescriptionTranslation, event.Description)
		assert.Equal(t, testEndTime.GoString(), event.EndTime.GoString())
		assert.Equal(t, fmt.Sprintf("%s/assets/%s", config.EnvVariable.DirectusBaseURI, testGallery), event.Gallery)
		assert.Equal(t, testCategory.ID, event.Category.ID)
		assert.Equal(t, testCategory.Translations[0].Name, event.Category.Value)
		assert.Equal(t, testParticipate.ID, event.Hosts[0].ID)
		assert.Equal(t, testParticipate.Translations[0].Name, event.Hosts[0].DisplayName)
		assert.Equal(t, testSpeaker.ID, event.Speakers[0].ID)
		assert.Equal(t, testSpeaker.Translations[0].Name, event.Speakers[0].DisplayName)
		assert.Equal(t, testRoom.ID, event.Rooms[0].ID)
		assert.Equal(t, fmt.Sprintf("%s/assets/%s", config.EnvVariable.DirectusBaseURI, testRoom.Gallery), event.Rooms[0].Gallery)
		assert.Equal(t, testRoom.Translations[0].Title, event.Rooms[0].Title)
		assert.Equal(t, testHashtag.ID, event.Hashtags[0].ID)
		assert.Equal(t, testHashtag.Name, event.Hashtags[0].Value)
		assert.Equal(t, os.Getenv("HUBS_BASE_URI")+"/"+testRoom.HubsID, event.Rooms[0].HubsURL)
	})
}

func TestPatchDirectusEvent(t *testing.T) {
	t.Run("Test patch directus event", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		testID := "test-event-id"
		mockDirectusResponse := dto.DirectusGetEventsResponse{
			Meta: dto.DirectusMeta{
				FilterCount: 4,
				TotalCount:  4,
			},
			Data: []dto.DirectusEventResponseData{
				*genDirEvent(testID, "10000"),
			},
		}

		getEventsResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			&mockDirectusResponse)

		httpmock.RegisterResponder(
			"PATCH",
			config.GetDirectusGetEventURISimple(testID),
			getEventsResponder)

		likeCountPatchBody := struct {
			LikeCount string `json:"like_count"`
		}{
			LikeCount: "100",
		}
		// verify service flow
		err := service.PatchDirectusEvent(testID, &likeCountPatchBody)
		assert.Nil(t, err)
	})
}
