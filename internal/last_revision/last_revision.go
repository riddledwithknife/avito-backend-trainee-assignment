package last_revision

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"avito-backend-trainee-assignment/internal/models"
)

func FetchDataFromDBAndCache(db *gorm.DB, redisClient *redis.Client) {
	var banners []models.Banner
	if err := db.Find(&banners).Error; err != nil {
		return
	}

	ctx := context.Background()

	for _, banner := range banners {
		bannerJSON, err := json.Marshal(banner)
		if err != nil {
			fmt.Println("Failed to marshal banner JSON:", err)
			continue
		}

		for _, tagID := range banner.TagIDs {
			if err := redisClient.Set(ctx, fmt.Sprintf("banner:%d:%d", banner.FeatureID, tagID), string(bannerJSON), 5*time.Minute).Err(); err != nil {
				fmt.Println("Failed to set banner data in Redis:", err)
			}
		}
	}
}

func PeriodicDataUpdate(db *gorm.DB, redisClient *redis.Client) {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	FetchDataFromDBAndCache(db, redisClient)

	for {
		select {
		case <-ticker.C:
			FetchDataFromDBAndCache(db, redisClient)
		}
	}
}
