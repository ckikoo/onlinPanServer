package oss

import (
	"io"
	"onlineCLoud/pkg/file"
	ossUtil "onlineCLoud/pkg/util/oss"
	"os"
)

// LocalCache 本地缓存
type OssCache struct {
	client *ossUtil.OssClient
}

func NewOssCache() *OssCache {
	c, _ := ossUtil.NewClient()
	return &OssCache{
		client: c,
	}
}
func (oc *OssCache) Get(filePath string) (*file.AbstractFile, error) {
	file := new(file.AbstractFile)
	f, err := oc.client.Get(filePath)
	if err != nil {
		return nil, err
	}
	file.Closer = f
	file.Reader = f

	return file, nil
}

func (oc *OssCache) Put(filePath string, content io.Reader) error {
	// 打开文件进行追加写入，如果文件不存在则会创建该文件

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将 content 的数据写入到文件中
	_, err = io.Copy(file, content)
	if err != nil {
		return err
	}

	return nil
}

func (oc *OssCache) Truncate(filepath string) error {

	return nil
}

func (oc *OssCache) Delete(filePath string) error {
	return oc.client.DeleteDir(filePath)
}
