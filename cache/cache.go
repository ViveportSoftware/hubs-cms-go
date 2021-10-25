package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

// Store is an in-memory key value pair for storing cache data
var Store *cache.Cache
var EventLikes *cache.Cache
var RoomLikes *cache.Cache

// Setup initialize the Cache object
func Setup() {
	Store = cache.New(86400*time.Second, 1800*time.Second)
	EventLikes = cache.New(cache.NoExpiration, 1800*time.Second)
	RoomLikes = cache.New(cache.NoExpiration, 1800*time.Second)
}
