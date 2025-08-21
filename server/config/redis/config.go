package redis

import (
	"fmt"

	"server/config/env"
	"server/config/log"
	"server/internal/utils/data"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func SetupRedisDatabase(cfg env.Redis) {
	var db int
	if env.Cfg.Server.Mode == data.DEVELOPMENT_MODE {
		db = 1
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", cfg.RHost, cfg.RPort),
		DB:   db,
	})

	_, err := rdb.Ping(rdb.Context()).Result()
	if err != nil {
		log.Log.Fatalf("Gagal terhubung ke Redis: %v", err)
	}

	RDB = rdb
}
