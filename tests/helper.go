package tests

import (
	"hubs-cms-go/dto"

	"fmt"
	"os"

	goCache "github.com/patrickmn/go-cache"
)

func getAssetPath() string {
	return os.Getenv("DIRECTUS_BASE_URI") + "/assets/"
}

func AddNumberToDirectusAvatarResponseData(n int) dto.DirectusAvatarResponseData {
	return dto.DirectusAvatarResponseData{
		ID:       fmt.Sprintf("Test-id-%d", n),
		Snapshot: "Test-snapshot-path",
		GLB:      "Test GLB",
		Owner:    "owner-id",
		Source:   "Test avatar source",
		Title:    "Test Avatar",
		IsPublic: true,
	}
}

func mapMockAvatarResponseDataList(numbers []int, fn func(int) dto.DirectusAvatarResponseData) []dto.DirectusAvatarResponseData {
	result := []dto.DirectusAvatarResponseData{}

	for _, n := range numbers {
		result = append(result, fn(n))
	}

	return result
}

func mockPage(start int64, limit int64) string {
	return fmt.Sprintf("/api/hubs-cms/v1/my-avatars?start=%d&limit=%d", start, limit)
}

func putItemToCache(cacheKey string, source *goCache.Cache, value int64) (string, error) {

	source.Add(cacheKey, value, goCache.NoExpiration)

	if item, ok := source.Get(cacheKey); ok {
		if cacheValue, ok := item.(int64); ok {
			return fmt.Sprintf("%v", cacheValue), nil
		}
	}

	return "", nil
}

func search(ss []string, value string) string {
	for _, s := range ss {
		if s == value {
			return s
		}
	}
	return ""
}

func findById(i interface{}, id string) interface{} {
	switch i.(type) {
	case []dto.Host:
		var hosts = i.([]dto.Host)
		for _, s := range hosts {
			if s.ID == id {
				return s
			}
		}
	case []dto.Speaker:
		var hosts = i.([]dto.Speaker)
		for _, s := range hosts {
			if s.ID == id {
				return s
			}
		}
	case []dto.Room:
		var hosts = i.([]dto.Room)
		for _, s := range hosts {
			if s.ID == id {
				return s
			}
		}
	case []dto.Type:
		var hosts = i.([]dto.Type)
		for _, s := range hosts {
			if s.ID == id {
				return s
			}
		}
	case []dto.Hashtag:
		var hosts = i.([]dto.Hashtag)
		for _, s := range hosts {
			if s.ID == id {
				return s
			}
		}
	}
	return nil
}
