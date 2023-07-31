package initialize

import (
	"fmt"

	"github.com/go-redis/redis"
)

var redisDb *redis.Client

// 初始化日志
func Init_Redis(cfg *RedisConfig) (err error) {
	redisDb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		//Password: fmt.Sprintf("%s", cfg.Password), // no password set
		DB:       cfg.Db, // use default DB
		PoolSize: cfg.PoolSize,
	})

	_, err = redisDb.Ping().Result()
	return err
}

func Close_Redis() {
	_ = redisDb.Close()
}
