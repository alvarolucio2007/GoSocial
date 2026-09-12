package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/alvarolucio2007/GoSocial/internal/store"
	"github.com/redis/go-redis/v9"
)

const UserExpTime = time.Hour

type UserCache interface {
	Set(ctx context.Context, user *store.User) error
	Get(ctx context.Context, userID int) (*store.User, error)
}
type RedisUserCache struct {
	rdb *redis.Client
}

func (r *RedisUserCache) Set(ctx context.Context, user *store.User) error {
	cacheKey := fmt.Sprintf("user-%v", user.ID)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, cacheKey, json, UserExpTime).Err()
}

func (r *RedisUserCache) Get(ctx context.Context, userID int) (*store.User, error) {
	cacheKey := fmt.Sprintf("user-%d", userID)
	data, err := r.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var user store.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, err
	}
	return &user, nil
}
