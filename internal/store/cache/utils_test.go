package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	goRedis "github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

var (
	testRedis *goRedis.Client
	testCache Storage
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	var testContainer *redis.RedisContainer
	var err error
	testRedis, testContainer, err = setupTestCache()
	if err != nil {
		log.Fatalf("unable to create redis and testcontainers %v", err)
	}
	if testRedis == nil {
		log.Fatalf("testRedis is nil\n")
	}
	if testContainer == nil {
		log.Fatalf("testContainer Redis is nil")
	}
	testCache = NewRedisStorage(testRedis)
	code := m.Run()
	if err := testContainer.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}
	if err := testRedis.Close(); err != nil {
		log.Printf("failed to terminate DB: %v", err)
	}
	os.Exit(code)
}

func setupTestCache() (*goRedis.Client, *redis.RedisContainer, error) {
	ctx := context.Background()
	rdc, err := redis.Run(ctx, "redis:8-alpine")
	if err != nil {
		return nil, nil, err
	}
	connStr, err := rdc.ConnectionString(ctx)
	if err != nil {
		return nil, nil, err
	}
	rdb := New(strings.TrimPrefix(connStr, "redis://"), "", 0)
	if rdb == nil {
		return nil, nil, fmt.Errorf("rdb is nil at setupTestCache")
	}
	return rdb, rdc, nil
}
