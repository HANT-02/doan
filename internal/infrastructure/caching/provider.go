package caching

import "github.com/google/wire"

var CacheManagerProvider = wire.NewSet(NewRedisCacheManager)
