package database

import (
	"context"
	"deca-task/internal/config"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var redisClent *redis.Client

func InitRedis(cfg *config.Config) *redis.Client {

	redisClent = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	_, err := redisClent.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	log.Println("✅ Redis connected")
	return redisClent
}
