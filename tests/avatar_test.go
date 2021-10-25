package tests

import (
	"errors"
	"hubs-cms-go/client"
	"hubs-cms-go/config"
	"hubs-cms-go/dto"
	hubsErrorInfo "hubs-cms-go/errors"
	"hubs-cms-go/service"
	"hubs-cms-go/utils"
	"net/http"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetDirectusAvatar(t *testing.T) {
	t.Run("Test get directus avatar", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		testAvatarID := gofakeit.UUID()

		mockDirectusAvatarResponseData := AddNumberToDirectusAvatarResponseData(1)

		mockDirectusGetAvatarResponse := dto.DirectusGetAvatarResponse{Data: mockDirectusAvatarResponseData}

		getDirectusAvatarResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			mockDirectusGetAvatarResponse)

		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusGetAvatarURI(testAvatarID),
			getDirectusAvatarResponder)

		result, err := service.GetDirectusAvatar(testAvatarID, false)
		assert.Nil(t, err)

		directusAssetsHost := getAssetPath()

		assert.Equal(t, mockDirectusAvatarResponseData.ID, result.ID)
		assert.Equal(t, directusAssetsHost+mockDirectusAvatarResponseData.Snapshot, result.Snapshot)
		assert.Equal(t, directusAssetsHost+mockDirectusAvatarResponseData.GLB, result.GLB)
		assert.Equal(t, mockDirectusAvatarResponseData.Owner, result.Owner)
		assert.Equal(t, mockDirectusAvatarResponseData.Source, result.Source)
		assert.Equal(t, mockDirectusAvatarResponseData.Title, result.Title)
		assert.Equal(t, mockDirectusAvatarResponseData.IsPublic, result.IsPublic)

	})
}

func TestGetDirectusAvatarError(t *testing.T) {
	t.Run("Test get directus avatar error", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		testAvatarID := gofakeit.UUID()

		mockErr := errors.New("some error")
		expectResult := mockErr.Error()
		testErrorResponder := httpmock.NewErrorResponder(mockErr)
		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusGetAvatarURI(testAvatarID),
			testErrorResponder)

		emptyResult, responseErr := service.GetDirectusAvatar(testAvatarID, false)
		assert.Equal(t, dto.DirectusAvatar{}, emptyResult)
		assert.True(t, strings.Contains(responseErr.Error(), expectResult))
	})
}

func TestGetPublicAvatarsWithPagination1(t *testing.T) {
	t.Run("get public avatars with pagination start=1&limit=3", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		start := int64(1)
		limit := int64(3)

		num := []int{1, 2}

		mockData := mapMockAvatarResponseDataList(num, AddNumberToDirectusAvatarResponseData)

		mockDirectusGetAvatarResponse := dto.DirectusGetAvatarsResponse{
			Data: mockData,
			Meta: dto.DirectusMeta{
				FilterCount: int64(len(mockData)),
				TotalCount:  int64(len(num)),
			},
		}

		getDirectusAvatarResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			mockDirectusGetAvatarResponse)

		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusGetPublicAvatarURI(start, limit),
			getDirectusAvatarResponder)

		result, err := service.GetPublicAvatars(start, limit, false)

		assert.Equal(t, hubsErrorInfo.ErrorInfo{}, err)
		assert.Equal(t, len(mockData), len(result.Results))
		assert.Equal(t, mockPage(0, limit), result.Pages.Prev)
		assert.Equal(t, "", result.Pages.Next)
	})
}

func TestGetPublicAvatarsWithPagination2(t *testing.T) {
	t.Run("get public avatars with pagination start=0&limit=1", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		start := int64(0)
		limit := int64(1)

		num := []int{1}

		mockData := mapMockAvatarResponseDataList(num, AddNumberToDirectusAvatarResponseData)

		mockDirectusGetAvatarResponse := dto.DirectusGetAvatarsResponse{
			Data: mockData,
			Meta: dto.DirectusMeta{
				FilterCount: int64(len(num)),
				TotalCount:  int64(len(num)),
			},
		}

		getDirectusAvatarResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			mockDirectusGetAvatarResponse)

		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusGetPublicAvatarURI(start, limit),
			getDirectusAvatarResponder)

		result, err := service.GetPublicAvatars(start, limit, false)

		assert.Equal(t, hubsErrorInfo.ErrorInfo{}, err)
		assert.Equal(t, len(mockData), len(result.Results))
		assert.Equal(t, "", result.Pages.Prev)
		assert.Equal(t, "", result.Pages.Next)
	})
}

func TestGetPublicAvatarsWithPagination3(t *testing.T) {
	t.Run("get public avatars with pagination start=3&limit=3", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		start := int64(2)
		limit := int64(1)

		num := []int{3, 4, 5}
		total := int64(len(num) * 2)
		mockData := mapMockAvatarResponseDataList(num, AddNumberToDirectusAvatarResponseData)

		mockDirectusGetAvatarResponse := dto.DirectusGetAvatarsResponse{
			Data: mockData,
			Meta: dto.DirectusMeta{
				FilterCount: total,
				TotalCount:  total,
			},
		}

		getDirectusAvatarResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			mockDirectusGetAvatarResponse)

		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusGetPublicAvatarURI(start, limit),
			getDirectusAvatarResponder)

		result, err := service.GetPublicAvatars(start, limit, false)

		assert.Equal(t, hubsErrorInfo.ErrorInfo{}, err)
		assert.Equal(t, len(mockData), len(result.Results))
		assert.Equal(t, mockPage(utils.Max(0, start-limit), limit), result.Pages.Prev)
		assert.Equal(t, mockPage(start+limit, limit), result.Pages.Next)
	})
}

func TestGetPublicAvatarsError(t *testing.T) {
	t.Run("Test get public avatar error", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		start := int64(0)
		limit := int64(1)

		mockErr := errors.New("some error")
		testErrorResponder := httpmock.NewErrorResponder(mockErr)
		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusGetPublicAvatarURI(start, limit),
			testErrorResponder)

		emptyResult, responseErr := service.GetPublicAvatars(start, limit, false)
		assert.Equal(t, dto.GetAvatarsResponse{}, emptyResult)
		assert.Equal(t, hubsErrorInfo.InternalError, responseErr)
	})
}

func TestGetMyAvatars(t *testing.T) {
	t.Run("Test get my avatars", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		start := int64(1)
		limit := int64(4)
		num := []int{1, 2}
		testAccountID := gofakeit.UUID()

		mockData := mapMockAvatarResponseDataList(num, AddNumberToDirectusAvatarResponseData)

		mockDirectusGetAvatarResponse := dto.DirectusGetAvatarsResponse{
			Data: mockData,
			Meta: dto.DirectusMeta{
				FilterCount: int64(len(num)),
				TotalCount:  int64(len(num)),
			},
		}

		getDirectusAvatarResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			mockDirectusGetAvatarResponse)

		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusGetMyAvatarURI(testAccountID, start, limit),
			getDirectusAvatarResponder)

		result, err := service.GetMyAvatars(testAccountID, start, limit, false)

		assert.Equal(t, hubsErrorInfo.ErrorInfo{}, err)
		assert.Equal(t, len(mockData), len(result.Results))
		assert.Equal(t, mockPage(0, limit), result.Pages.Prev)
		assert.Equal(t, "", result.Pages.Next)
	})

}

func TestGetMyAvatarsError(t *testing.T) {
	t.Run("Test get my avatars error", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		accountID := gofakeit.UUID()
		start := int64(0)
		limit := int64(1)

		mockErr := errors.New("some error")
		testErrorResponder := httpmock.NewErrorResponder(mockErr)
		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusGetMyAvatarURI(accountID, start, limit),
			testErrorResponder)

		emptyResult, responseErr := service.GetMyAvatars(accountID, start, limit, false)
		assert.Equal(t, dto.GetAvatarsResponse{}, emptyResult)
		assert.Equal(t, hubsErrorInfo.InternalError, responseErr)
	})
}

func TestImportAsset(t *testing.T) {
	t.Run("Test upload asset", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		mockDirectusUploadAssetResponse := dto.DirectusUploadAssetResponse{
			Data: dto.DirectusUploadAssetResponseData{
				ID: gofakeit.UUID(),
			},
		}

		importAssetResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			mockDirectusUploadAssetResponse)

		httpmock.RegisterResponder(
			"POST",
			config.GetDirectusImportAssetURI(),
			importAssetResponder)

		mockAssetURL := getAssetPath() + "Test-GLB"

		glbID, err := service.ImportAsset(mockAssetURL, false)
		assert.Nil(t, err)
		assert.Equal(t, mockDirectusUploadAssetResponse.Data.ID, glbID)
	})
}

func TestImportAssetError(t *testing.T) {
	t.Run("Test import asset error", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		mockErr := errors.New("some error")
		expectResult := mockErr.Error()
		testErrorResponder := httpmock.NewErrorResponder(mockErr)
		httpmock.RegisterResponder(
			"POST",
			config.GetDirectusImportAssetURI(),
			testErrorResponder)

		mockAssetURL := getAssetPath() + "Test-GLB"

		glbID, err := service.ImportAsset(mockAssetURL, false)

		assert.True(t, strings.Contains(err.Error(), expectResult))
		assert.Equal(t, "", glbID)
	})
}

func TestCreateAvatar(t *testing.T) {
	t.Run("Test create avatar", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		mockDirectusAvatarResponseData := AddNumberToDirectusAvatarResponseData(1)
		mockAvatarResponse := dto.DirectusGetAvatarResponse{Data: mockDirectusAvatarResponseData}

		getDirectusAvatarResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			mockAvatarResponse)

		httpmock.RegisterResponder(
			"POST",
			config.GetDirectusCreateAvatarURI(),
			getDirectusAvatarResponder)

		mockRequest := dto.DirectusCreateAvatarRequest{
			Snapshot: "Test-snapshot-path",
			GLB:      "Test-GLB",
			Owner:    "owner-id",
			Source:   "Test avatar source",
			Title:    "Test Avatar",
			IsPublic: true,
		}

		createdAvatar, err := service.CreateAvatar(mockRequest, false)

		directusAssetsHost := getAssetPath()

		assert.Nil(t, err)
		assert.Equal(t, directusAssetsHost+mockAvatarResponse.Data.GLB, createdAvatar.GLB)
		assert.Equal(t, mockAvatarResponse.Data.Owner, createdAvatar.Owner)
		assert.Equal(t, mockAvatarResponse.Data.IsPublic, createdAvatar.IsPublic)
		assert.Equal(t, mockAvatarResponse.Data.Source, createdAvatar.Source)
		assert.Equal(t, mockAvatarResponse.Data.Title, createdAvatar.Title)
	})
}

func TestCreateAvatarError(t *testing.T) {
	t.Run("Test create avatar error", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		mockErr := errors.New("some error")
		expectResult := mockErr.Error()
		testErrorResponder := httpmock.NewErrorResponder(mockErr)
		httpmock.RegisterResponder(
			"POST",
			config.GetDirectusCreateAvatarURI(),
			testErrorResponder)

		mockRequest := dto.DirectusCreateAvatarRequest{
			Snapshot: "Test-snapshot-path",
			GLB:      "Test-GLB",
			Owner:    "owner-id",
			Source:   "Test avatar source",
			Title:    "Test Avatar",
			IsPublic: true,
		}
		emptyAvatar, err := service.CreateAvatar(mockRequest, false)

		assert.True(t, strings.Contains(err.Error(), expectResult))
		assert.Equal(t, dto.DirectusAvatar{}, emptyAvatar)
	})
}

func TestGetAvatar(t *testing.T) {
	t.Run("Test get avatar", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		testAvatarID := gofakeit.UUID()

		mockDirectusAvatarResponseData := AddNumberToDirectusAvatarResponseData(1)

		mockAvatarResponse := dto.DirectusGetAvatarResponse{Data: mockDirectusAvatarResponseData}

		getDirectusAvatarResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			mockAvatarResponse)

		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusSingleAvatarURI(testAvatarID),
			getDirectusAvatarResponder)

		getAvatar, err := service.GetAvatar(testAvatarID, false)
		assert.Nil(t, err)
		assert.Equal(t, mockAvatarResponse.Data.ID, getAvatar.Data.ID)
		assert.Equal(t, mockAvatarResponse.Data.IsPublic, getAvatar.Data.IsPublic)
	})
}

func TestGetAvatarError(t *testing.T) {
	t.Run("Test get avatar error", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		testAvatarID := gofakeit.UUID()

		mockErr := errors.New("some error")
		expectResult := mockErr.Error()
		testErrorResponder := httpmock.NewErrorResponder(mockErr)
		httpmock.RegisterResponder(
			"GET",
			config.GetDirectusSingleAvatarURI(testAvatarID),
			testErrorResponder)

		emptyAvatar, err := service.GetAvatar(testAvatarID, false)
		assert.True(t, strings.Contains(err.Error(), expectResult))
		assert.Equal(t, dto.DirectusGetAvatarResponse{}, emptyAvatar)
	})
}

func TestDeleteAvatar(t *testing.T) {
	t.Run("Test delete avatar", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		testAvatarID := gofakeit.UUID()

		mockDeleteAvatarRequest := dto.DirectusDeleteAvatarRequest{
			ID: testAvatarID,
		}

		deleteAvatarResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			mockDeleteAvatarRequest)

		httpmock.RegisterResponder(
			"DELETE",
			config.GetDirectusSingleAvatarURI(testAvatarID),
			deleteAvatarResponder)

		result := service.DeleteAvatar(mockDeleteAvatarRequest.ID, false)
		assert.Nil(t, result)
	})
}

func TestDeleteAvatarError(t *testing.T) {
	t.Run("Test delete avatar error", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		testAvatarID := gofakeit.UUID()

		mockErr := errors.New("some error")
		expectResult := mockErr.Error()
		testErrorResponder := httpmock.NewErrorResponder(mockErr)
		httpmock.RegisterResponder(
			"DELETE",
			config.GetDirectusSingleAvatarURI(testAvatarID),
			testErrorResponder)

		result := service.DeleteAvatar(testAvatarID, false)
		assert.True(t, strings.Contains(result.Error(), expectResult))
	})
}
