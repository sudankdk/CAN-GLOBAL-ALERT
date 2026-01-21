package my_redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sudankdk/cga/configs"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis(cfg configs.RedisConfig) {
	redisHost := cfg.Host
	redisPort := cfg.Port
	RedisClient = redis.NewClient(&redis.Options{Addr: fmt.Sprintf("%s:%s", redisHost, redisPort)})
}
