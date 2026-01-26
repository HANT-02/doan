package caching

import (
	"context"
	"doan/internal/caching"
)

type redisCacheManager struct {
}

func (r *redisCacheManager) GetString(ctx context.Context, key string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewRedisCacheManager() caching.CacheManager {
	return &redisCacheManager{}
}
