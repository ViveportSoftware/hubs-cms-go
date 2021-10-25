package main

import (
	"errors"
	"fmt"
	"hubs-cms-go/cache"
	"hubs-cms-go/client"
	"hubs-cms-go/config"
	"hubs-cms-go/constant"
	"hubs-cms-go/dto"
	"hubs-cms-go/router"
	"hubs-cms-go/service"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

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

func TestHealth(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestVersion(t *testing.T) {
	r := router.SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/version", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetDirectusAccessToken(t *testing.T) {
	t.Run("Get directus access token", func(t *testing.T) {
		cache.Setup()
		cache.Store.Delete(constant.CacheKeyDirectusAccessToken)
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

		cachedAccessToken, found := cache.Store.Get(constant.CacheKeyDirectusAccessToken)
		assert.True(t, found)
		assert.Equal(t, "Bearer xyz", cachedAccessToken.(string))
	})
}

func TestGetDirectusAccessTokenError(t *testing.T) {
	t.Run("Test directus access token error", func(t *testing.T) {
		cache.Setup()
		cache.Store.Delete(constant.CacheKeyDirectusAccessToken)
		httpmock.ActivateNonDefault(client.RestyClient.GetClient())
		defer httpmock.DeactivateAndReset()

		mockErr := errors.New("some error")
		expectResult := mockErr.Error()
		testErrorResponder := httpmock.NewErrorResponder(mockErr)
		httpmock.RegisterResponder(
			"POST",
			config.GetDirectusAccessTokenURI(),
			testErrorResponder)

		emptyResult, responseErr := service.GetDirectusAccessToken(true)
		assert.Equal(t, "", emptyResult)
		assert.True(t, strings.Contains(responseErr.Error(), expectResult))
	})
}
