package cache

import (
	"errors"

	"github.com/redis/go-redis/v9"
)

var (
	ErrSet = errors.New("unable to set redis value")
	ErrGet = errors.New("unable to get redis value")
)

func New(address, password string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
		Protocol: 2,
	})
	return rdb
}
