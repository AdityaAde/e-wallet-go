package component

import (
	"context"
	"log"
	"time"

	"adityaad.id/belajar-auth/domain"
	"github.com/allegro/bigcache/v3"
)

func GetCacheConnection() domain.CacheRepository {
	cache, err := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	if err != nil {
		log.Fatalf("Error connecting to cache:", err.Error())
	}

	return cache
}
