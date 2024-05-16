package api

import (
	"fmt"
	"io"
	"net/http"
	"onlineCLoud/internel/app/config"
	"onlineCLoud/internel/app/ginx"
	"onlineCLoud/internel/app/service"
	"onlineCLoud/pkg/cache"
	"onlineCLoud/pkg/contextx"
	"onlineCLoud/pkg/file"
	"onlineCLoud/pkg/timer"
	hdfsUtil "onlineCLoud/pkg/util/hdfs"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type DownLoadApi struct {
	Srv     *service.DownLoadSrv
	FileSrv *service.FileSrv
}

func (d DownLoadApi) Download(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	if len(code) == 0 {
		ginx.ResFail(c)
		return
	}

	fmt.Printf("code: %v\n", code)
	path, err := d.Srv.FindDownloadByCode(ctx, code)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}
	cr := cache.NewCacheReader(path)
	reader, err := cr.Read()
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}
	defer reader.Close()

	info, err := os.Stat(path)
	if err != nil {
		ginx.ResFailWithMessage(c, err.Error())
		return
	}

	timer.GetInstance().Add("local_cache_delete"+path, time.Now().Add(time.Hour*24), func() {
		d.Srv.Delete(ctx, code)
		os.RemoveAll(path)
	})
	timer.GetInstance().Add("hdfs_cache_delete"+path, time.Now().Add(time.Hour*24*7), func() {
		client, err := hdfsUtil.NewClient(config.C.Hadoop.Host)
		if err != nil {
			panic(err)
		}
		client.DeleteFile("/" + path)
	})

	limit, _ := d.Srv.GetDownLoadSpeed(ctx, contextx.FromUserID(ctx))
	if limit <= 100 {
		limit = 100
	}

	fmt.Printf("limit: %v\n", limit)

	// 限制每次读取的字节数为 100 KB
	limitedReader := &RateLimitedReader{
		R:     reader,
		Limit: int64(limit) * 1024,
	}
	c.Header("x-cookie", "1")

	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", info.Name()))
	c.Header("Accept-Ranges", "bytes")
	c.Header("Etag", info.ModTime().String())
	http.ServeContent(c.Writer, c.Request, info.Name(), info.ModTime(), limitedReader)
}

// RateLimitedReader 实现了 io.Reader 接口，用于限制读取速度。
type RateLimitedReader struct {
	R       io.Reader // 1原始的 io.Reader
	Limit   int64     // 每秒读取的字节数限制
	LastSec int64     // 上次读取的时间戳
	ReadCnt int64     // 当前秒内已读取的字节数
}

func (r *RateLimitedReader) Seek(offset int64, whence int) (int64, error) {
	return r.R.(*file.AbstractFile).Seek(offset, whence)
}

func (r *RateLimitedReader) Read(p []byte) (n int, err error) {
	// 获取当前时间戳
	now := time.Now().Unix()

	// 如果距离上次读取的时间大于1秒，则重置读取字节数
	if now-r.LastSec >= 1 {
		r.ReadCnt = 0
		r.LastSec = now
	}

	// 计算当前秒内还可以读取的字节数
	remaining := r.Limit - r.ReadCnt
	// 如果剩余可读字节数为0，则等待1秒后重新计算
	if remaining <= 0 {
		return 0, nil
	}

	// 限制每次读取的字节数不超过剩余可读字节数
	if int64(len(p)) > remaining {
		p = p[:remaining]
	}

	// 读取数据
	n, err = r.R.Read(p)
	if err != nil {
		return n, err
	}

	// 更新已读取字节数
	r.ReadCnt += int64(n)
	return n, nil
}
