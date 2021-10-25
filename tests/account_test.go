package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hubs-cms-go/client"
	"hubs-cms-go/config"
	"hubs-cms-go/constant"
	"hubs-cms-go/dto"
	"hubs-cms-go/router"
	"hubs-cms-go/service"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetDirectusAccountData(t *testing.T) {
	t.Run("Test get directus account data", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		mastodonAccount := gofakeit.Email()
		accountData := []dto.DirectusAccountResponseData{
			dto.DirectusAccountResponseData{
				ID:              gofakeit.UUID(),
				MastodonAccount: mastodonAccount,
				MastodonAvatar:  gofakeit.URL(),
				DisplayName:     gofakeit.FirstName(),
				IsAdmin:         false,
				ActiveAvatar:    AddNumberToDirectusAvatarResponseData(1),
				LikedRooms: []dto.DirectusAccountLikedRoom{
					dto.DirectusAccountLikedRoom{
						ID:     "10000",
						RoomID: gofakeit.UUID(),
					},
				},
				LikedEvents: []dto.DirectusAccountLikedEvent{
					dto.DirectusAccountLikedEvent{
						ID:      "20000",
						EventID: gofakeit.UUID(),
					},
				},
			},
		}

		mockDirectusGetResponse := dto.DirectusGetResponse{Data: &accountData}
		setUpResponder(http.StatusOK, mockDirectusGetResponse, http.MethodGet, config.GetDirectusGetAccountURI(mastodonAccount))

		response, err := service.GetDirectusAccountData(mastodonAccount)
		assert.Nil(t, err)

		assert.Equal(t, accountData[0].ID, response.ID)
		assert.Equal(t, accountData[0].MastodonAccount, response.MastodonAccount)
		assert.Equal(t, accountData[0].MastodonAvatar, response.MastodonAvatar)
		assert.Equal(t, accountData[0].DisplayName, response.DisplayName)
		assert.Equal(t, accountData[0].IsAdmin, response.IsAdmin)
		assert.Equal(t, accountData[0].ActiveAvatar.ID, response.ActiveAvatar.ID)
		assert.Equal(t, accountData[0].ActiveAvatar.Snapshot, response.ActiveAvatar.Snapshot)
		assert.Equal(t, accountData[0].ActiveAvatar.GLB, response.ActiveAvatar.GLB)
		assert.Equal(t, accountData[0].ActiveAvatar.Owner, response.ActiveAvatar.Owner)
		assert.Equal(t, accountData[0].ActiveAvatar.Source, response.ActiveAvatar.Source)
		assert.Equal(t, accountData[0].ActiveAvatar.Title, response.ActiveAvatar.Title)
		assert.Equal(t, accountData[0].ActiveAvatar.IsPublic, response.ActiveAvatar.IsPublic)
		assert.Equal(t, 1, len(response.LikedRooms))
		assert.Equal(t, accountData[0].LikedRooms[0].ID, response.LikedRooms[0].ID)
		assert.Equal(t, accountData[0].LikedRooms[0].RoomID, response.LikedRooms[0].RoomID)
		assert.Equal(t, 1, len(response.LikedEvents))
		assert.Equal(t, accountData[0].LikedEvents[0].ID, response.LikedEvents[0].ID)
		assert.Equal(t, accountData[0].LikedEvents[0].EventID, response.LikedEvents[0].EventID)

	})
}

func TestGetDirectusAccount(t *testing.T) {
	t.Run("Test get directus account", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		mastodonAccount := gofakeit.Email()
		accountData := dto.DirectusGetAccountResponse{
			Data: []dto.DirectusAccountResponseData{
				dto.DirectusAccountResponseData{
					ID:              gofakeit.UUID(),
					MastodonAccount: mastodonAccount,
					MastodonAvatar:  gofakeit.URL(),
					DisplayName:     gofakeit.FirstName(),
					IsAdmin:         false,
					ActiveAvatar:    AddNumberToDirectusAvatarResponseData(1),
					LikedRooms: []dto.DirectusAccountLikedRoom{
						dto.DirectusAccountLikedRoom{
							ID:     "10000",
							RoomID: gofakeit.UUID(),
						},
					},
					LikedEvents: []dto.DirectusAccountLikedEvent{
						dto.DirectusAccountLikedEvent{
							ID:      "20000",
							EventID: gofakeit.UUID(),
						},
					},
				},
			},
		}

		setUpResponder(http.StatusOK, accountData, http.MethodGet, config.GetDirectusGetAccountURI(mastodonAccount))

		response, err := service.GetDirectusAccount(mastodonAccount, true)
		assert.Nil(t, err)

		assert.Equal(t, accountData.Data[0].ID, response.ID)
		assert.Equal(t, accountData.Data[0].MastodonAccount, response.MastodonAccount)
		assert.Equal(t, accountData.Data[0].ID, response.ID)
		assert.Equal(t, accountData.Data[0].DisplayName, response.DisplayName)
		assert.Equal(t, accountData.Data[0].IsAdmin, response.IsAdmin)
		assert.Equal(t, accountData.Data[0].ActiveAvatar.ID, response.ActiveAvatar.ID)
		assert.Equal(t, getAssetPath()+accountData.Data[0].ActiveAvatar.Snapshot, response.ActiveAvatar.Snapshot)
		assert.Equal(t, getAssetPath()+accountData.Data[0].ActiveAvatar.GLB, response.ActiveAvatar.GLB)
		assert.Equal(t, accountData.Data[0].ActiveAvatar.Owner, response.ActiveAvatar.Owner)
		assert.Equal(t, accountData.Data[0].ActiveAvatar.Source, response.ActiveAvatar.Source)
		assert.Equal(t, accountData.Data[0].ActiveAvatar.Title, response.ActiveAvatar.Title)
		assert.Equal(t, accountData.Data[0].ActiveAvatar.IsPublic, response.ActiveAvatar.IsPublic)
		assert.Equal(t, 1, len(response.LikedRooms))
		assert.Equal(t, accountData.Data[0].LikedRooms[0].ID, accountData.Data[0].LikedRooms[0].ID)
		assert.Equal(t, accountData.Data[0].LikedRooms[0].RoomID, accountData.Data[0].LikedRooms[0].RoomID)
		assert.Equal(t, 1, len(response.LikedEvents))
		assert.Equal(t, accountData.Data[0].LikedEvents[0].ID, accountData.Data[0].LikedEvents[0].ID)
		assert.Equal(t, accountData.Data[0].LikedEvents[0].EventID, accountData.Data[0].LikedEvents[0].EventID)
	})
}

func TestCreateDirectusAccount(t *testing.T) {
	t.Run("Test create directus account", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		testMastodonAccount := gofakeit.Email()
		testDisplayName := gofakeit.FirstName()
		testID := gofakeit.UUID()
		testOwner := gofakeit.Name()
		testMastodonAvatar := gofakeit.URL()

		mockAccountData := dto.DirectusUpsertAccountResponse{
			Data: dto.DirectusAccountResponseData{
				ID:              testID,
				MastodonAccount: testMastodonAccount,
				MastodonAvatar:  testMastodonAvatar,
				DisplayName:     testDisplayName,
				IsAdmin:         false,
				ActiveAvatar:    AddNumberToDirectusAvatarResponseData(1),
				LikedRooms: []dto.DirectusAccountLikedRoom{
					dto.DirectusAccountLikedRoom{
						ID:     "10000",
						RoomID: gofakeit.UUID(),
					},
				},
				LikedEvents: []dto.DirectusAccountLikedEvent{
					dto.DirectusAccountLikedEvent{
						ID:      "20000",
						EventID: gofakeit.UUID(),
					},
				},
			},
		}

		setUpResponder(http.StatusOK, mockAccountData, http.MethodPost, config.GetDirectusCreateAccountURI())

		mockMastodonAccountInfo := dto.MastodonVerifyCredentialsResponse{
			ID:              testID,
			UserName:        testOwner,
			MastodonAccount: testMastodonAccount,
			DisplayName:     testDisplayName,
			MastodonAvatar:  testMastodonAvatar,
			MastodonToken:   "testMastodonToken",
		}

		result, err := service.CreateDirectusAccount(mockMastodonAccountInfo, false)
		assert.Nil(t, err)

		assert.Equal(t, mockMastodonAccountInfo.ID, result.ID)
		assert.Equal(t, mockMastodonAccountInfo.MastodonAccount, result.MastodonAccount)
		assert.Equal(t, mockMastodonAccountInfo.DisplayName, result.DisplayName)
		assert.Equal(t, mockMastodonAccountInfo.MastodonAvatar, result.MastodonAvatar)
		assert.Equal(t, mockAccountData.Data.IsAdmin, result.IsAdmin)
	})
}

func TestPatchDirectusAccount(t *testing.T) {
	t.Run("Test patch directus account", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()
		regDirTokenRes()

		testMastodonAccount := gofakeit.Email()
		testID := gofakeit.UUID()

		testNewMastodonAvatar := gofakeit.URL()
		testNewActiveAvatarID := gofakeit.UUID()
		testNewDisplayName := gofakeit.FirstName()

		mockPatchAccountRequestBody := dto.PatchAccountRequestBody{
			DisplayName:    testNewDisplayName,
			ActiveAvatarID: testNewActiveAvatarID,
			MastodonAvatar: testNewMastodonAvatar,
		}

		mockAccountData := dto.DirectusUpsertAccountResponse{
			Data: dto.DirectusAccountResponseData{
				ID:              testID,
				MastodonAccount: testMastodonAccount,
				MastodonAvatar:  testNewMastodonAvatar,
				DisplayName:     testNewDisplayName,
				IsAdmin:         false,
				ActiveAvatar:    AddNumberToDirectusAvatarResponseData(1),
				LikedRooms: []dto.DirectusAccountLikedRoom{
					dto.DirectusAccountLikedRoom{
						ID:     "10000",
						RoomID: gofakeit.UUID(),
					},
				},
				LikedEvents: []dto.DirectusAccountLikedEvent{
					dto.DirectusAccountLikedEvent{
						ID:      "20000",
						EventID: gofakeit.UUID(),
					},
				},
			},
		}
		setUpResponder(http.StatusOK, mockAccountData, http.MethodPatch, config.GetDirectusPatchAccountURI(testID))

		result, err := service.PatchDirectusAccount(testID, &mockPatchAccountRequestBody, false)
		assert.Nil(t, err)

		assert.Equal(t, mockPatchAccountRequestBody.DisplayName, result.DisplayName)
		assert.Equal(t, mockPatchAccountRequestBody.MastodonAvatar, result.MastodonAvatar)
		assert.Equal(t, mockAccountData.Data.ID, result.ID)
		assert.Equal(t, mockAccountData.Data.MastodonAccount, result.MastodonAccount)
		assert.Equal(t, mockAccountData.Data.IsAdmin, result.IsAdmin)
	})
}

func TestGetProfileMeAndUpdateDirectusAccount(t *testing.T) {
	t.Run("Test get profile me", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		testMastodonAccount := gofakeit.Email()
		testDisplayName := gofakeit.FirstName()
		testID := gofakeit.UUID()
		testOwner := gofakeit.Name()
		testMastodonAvatar := gofakeit.URL()

		mockMastodonAccountInfo := dto.MastodonVerifyCredentialsResponse{
			ID:              testID,
			UserName:        testOwner,
			MastodonAccount: testMastodonAccount,
			DisplayName:     testDisplayName,
			MastodonAvatar:  testMastodonAvatar,
			MastodonToken:   "testMastodonToken",
		}

		setUpResponder(http.StatusOK, mockMastodonAccountInfo, http.MethodGet, config.GetMastodonVerifyCredentialsURI())

		accountData := dto.DirectusGetAccountResponse{
			Data: []dto.DirectusAccountResponseData{
				dto.DirectusAccountResponseData{
					ID:              testID,
					MastodonAccount: testMastodonAccount,
					MastodonAvatar:  testMastodonAvatar,
					DisplayName:     testDisplayName,
					IsAdmin:         false,
					ActiveAvatar:    AddNumberToDirectusAvatarResponseData(1),
					LikedRooms: []dto.DirectusAccountLikedRoom{
						dto.DirectusAccountLikedRoom{
							ID:     "10000",
							RoomID: gofakeit.UUID(),
						},
					},
					LikedEvents: []dto.DirectusAccountLikedEvent{
						dto.DirectusAccountLikedEvent{
							ID:      "20000",
							EventID: gofakeit.UUID(),
						},
					},
				},
			},
		}

		setUpResponder(http.StatusOK, accountData, http.MethodGet, config.GetDirectusGetAccountURI(testMastodonAccount))

		mockCreateAccountData := dto.DirectusUpsertAccountResponse{
			Data: dto.DirectusAccountResponseData{
				ID:              testID,
				MastodonAccount: testMastodonAccount,
				MastodonAvatar:  testMastodonAvatar,
				DisplayName:     testDisplayName,
				IsAdmin:         false,
				ActiveAvatar:    AddNumberToDirectusAvatarResponseData(1),
				LikedRooms: []dto.DirectusAccountLikedRoom{
					dto.DirectusAccountLikedRoom{
						ID:     "10000",
						RoomID: gofakeit.UUID(),
					},
				},
				LikedEvents: []dto.DirectusAccountLikedEvent{
					dto.DirectusAccountLikedEvent{
						ID:      "20000",
						EventID: gofakeit.UUID(),
					},
				},
			},
		}

		setUpResponder(http.StatusOK, mockCreateAccountData, http.MethodPatch, config.GetDirectusPatchAccountURI(testID))

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/me")

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		req, err := http.NewRequest(http.MethodGet, testApi, nil)
		req.Header.Add(constant.HeaderAuthorization, "Bearer xyz")

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestPatchAccount(t *testing.T) {
	t.Run("Test patch account", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		regDirTokenRes()

		testMastodonAccount := gofakeit.Email()
		testDisplayName := gofakeit.FirstName()
		testID := gofakeit.UUID()
		testOwner := gofakeit.Name()
		testMastodonAvatar := gofakeit.URL()

		mockMastodonAccountInfo := dto.MastodonVerifyCredentialsResponse{
			ID:              testID,
			UserName:        testOwner,
			MastodonAccount: testMastodonAccount,
			DisplayName:     testDisplayName,
			MastodonAvatar:  testMastodonAvatar,
			MastodonToken:   "testMastodonToken",
		}

		mockActivateAvatar := AddNumberToDirectusAvatarResponseData(1)

		setUpResponder(http.StatusOK, mockMastodonAccountInfo, http.MethodGet, config.GetMastodonVerifyCredentialsURI())

		accountData := dto.DirectusGetAccountResponse{
			Data: []dto.DirectusAccountResponseData{
				dto.DirectusAccountResponseData{
					ID:              testID,
					MastodonAccount: testMastodonAccount,
					MastodonAvatar:  testMastodonAvatar,
					DisplayName:     testDisplayName,
					IsAdmin:         false,
					ActiveAvatar:    mockActivateAvatar,
					LikedRooms: []dto.DirectusAccountLikedRoom{
						dto.DirectusAccountLikedRoom{},
					},
					LikedEvents: []dto.DirectusAccountLikedEvent{
						dto.DirectusAccountLikedEvent{},
					},
				},
			},
		}

		setUpResponder(http.StatusOK, accountData, http.MethodGet, config.GetDirectusGetAccountURI(testMastodonAccount))

		testNewActiveAvatarID := gofakeit.UUID()
		testNewMastodonAvatar := gofakeit.URL()
		testNewDisplayName := gofakeit.FirstName()

		mockPatchAccountData := dto.DirectusUpsertAccountResponse{
			Data: dto.DirectusAccountResponseData{
				ID:              testID,
				MastodonAccount: testMastodonAccount,
				MastodonAvatar:  testNewMastodonAvatar,
				DisplayName:     testNewDisplayName,
				IsAdmin:         false,
				ActiveAvatar:    mockActivateAvatar,
				LikedRooms: []dto.DirectusAccountLikedRoom{
					dto.DirectusAccountLikedRoom{},
				},
				LikedEvents: []dto.DirectusAccountLikedEvent{
					dto.DirectusAccountLikedEvent{},
				},
			},
		}

		setUpResponder(http.StatusOK, mockPatchAccountData, http.MethodPatch, config.GetDirectusPatchAccountURI(testID))

		mockDirectusGetAvatarResponse := dto.DirectusGetAvatarResponse{Data: mockActivateAvatar}
		setUpResponder(http.StatusOK, mockDirectusGetAvatarResponse, http.MethodGet, config.GetDirectusGetAvatarURI(testNewActiveAvatarID))

		setUpResponder(http.StatusOK, mockMastodonAccountInfo, http.MethodPatch, config.GetMastodonUpdateCredentialsURI())

		// verify api flow from handler
		testApi := fmt.Sprintf("/api/hubs-cms/v1/accounts/%s", testID)

		r := router.SetupRouter()
		w := httptest.NewRecorder()

		mockPatchAccountRequestBody := dto.PatchAccountRequestBody{
			DisplayName:    testNewDisplayName,
			ActiveAvatarID: testNewActiveAvatarID,
			MastodonAvatar: testNewMastodonAvatar,
		}
		requestBody := new(bytes.Buffer)
		json.NewEncoder(requestBody).Encode(mockPatchAccountRequestBody)

		req, err := http.NewRequest(http.MethodPatch, testApi, requestBody)
		req.Header.Set(constant.HeaderAuthorization, "Bearer xyz")
		req.Header.Set("Content-Type", "application/json")

		if err != nil {
			t.Errorf("NewRequest，err:%v\n", err)
		}

		r.ServeHTTP(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

		res := w.Result()
		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)
		responseBody := &dto.DirectusAccount{}

		if err := json.Unmarshal(data, responseBody); err != nil {
			t.Errorf("Unmarshal，err:%v\n", err)
		}

		assert.Equal(t, mockPatchAccountRequestBody.DisplayName, responseBody.DisplayName)
		assert.Equal(t, mockPatchAccountRequestBody.MastodonAvatar, responseBody.MastodonAvatar)
		assert.Equal(t, mockPatchAccountData.Data.ID, responseBody.ID)
		assert.Equal(t, mockPatchAccountData.Data.MastodonAccount, responseBody.MastodonAccount)
		assert.Equal(t, mockPatchAccountData.Data.IsAdmin, responseBody.IsAdmin)
	})
}
