package local

import (
	"io"
	"onlineCLoud/pkg/file"
	fileUtil "onlineCLoud/pkg/util/file"
	"os"
)

// LocalCache 本地缓存
type LocalCache struct{}

func (lc *LocalCache) Get(filePath string) (*file.AbstractFile, error) {
	fileabs := new(file.AbstractFile)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	fileabs.Closer = f
	fileabs.Reader = f
	fileabs.Seeker = f
	return fileabs, nil
}

func (lc *LocalCache) Put(filePath string, content io.Reader) error {
	// 打开文件进行追加写入，如果文件不存在则会创建该文件
	file, err := fileUtil.FileCreate(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND)
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

func (lc *LocalCache) Truncate(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}
func (lc *LocalCache) Delete(filepath string) error {
	err := os.RemoveAll(filepath)
	if err != nil {
		return err
	}

	return nil
}
