package cache

import (
	"github.com/go-redis/redis/v8"
	"learn_together/init"
)

var Rdb *redis.Client

func InitRedis(config *init.Config) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.PWD, // no password set
		DB:       config.Redis.DB,  // use default DB
	})
}
