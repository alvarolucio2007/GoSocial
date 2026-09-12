package cache

import (
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Users UserCache
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &RedisUserCache{rdb: rdb},
	}
}
