package file

import (
	"context"
	"errors"
	"fmt"
	"log"

	"onlineCLoud/internel/app/dao/redisx"
	"onlineCLoud/pkg/util/json"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type CacheFileServer struct {
	DB *gorm.DB
	RD *redisx.Redisx
}

var (
	CachaNoFound                     = errors.New("cache not found")
	CacheIsEmpty                     = errors.New("cache empty")
	CacheFileKey                     = "online:file_cache:%v"
	FileCache       *CacheFileServer = nil
	cache_file_once sync.Once
)

func Init(ctx context.Context, db *gorm.DB, rd *redisx.Redisx) {
	cache_file_once.Do(func() {
		FileCache = &CacheFileServer{
			DB: db,
			RD: rd,
		}
		FileCache.cacheInit(ctx)
	})

}

func (c *CacheFileServer) cacheInit(ctx context.Context) {
	files := make([]File, 0)
	batchSize := 20
	offset := 0
	for {
		files = files[:0]
		if err := c.DB.Raw("SELECT tb_file.* FROM tb_file JOIN (SELECT file_md5, MIN(file_id) AS min_file_id FROM tb_file GROUP BY file_md5) AS subquery ON tb_file.file_md5 = subquery.file_md5 AND tb_file.file_id = subquery.min_file_id LIMIT ? OFFSET ?", batchSize, offset).Find(&files).Error; err != nil {
			panic(err)
		}
		if len(files) == 0 {
			break // 如果没有更多记录可查询，退出循环
		}
		Md5s := make([]string, 0)

		for _, file := range files {
			Md5s = append(Md5s, file.FileMd5)
		}
		if err := c.RD.ZsetWithTimestamps(ctx, fmt.Sprintf(CacheFileKey, "key"), Md5s, 20, time.Hour*24*31); err != nil {
			log.Default().Printf("error load file")
		}

		_, err := c.RD.HMset(ctx, fmt.Sprintf(CacheFileKey, "file"), ToMd5Map(files))
		if err != nil {
			log.Println(err)
		}
		offset += batchSize
		time.Sleep(time.Second * 1)
	}

}

func (c *CacheFileServer) AddFile(ctx context.Context, f File) {
	strs := make([]string, 1)
	strs[0] = f.FileMd5
	c.RD.ZsetWithTimestamps(ctx, fmt.Sprintf(CacheFileKey, "key"), strs, 1, time.Hour*24*31)

	c.RD.HSet(ctx, fmt.Sprintf(CacheFileKey, "file"), f.FileMd5, f)
}

func (c *CacheFileServer) FindMd5(ctx context.Context, md5str string) (*File, error) {
	m, err := c.RD.HGet(ctx, fmt.Sprintf(CacheFileKey, "file"), md5str)
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err == redis.Nil {
		return nil, CachaNoFound
	}
	if m == "" {
		return nil, CacheIsEmpty
	}

	f := new(File)
	err = json.Unmarshal([]byte(m), f)
	if err != nil {
		log.Default().Println(err)
		return nil, err
	}

	return f, nil
}

func (c *CacheFileServer) DeleteFile(ctx context.Context, md5str string) error {
	_, err := c.RD.Zrem(ctx, fmt.Sprintf(CacheFileKey, "file"), []string{md5str})
	if err != nil && err != redis.Nil {
		return err
	}

	_, err = c.RD.HDel(ctx, fmt.Sprintf(CacheFileKey, "file"), []string{md5str})
	if err != nil {
		log.Default().Println(err)
		return err
	}

	return nil
}
