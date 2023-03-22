package common

import (
	"go-high-currency/config"
	"log"

	redis "github.com/go-redis/redis/v8"
)

func NewRedisClient(config *config.RedisConf) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password, // 没有密码，默认值
		DB:       0,               // 默认DB 0
	})
	log.Println("get rdb success")
	return rdb
}
