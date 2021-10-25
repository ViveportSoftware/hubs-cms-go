package jobs

import (
	"hubs-cms-go/cache"
	"hubs-cms-go/config"
	"hubs-cms-go/logger"
	"hubs-cms-go/service"
	"sync"
	"time"

	goCache "github.com/patrickmn/go-cache"
	"github.com/robfig/cron"
)

func Setup() {
	initialCache()
	startCron()
}

func initialCache() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		startTime := time.Now()
		likedEventCount, totalLikes, _ := RestoreEventLikeCount()
		logger.Debug.Printf("[initialCache] %v liked events has been restored, total likes=%v, duration=%v", likedEventCount, totalLikes, time.Since(startTime))
	}()
	go func() {
		defer wg.Done()
		startTime := time.Now()
		likedRoomCount, totalLikes, _ := RestoreRoomLikeCount()
		logger.Debug.Printf("[initialCache] %v liked rooms has been restored, total likes=%v, duration=%v", likedRoomCount, totalLikes, time.Since(startTime))
	}()
	wg.Wait()
}

func startCron() {

	c := cron.New()

	c.AddFunc(config.EnvVariable.EventBackupInterval, func() {
		startTime := time.Now()
		count, likes := service.BackupLikeCount("event", cache.EventLikes.Items(), "room", cache.RoomLikes.Items())
		logger.Debug.Printf("[startCron] Backup %v items(event+room) %v likes, duration=%v", count, likes, time.Since(startTime))
	})

	c.Start()
}

func RestoreEventLikeCount() (eventCount, totalLikes int, err error) {

	var offset = int64(0)
	var pageSize = int64(100)

	for {
		accounts, filterCount, _ := service.GetAccountsLikedStuff("event", offset, pageSize)
		if len(accounts) > 0 {
			for _, account := range accounts {
				for _, event := range account.LikedEvents {
					if _, err := cache.EventLikes.IncrementInt64(event.EventID, 1); err != nil {
						//id does not exist, create it
						if err = cache.EventLikes.Add(event.EventID, int64(1), goCache.NoExpiration); err != nil {
							//already exists, should not happen
							logger.Error.Printf("[RestoreEventLikeCount] like+1 [%+v]\n", err)
							continue
						}
						eventCount++
					}
					totalLikes++
				}
			}
		}
		if offset+pageSize >= filterCount {
			break
		}
		offset = offset + pageSize
	}
	return
}

func RestoreRoomLikeCount() (eventCount, totalLikes int, err error) {

	var offset = int64(0)
	var pageSize = int64(100)

	for {
		accounts, filterCount, _ := service.GetAccountsLikedStuff("room", offset, pageSize)
		if len(accounts) > 0 {
			for _, account := range accounts {
				for _, event := range account.LikedRooms {
					if _, err := cache.RoomLikes.IncrementInt64(event.RoomID, 1); err != nil {
						//id does not exist, create it
						if err = cache.RoomLikes.Add(event.RoomID, int64(1), goCache.NoExpiration); err != nil {
							//already exists, should not happen
							logger.Error.Printf("[RestoreRoomLikeCount] like+1 [%+v]\n", err)
							continue
						}
						eventCount++
					}
					totalLikes++
				}
			}
		}
		if offset+pageSize >= filterCount {
			break
		}
		offset = offset + pageSize
	}
	return
}
