package cache

import (
	"context"

	"github.com/alvarolucio2007/GoSocial/internal/store"
)

type MockCacheStore struct{}

func NewMockCache(mockCache MockCacheStore) Storage {
	return Storage{Users: mockCache}
}

func (m MockCacheStore) Get(ctx context.Context, id int) (*store.User, error) {
	return nil, nil
}

func (m MockCacheStore) Set(ctx context.Context, user *store.User) error {
	return nil
}
