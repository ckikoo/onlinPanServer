package cache

import (
	"context"
	"onlineCLoud/internel/app/dao/redisx"
	"sync"

	"gorm.io/gorm"
)

type CachePack struct {
	DB *gorm.DB
	RD *redisx.Redisx
}

var (
	CachePackKey               = "online:file_cache:%v"
	PackCache       *CachePack = nil
	cache_file_once sync.Once
)

// TODP wait to add cache with cache
func Init(ctx context.Context, db *gorm.DB, rd *redisx.Redisx) {
	cache_file_once.Do(func() {
		PackCache = &CachePack{
			DB: db,
			RD: rd,
		}
		// FileCache.cacheInit(ctx)
	})

}
