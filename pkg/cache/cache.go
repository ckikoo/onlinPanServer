package cache

import (
	"bytes"
	"fmt"
	"io"
	"onlineCLoud/pkg/cache/file/hadoop"
	"onlineCLoud/pkg/cache/file/local"
	"onlineCLoud/pkg/cache/file/oss"
	"onlineCLoud/pkg/errors"
	"onlineCLoud/pkg/file"
)

// CacheReader 定义了一个缓存读取器
type CacheReader struct {
	FilePath string
	Caches   []Cache // 用于配置不同级别的缓存
}

// Cache 接口定义了缓存类型需要实现的方法
type Cache interface {
	Get(filePath string) (*file.AbstractFile, error)
	Put(filePath string, content io.Reader) error
	Truncate(filepath string) error
	Delete(filePath string) error
}

func NewCacheReader(filePath string) *CacheReader {
	cache := make([]Cache, 3)
	cache[0] = &local.LocalCache{}
	cache[1] = &hadoop.HadoopCache{}
	cache[2] = oss.NewOssCache()

	return &CacheReader{
		FilePath: filePath,
		Caches:   cache,
	}
}

func (cr *CacheReader) out(cache Cache, file *file.AbstractFile) error {
	for {
		buf := make([]byte, 16*1024*1024)
		cnt, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("error reading file: %w", err)
		}
		if cnt == 0 {
			break
		}
		err = cache.Put(cr.FilePath, bytes.NewReader(buf[:cnt]))
		if err != nil {
			return err
		}
	}
	return nil
}

// resetCache 从文件中读取内容，并写入低级别缓存中
func (cr *CacheReader) resetCache(file *file.AbstractFile, begin int, limit int) error {

	// 从文件中读取内容，并写入低级别缓存中
	for i, cache := range cr.Caches {
		if i < begin || i >= limit {
			continue
		}
		file.Seek(0, io.SeekStart)
		cr.out(cache, file)
	}
	return fmt.Errorf("cache not found")
}

// Read 从缓存中读取文件内容
func (cr *CacheReader) Read() (*file.AbstractFile, error) {
	for pos, cache := range cr.Caches {
		r, err := cache.Get(cr.FilePath)
		if err != nil {
			fmt.Println(pos)
			continue
		}
		if nil == err {
			if pos == 0 {
				return r, nil
			} else {
				cr.out(cr.Caches[0], r)
				r.Close()
				r, err := cr.Caches[0].Get(cr.FilePath)
				if err != nil {
					panic(err)
				}
				go cr.resetCache(r, 1, pos)
				r.Seeker.Seek(0, io.SeekStart)
				return r, nil
			}
		}
	}

	return nil, errors.New("file not found in any cache")
}

func (cr *CacheReader) Delete(filePath string) error {

	for _, do := range cr.Caches {

		if err := do.Delete(filePath); err != nil {
			// 记录日志
			continue
		}

	}

	return nil
}
