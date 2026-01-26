package caching

import "context"

type CacheManager interface {
	GetString(ctx context.Context, key string) (string, error)
}
