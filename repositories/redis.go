// repositories/redis.go
package repositories

import (
	"github.com/go-redis/redis/v8"
	"github.com/programmercintasunnah/go-todolist-ilcs/config"
)

func InitRedis(cfg config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return rdb
}
