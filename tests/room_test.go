package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"hubs-cms-go/cache"
	"hubs-cms-go/client"
	"hubs-cms-go/config"
	"hubs-cms-go/dto"
	hubsErrorInfo "hubs-cms-go/errors"
	"hubs-cms-go/logger"
	"hubs-cms-go/router"
	"hubs-cms-go/service"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	os.Setenv("GO_HTTP_PORT", "9999")
	os.Setenv("LOG_LEVEL", "INFO")
	os.Setenv("ENVIRONMENT", "DEVELOP")
	os.Setenv("MASTODON_BASE_URI", "https://mastodon.test.com")
	os.Setenv("DIRECTUS_BASE_URI", "https://test.directus.app")
	os.Setenv("DIRECTUS_ADMIN_EMAIL", gofakeit.Email())
	os.Setenv("DIRECTUS_ADMIN_PASSWORD", gofakeit.Password(true, true, true, true, true, 10))
	os.Setenv("HUBS_BASE_URI", "https://test.com")

	logger.Setup(config.EnvVariable.LogLevel)
	cache.Setup()

	r := m.Run()

	if r == 0 && testing.CoverMode() != "" {
		c := testing.Coverage() * 100
		l := 0.00
		fmt.Println("=================================================")
		fmt.Println("||               Coverage Report               ||")
		fmt.Println("=================================================")
		fmt.Printf("Cover mode: %s\n", testing.CoverMode())
		fmt.Printf("Coverage  : %.2f %% (Threshold: %.2f %%)\n\n", c, l)
		if c < l {
			fmt.Println("[Tests passed but coverage failed]")
			r = -1
		}
	}

	os.Exit(r)
}

func setPreconditionForViewCount(testID string, isPublic bool, viewCount json.Number) dto.DierctusRoomData {
	nft := dto.DirectusNFTContract{
		ID:         gofakeit.Noun(),
		Name:       "nft-name",
		Source:     "nft-source",
		Platform:   "platform",
		Standard:   "standard",
		Blockchain: "blockchain",
		Address:    "test-address",
	}

	ret := dto.DierctusRoomData{
		ID:          testID,
		Title:       "My Room",
		Description: "My Room Description",
		ViewCount:   viewCount,
		HasNFT:      true,
		IsPublic:    isPublic,
		Passcode:    "0000",
		Translations: []dto.DirectusRoomL10N{
			dto.DirectusRoomL10N{
				ID:          "20",
				Title:       "[zh-TW] My Room Title",
				Description: "[zh-TW] My Room Description",
			},
		},
		Gallery:     dto.DirectusFile{},
		Owner:       gofakeit.Name(),
		HubsID:      gofakeit.UUID(),
		NFTContract: &nft,
		JoinedEvents: []dto.DirectusJoinedEvent{
			dto.DirectusJoinedEvent{
				EventID: gofakeit.UUID(),
			},
			dto.DirectusJoinedEvent{
				EventID: gofakeit.UUID(),
			},
		},
	}

	likes, _ := putItemToCache(testID, cache.RoomLikes, int64(77000))
	ret.LikeCount = json.Number(fmt.Sprintf("%v", likes))
	return ret
}

func TestCheckHubsPasscode(t *testing.T) {
	t.Run("Test check hubs passcode from handler", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testUserID := gofakeit.UUID()
		testLocale := ""
		start := int64(0)
		limit := int64(0)
		testHubsID := gofakeit.Noun()

		mockHubsRequest := dto.HubsPasscodeRequest{
			Passcode: gofakeit.Word(),
		}

		mockRoomResponseA := setPreconditionForViewCount(testUserID, true, "10000")
		mockRoomResponseA.HubsID = testHubsID
		mockRoomResponseA.Passcode = mockHubsRequest.Passcode
		mockResponse := dto.DirectusGetResponse{Data: []dto.DierctusRoomData{
			mockRoomResponseA,
		},
		}
		setUpResponder(http.StatusOK, mockResponse, http.MethodGet, config.GetDirectusGetRoomListURI(nil, testHubsID, testLocale, start, limit))

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/passcode/%s", testHubsID)

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(mockHubsRequest)

		req, err := http.NewRequest(http.MethodPost, testApi, requestBody)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestCheckInvalidHubsPasscode(t *testing.T) {
	t.Run("Handle invalid hubs passcode", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testInvalidHubsID := gofakeit.UUID()

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/passcode/%s", testInvalidHubsID)

		mockHubsRequest := dto.HubsPasscodeRequest{
			Passcode: gofakeit.Word(),
		}

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(mockHubsRequest)

		req, err := http.NewRequest(http.MethodPost, testApi, requestBody)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCheckInvalidHubsPasscodePayload(t *testing.T) {
	t.Run("Handle invalid hubs payload", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testHubsID := gofakeit.Word()

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/passcode/%s", testHubsID)

		mockInvalidRequest := dto.HubsPasscodeRequest{
			Passcode: "",
		}

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(mockInvalidRequest)

		req, err := http.NewRequest(http.MethodPost, testApi, requestBody)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestCheckHubsPasscodeWithInternalError(t *testing.T) {
	t.Run("Test get room list error while check hubs passcode", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testLocale := ""
		start := int64(0)
		limit := int64(0)
		testHubsID := gofakeit.Noun()

		mockErr := errors.New("some error")
		setUpErrorResponder(mockErr, http.MethodGet, config.GetDirectusGetRoomListURI(nil, testHubsID, testLocale, start, limit))

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/passcode/%s", testHubsID)

		mockHubsRequest := dto.HubsPasscodeRequest{
			Passcode: gofakeit.Word(),
		}

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(mockHubsRequest)

		req, err := http.NewRequest(http.MethodPost, testApi, requestBody)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestMultipleHubsIDQueryResult(t *testing.T) {
	t.Run("hubsID should be unique", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testUserID := gofakeit.UUID()
		testLocale := ""
		start := int64(0)
		limit := int64(0)
		testHubsID := gofakeit.Noun()

		mockHubsRequest := dto.HubsPasscodeRequest{
			Passcode: gofakeit.Word(),
		}

		mockRoomResponseA := setPreconditionForViewCount(testUserID, true, "10000")
		mockRoomResponseA.HubsID = testHubsID
		mockRoomResponseA.Passcode = gofakeit.Word()
		mockRoomResponseB := setPreconditionForViewCount(testUserID, true, "20000")
		mockRoomResponseB.HubsID = testHubsID
		mockRoomResponseB.Passcode = mockHubsRequest.Passcode
		mockResponse := dto.DirectusGetResponse{
			Data: []dto.DierctusRoomData{
				mockRoomResponseA,
				mockRoomResponseB,
			},
			Meta: dto.DirectusMeta{
				FilterCount: int64(2),
				TotalCount:  int64(2),
			},
		}

		setUpResponder(http.StatusOK, mockResponse, http.MethodGet, config.GetDirectusGetRoomListURI(nil, testHubsID, testLocale, start, limit))

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/passcode/%s", testHubsID)

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(mockHubsRequest)

		req, err := http.NewRequest(http.MethodPost, testApi, requestBody)
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestPatchDirectusRoom(t *testing.T) {
	t.Run("Test patch directus room from service", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		testID := gofakeit.UUID()

		mockRoomResponse := setPreconditionForViewCount(testID, true, "80000")
		setUpResponder(http.StatusOK, dto.DirectusGetResponse{Data: &mockRoomResponse}, http.MethodPatch, config.GetDirectusGetRoomURISimple(testID))

		likeCountPatchBody := struct {
			LikeCount string `json:"like_count"`
		}{
			LikeCount: "77000",
		}
		// verify service flow
		err := service.PatchDirectusRoom(testID, &likeCountPatchBody)
		assert.Nil(t, err)
	})
}

func TestGetDirectusRoom(t *testing.T) {
	t.Run("Test get directus room from service", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		testID := gofakeit.UUID()
		testLocale := "zh-TW"

		mockRoomResponse := setPreconditionForViewCount(testID, true, "80000")
		setUpResponder(http.StatusOK, dto.DirectusGetResponse{Data: &mockRoomResponse}, http.MethodGet, config.GetDirectusGetRoomURI(testID, testLocale))

		// verify service flow
		responseFormat, err := service.GetDirectusRoom(testID, testLocale)
		assert.Nil(t, err)

		assert.Equal(t, mockRoomResponse.ID, responseFormat.ID)
		assert.Equal(t, "[zh-TW] My Room Title", responseFormat.Title)
		assert.Equal(t, "[zh-TW] My Room Description", responseFormat.Description)
		assert.Equal(t, mockRoomResponse.LikeCount, responseFormat.LikeCount)
		assert.Equal(t, mockRoomResponse.ViewCount, responseFormat.ViewCount)
		assert.Equal(t, mockRoomResponse.IsPublic, responseFormat.IsPublic)
		assert.Equal(t, mockRoomResponse.Owner, responseFormat.Owner)
		assert.Equal(t, mockRoomResponse.NFTContract, responseFormat.NFTContract)
	})
}
func TestGetMyRoomList(t *testing.T) {
	t.Run("Test get my room list from service", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		testAccountID := gofakeit.UUID()
		testRoomID := gofakeit.UUID()
		testLocale := "zh-TW"
		start := int64(0)
		limit := int64(10)
		mockRoomResponse := setPreconditionForViewCount(testRoomID, true, "10000")

		mockResponse := dto.DirectusGetResponse{Data: []dto.DierctusRoomData{
			mockRoomResponse,
		},
		}

		setUpResponder(http.StatusOK, mockResponse, http.MethodGet, config.GetDirectusGetMyRoomListURI(testAccountID, testLocale, start, limit))

		// verify service flow
		directusRoomList, _, err := service.GetDirectusMyRoomList(testAccountID, testLocale, start, limit)
		assert.Equal(t, 1, len(directusRoomList))
		assert.Nil(t, err)
	})
}

func TestGetRoomListService(t *testing.T) {
	t.Run("Test get room lists from service", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		testIDA := gofakeit.UUID()
		testIDB := gofakeit.UUID()
		testIDC := gofakeit.UUID()
		testLocale := "zh-TW"
		hasNFT := false
		start := int64(0)
		limit := int64(10)
		testHubsID := gofakeit.UUID()

		mockRoomResponseA := setPreconditionForViewCount(testIDA, true, "10000")
		mockRoomResponseB := setPreconditionForViewCount(testIDB, true, "20000")
		mockRoomResponseC := setPreconditionForViewCount(testIDC, false, "0")

		mockResponse := dto.DirectusGetResponse{Data: []dto.DierctusRoomData{
			mockRoomResponseA,
			mockRoomResponseB,
			mockRoomResponseC,
		},
		}

		setUpResponder(http.StatusOK, mockResponse, http.MethodGet, config.GetDirectusGetRoomListURI(&hasNFT, testHubsID, testLocale, start, limit))

		// verify service flow
		directusRoomList, _, err := service.GetDirectusRoomList(&hasNFT, testHubsID, testLocale, start, limit)

		assert.Equal(t, 3, len(directusRoomList))
		assert.Nil(t, err)
	})
}

func TestGetRoomListAPI(t *testing.T) {
	t.Run("Retrieve room list", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testID := gofakeit.UUID()
		testLocale := "zh-TW"
		hasNFT := false
		start := int64(0)
		limit := int64(10)
		testHubsID := gofakeit.UUID()

		mockRoomResponseA := setPreconditionForViewCount(testID, true, "10000")
		mockRoomResponseA.HubsID = testHubsID

		mockResponse := dto.DirectusGetResponse{Data: []dto.DierctusRoomData{
			mockRoomResponseA,
		},
		}

		setUpResponder(http.StatusOK, mockResponse, http.MethodGet, config.GetDirectusGetRoomListURI(&hasNFT, testHubsID, testLocale, start, limit))

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/rooms")

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, testApi, nil)

		q := req.URL.Query()
		q.Add("has_nft", strconv.FormatBool(hasNFT))
		q.Add("hubs_id", testHubsID)
		q.Add("locale", testLocale)
		q.Add("start", strconv.FormatInt(start, 10))
		q.Add("limit", strconv.FormatInt(limit, 10))
		req.URL.RawQuery = q.Encode()

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestGetRoomListAPIError(t *testing.T) {
	t.Run("Retrieve room list Error", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testLocale := "zh-TW"
		hasNFT := false
		start := int64(0)
		limit := int64(10)
		testHubsID := gofakeit.UUID()

		mockErr := errors.New("some error")
		setUpErrorResponder(mockErr, http.MethodGet, config.GetDirectusGetRoomListURI(nil, testHubsID, testLocale, start, limit))

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/rooms")

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, testApi, nil)

		q := req.URL.Query()
		q.Add("has_nft", strconv.FormatBool(hasNFT))
		q.Add("hubs_id", testHubsID)
		q.Add("locale", testLocale)
		q.Add("start", strconv.FormatInt(start, 10))
		q.Add("limit", strconv.FormatInt(limit, 10))
		req.URL.RawQuery = q.Encode()

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestGetRoomListAPIWithInvalidStartRange(t *testing.T) {
	t.Run("Handle invalid request of start range", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testLocale := "zh-TW"
		hasNFT := "false"
		start := "-1"
		testHubsID := gofakeit.UUID()

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/rooms")

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, testApi, nil)

		q := req.URL.Query()
		q.Add("has_nft", hasNFT)
		q.Add("hubs_id", testHubsID)
		q.Add("locale", testLocale)
		q.Add("start", start)
		req.URL.RawQuery = q.Encode()

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetRoomListAPIWithInvalidLimitRange(t *testing.T) {
	t.Run("Handle invalid request of limit range", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testLocale := "zh-TW"
		hasNFT := "false"
		testHubsID := gofakeit.UUID()

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/rooms")

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, testApi, nil)

		q := req.URL.Query()
		q.Add("has_nft", hasNFT)
		q.Add("hubs_id", testHubsID)
		q.Add("locale", testLocale)
		q.Add("limit", "-3")
		req.URL.RawQuery = q.Encode()

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetRoomListAPIWithInvalidRequests(t *testing.T) {
	t.Run("Retrieve room detail with invalid requests", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testHubsID := gofakeit.UUID()

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/rooms")

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, testApi, nil)

		q := req.URL.Query()
		q.Add("has_nft", gofakeit.Word())
		q.Add("hubs_id", testHubsID)
		q.Add("locale", gofakeit.Word())
		req.URL.RawQuery = q.Encode()

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetRoomListStartRangeExceedTotalRange(t *testing.T) {
	t.Run("Request room list and make start range exceed total range", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testLocale := "zh-TW"
		hasNFT := false
		start := int64(3)
		limit := int64(10)
		testHubsID := gofakeit.UUID()

		mockResponse := dto.DirectusGetResponse{
			Data: []dto.DierctusRoomData{},
			Meta: dto.DirectusMeta{
				FilterCount: int64(2),
				TotalCount:  int64(2),
			},
		}

		setUpResponder(http.StatusOK, mockResponse, http.MethodGet, config.GetDirectusGetRoomListURI(&hasNFT, testHubsID, testLocale, start, limit))

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/rooms")

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, testApi, nil)

		q := req.URL.Query()
		q.Add("has_nft", strconv.FormatBool(hasNFT))
		q.Add("hubs_id", testHubsID)
		q.Add("locale", testLocale)
		q.Add("start", strconv.FormatInt(start, 10))
		q.Add("limit", strconv.FormatInt(limit, 10))
		req.URL.RawQuery = q.Encode()

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		res := w.Result()
		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)
		responseBody := dto.DierctusRoomData{}

		if err := json.Unmarshal(data, &responseBody); err != nil {
			t.Errorf("Unmarshal，err:%v\n", err)
		}

		assert.Equal(t, dto.DierctusRoomData{}, responseBody)
	})
}

/**
*
* GET /api/hubs-cms/v1/rooms/{id}
*
**/
func TestGetRoomDetailByID(t *testing.T) {
	t.Run("Get room by id ", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		testID := gofakeit.UUID()
		testLocale := "zh-TW"

		mockRoomResponse := setPreconditionForViewCount(testID, true, "80000")
		setUpResponder(http.StatusOK, dto.DirectusGetResponse{Data: &mockRoomResponse}, http.MethodGet, config.GetDirectusGetRoomURI(testID, testLocale))

		serviceResult, err := service.GetDirectusRoom(testID, testLocale)
		assert.Nil(t, err)
		assert.Equal(t, mockRoomResponse.LikeCount, serviceResult.LikeCount)

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/rooms/%s?locale=%s", testID, testLocale)

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, countErr := http.NewRequest(http.MethodGet, testApi, nil)

		if countErr != nil {
			t.Errorf("NewRequest，err:%v\n", countErr)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, countErr)
		assert.Equal(t, http.StatusOK, w.Code)

		res := w.Result()
		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)
		responseFormat := &dto.GetRoomResponse{}

		if err := json.Unmarshal(data, responseFormat); err != nil {
			t.Errorf("Unmarshal，err:%v\n", err)
		}

		assert.Equal(t, mockRoomResponse.ID, responseFormat.ID)
		assert.Equal(t, "[zh-TW] My Room Title", responseFormat.Title)
		assert.Equal(t, "[zh-TW] My Room Description", responseFormat.Description)
		assert.Equal(t, mockRoomResponse.LikeCount, responseFormat.LikeCount)
		assert.Equal(t, mockRoomResponse.ViewCount, responseFormat.ViewCount)
		assert.Equal(t, mockRoomResponse.IsPublic, responseFormat.IsPublic)
		assert.Equal(t, mockRoomResponse.Owner, responseFormat.Owner)
		assert.Equal(t, mockRoomResponse.NFTContract.Standard, responseFormat.NFT.Standard)
		assert.Equal(t, mockRoomResponse.NFTContract.Blockchain, responseFormat.NFT.Blockchain)
		assert.Equal(t, mockRoomResponse.NFTContract.Address, responseFormat.NFT.Address)
		mockJoinedEvents := mockRoomResponse.JoinedEvents
		assert.Equal(t, []string([]string{mockJoinedEvents[0].EventID, mockJoinedEvents[1].EventID}), responseFormat.Events)
	})
}

func TestGetNonPublicRoomDetailByID(t *testing.T) {
	t.Run("Unable to get non public room by id", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testID := gofakeit.UUID()
		testLocale := "zh-TW"

		mockRoomResponse := setPreconditionForViewCount(testID, false, "80000")
		setUpResponder(http.StatusOK, mockRoomResponse, http.MethodGet, config.GetDirectusGetRoomURI(testID, testLocale))

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/rooms/%s?locale=%s", testID, testLocale)

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, testApi, nil)

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

/**
*
* GET /api/hubs-cms/v1/rooms/{id}/viewed
*
**/
func TestPatchedViewCountForPublicRoom(t *testing.T) {
	t.Run("Get and patch view count ", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		testID := gofakeit.UUID()
		testLocale := "zh-TW"

		mockRoomResponse := setPreconditionForViewCount(testID, true, "80000")
		setUpResponder(http.StatusOK, dto.DirectusGetResponse{Data: &mockRoomResponse}, http.MethodGet, config.GetDirectusGetRoomURI(testID, testLocale))

		mockPatchRoomResponse := setPreconditionForViewCount(testID, true, "80001")
		setUpResponder(http.StatusOK, dto.DirectusGetResponse{Data: &mockPatchRoomResponse}, http.MethodPatch, config.GetDirectusGetRoomURI(testID, testLocale))

		// verify service flow
		result, err := service.PostRoomViewCount(mockRoomResponse, testLocale, false)
		assert.Equal(t, hubsErrorInfo.ErrorInfo{}, err)
		assert.Equal(t, mockPatchRoomResponse, result)

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/rooms/%s/viewed?locale=%s", testID, testLocale)

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, countErr := http.NewRequest(http.MethodPost, testApi, nil)

		if countErr != nil {
			t.Errorf("NewRequest，err:%v\n", countErr)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, countErr)
		assert.Equal(t, http.StatusOK, w.Code)

		res := w.Result()
		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)
		responseFormat := &dto.GetRoomResponse{}

		if err := json.Unmarshal(data, responseFormat); err != nil {
			t.Errorf("Unmarshal，err:%v\n", err)
		}

		assert.Equal(t, json.Number("80001"), responseFormat.ViewCount)
	})
}

func TestPatchedViewCountForNotPublicRoom(t *testing.T) {
	t.Run("Get and patch view count for room not public", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		testID := gofakeit.UUID()
		testLocale := "zh-TW"

		mockRoomResponse := setPreconditionForViewCount(testID, false, "80000")
		setUpResponder(http.StatusOK, dto.DirectusGetResponse{Data: &mockRoomResponse}, http.MethodGet, config.GetDirectusGetRoomURI(testID, testLocale))

		mockPatchRoomResponse := setPreconditionForViewCount(testID, false, "80001")
		setUpResponder(http.StatusOK, dto.DirectusGetResponse{Data: &mockPatchRoomResponse}, http.MethodPatch, config.GetDirectusGetRoomURI(testID, testLocale))

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/rooms/%s/viewed?locale=%s", testID, testLocale)

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, countErr := http.NewRequest(http.MethodPost, testApi, nil)

		if countErr != nil {
			t.Errorf("NewRequest，err:%v\n", countErr)
		}
		assert.Nil(t, countErr)

		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)

	})
}

func TestRoomIdNotFound(t *testing.T) {
	t.Run("Test room id not found", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testID := gofakeit.UUID()
		testLocale := "zh-TW"

		mockResponse := dto.DirectusErrorResponse{
			Status: http.StatusForbidden,
			Errors: []dto.DirectusError{
				{
					Message: "You don't have permission to access this.",
					Extensions: dto.DirectusErrorExtension{
						Code: "FORBIDDEN",
					},
				},
			},
		}

		setUpResponder(http.StatusForbidden, mockResponse, http.MethodGet, config.GetDirectusGetRoomURI(testID, testLocale))

		r := router.SetupRouter()
		w := httptest.NewRecorder()
		req, countErr := http.NewRequest(http.MethodPost, "/api/hubs-cms/v1/rooms/"+testID+"/viewed", nil)
		q := req.URL.Query()
		q.Add("locale", testLocale)
		req.URL.RawQuery = q.Encode()

		r.ServeHTTP(w, req)
		assert.Nil(t, countErr)
		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestRoomViewedCountAPIWithInvalidRoomIDFormat(t *testing.T) {
	t.Run("Handle invalid room id format", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testID := gofakeit.Word()
		testLocale := "zh-TW"

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/rooms/%s/viewed", testID)

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodPost, testApi, nil)
		q := req.URL.Query()
		q.Add("locale", testLocale)
		req.URL.RawQuery = q.Encode()

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestRoomViewedCountAPIWithInvalidQuery(t *testing.T) {
	t.Run("Handle invalid query for room viewed count API", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testID := gofakeit.UUID()
		testLocale := gofakeit.CountryAbr() + gofakeit.CountryAbr()

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/rooms/%s/viewed", testID)

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodPost, testApi, nil)
		q := req.URL.Query()
		q.Add("locale", testLocale)
		req.URL.RawQuery = q.Encode()

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAddRoomViewCountAPIError(t *testing.T) {
	t.Run("Handle error for room view count api", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		testID := gofakeit.UUID()
		testLocale := "zh-TW"

		mockErr := errors.New("some error")
		setUpErrorResponder(mockErr, http.MethodPost, config.GetDirectusGetRoomURI(testID, testLocale))

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/rooms/%s/viewed", testID)

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodPost, testApi, nil)

		q := req.URL.Query()
		q.Add("locale", testLocale)
		req.URL.RawQuery = q.Encode()

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
