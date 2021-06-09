package redis

import (
	"cacheflusher/config/apollo"
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"sync"
	"time"
)

var rc *redis.Client

var once sync.Once

func GetRedisInstance() *redis.Client {
	once.Do(func() {
		client := redis.NewClient(&redis.Options{
			//Addr:       "127.0.0.1:6379",
			Addr:       apollo.ApolloConfig.REDIS_ADDR,
			//Password:   apollo.ApolloConfig.REDIS_PWD,
			Password:   "",
			DB:         apollo.ApolloConfig.REDIS_DB,
			MaxRetries: 1,
			PoolSize:   200,
			DialTimeout : time.Second * 60,
			ReadTimeout : time.Second * 60,
			WriteTimeout: time.Second * 60,
		})
		_, rErr := client.Ping(context.TODO()).Result()
		if rErr != nil {
			log.Printf("redis链接失败,链接字符串: [%v], error: [%v]", apollo.ApolloConfig.REDIS_ADDR, rErr)
			panic("数据库类型 [redis] 链接失败\n")
		}
		rc = client
	})
	return rc
}