package tests

import (
	"errors"
	"hubs-cms-go/client"
	"hubs-cms-go/config"
	"hubs-cms-go/dto"
	"hubs-cms-go/service"
	"net/http"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetMastodonVerifyCredentials(t *testing.T) {
	t.Run("Test get mastodon verify credentials", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		getDirectusAccessTokenJsonResponder, _ := httpmock.NewJsonResponder(200,
			dto.DirectusAuthLoginResponse{
				Data: dto.DirectusAuthLoginResponseData{
					AccessToken:  "xyz",
					Expires:      900000,
					RefreshToken: "abc"}})

		httpmock.RegisterResponder(
			"POST",
			config.GetDirectusAccessTokenURI(),
			getDirectusAccessTokenJsonResponder)

		directusAccessToken, err := service.GetDirectusAccessToken(false)
		assert.Nil(t, err)
		assert.Equal(t, "Bearer xyz", directusAccessToken)

		testMastodonAccount := "mastodon-test-account"
		testDisplayName := "test-name"
		testID := "test-id"
		testOwner := "Test owner"
		testMastodonAvatar := "test-mastodon-avatar"

		mockMastodonAccountInfo := dto.MastodonVerifyCredentialsResponse{
			ID:              testID,
			UserName:        testOwner,
			MastodonAccount: testMastodonAccount,
			DisplayName:     testDisplayName,
			MastodonAvatar:  testMastodonAvatar,
			MastodonToken:   "testMastodonToken",
		}

		getDirectusAccountDataResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			mockMastodonAccountInfo)

		httpmock.RegisterResponder(
			"GET",
			config.GetMastodonVerifyCredentialsURI(),
			getDirectusAccountDataResponder)

		result, err := service.GetMastodonVerifyCredentials(directusAccessToken)
		assert.Nil(t, err)

		assert.Equal(t, mockMastodonAccountInfo.ID, result.ID)
		assert.Equal(t, mockMastodonAccountInfo.MastodonAccount, result.MastodonAccount)
		assert.Equal(t, mockMastodonAccountInfo.DisplayName, result.DisplayName)
		assert.Equal(t, mockMastodonAccountInfo.MastodonAvatar, result.MastodonAvatar)
		assert.Equal(t, mockMastodonAccountInfo.UserName, result.UserName)
	})
}

func TestGetMastodonVerifyCredentialsError(t *testing.T) {
	t.Run("Test get mastodon verify credentials error", func(t *testing.T) {
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		getDirectusAccessTokenJsonResponder, _ := httpmock.NewJsonResponder(200,
			dto.DirectusAuthLoginResponse{
				Data: dto.DirectusAuthLoginResponseData{
					AccessToken:  "xyz",
					Expires:      900000,
					RefreshToken: "abc"}})

		httpmock.RegisterResponder(
			"POST",
			config.GetDirectusAccessTokenURI(),
			getDirectusAccessTokenJsonResponder)

		directusAccessToken, err := service.GetDirectusAccessToken(false)
		assert.Nil(t, err)
		assert.Equal(t, "Bearer xyz", directusAccessToken)

		mockErr := errors.New("some error")
		expectResult := mockErr.Error()
		testErrorResponder := httpmock.NewErrorResponder(mockErr)
		httpmock.RegisterResponder(
			"GET",
			config.GetMastodonVerifyCredentialsURI(),
			testErrorResponder)

		emptyResult, responseErr := service.GetMastodonVerifyCredentials(directusAccessToken)
		assert.Equal(t, dto.MastodonVerifyCredentialsResponse{}, emptyResult)
		assert.True(t, strings.Contains(responseErr.Error(), expectResult))
	})
}

func TestPatchMastodonAccount(t *testing.T) {
	t.Run("Test patch mastodon verify credentials", func(t *testing.T) {
		Init()

		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		getDirectusAccessTokenJsonResponder, _ := httpmock.NewJsonResponder(200,
			dto.DirectusAuthLoginResponse{
				Data: dto.DirectusAuthLoginResponseData{
					AccessToken:  "xyz",
					Expires:      900000,
					RefreshToken: "abc"}})

		httpmock.RegisterResponder(
			"POST",
			config.GetDirectusAccessTokenURI(),
			getDirectusAccessTokenJsonResponder)

		directusAccessToken, err := service.GetDirectusAccessToken(false)
		assert.Nil(t, err)
		assert.Equal(t, "Bearer xyz", directusAccessToken)

		testMastodonAccount := "mastodon-test-account"
		testNewDisplayName := "test-name-2"
		testID := "test-id"
		testOwner := "Test owner"
		testMastodonAvatar := "test-mastodon-avatar"

		mockMastodonAccountInfo := dto.MastodonVerifyCredentialsResponse{
			ID:              testID,
			UserName:        testOwner,
			MastodonAccount: testMastodonAccount,
			DisplayName:     testNewDisplayName,
			MastodonAvatar:  testMastodonAvatar,
			MastodonToken:   "testMastodonToken",
		}

		getDirectusAccountDataResponder, _ := httpmock.NewJsonResponder(http.StatusOK,
			mockMastodonAccountInfo)

		httpmock.RegisterResponder(
			"PATCH",
			config.GetMastodonUpdateCredentialsURI(),
			getDirectusAccountDataResponder)

		mockRequestBody := dto.MastodonPatchAccountRequestBody{
			DisplayName: testNewDisplayName,
		}

		result, err := service.PatchMastodonAccount(directusAccessToken, mockRequestBody)
		assert.Nil(t, err)

		assert.Equal(t, mockMastodonAccountInfo.ID, result.ID)
		assert.Equal(t, mockMastodonAccountInfo.MastodonAccount, result.MastodonAccount)
		assert.Equal(t, mockMastodonAccountInfo.DisplayName, result.DisplayName)
		assert.Equal(t, mockMastodonAccountInfo.MastodonAvatar, result.MastodonAvatar)
		assert.Equal(t, mockMastodonAccountInfo.UserName, result.UserName)

	})
}
