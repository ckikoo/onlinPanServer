package hadoop

import (
	"io"
	"onlineCLoud/pkg/file"
	hdfsUtil "onlineCLoud/pkg/util/hdfs"
)

// HadoopCache Hadoop 缓存
type HadoopCache struct{}

func (hc *HadoopCache) Get(filePath string) (*file.AbstractFile, error) {
	if filePath[0] != '/' {
		filePath = "/" + filePath
	}

	fileAbs := new(file.AbstractFile)

	client, err := hdfsUtil.NewClient("172.20.0.2:9000")
	if err != nil {
		return nil, err
	}
	reader, err := client.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	fileAbs.Closer = reader
	fileAbs.Reader = reader
	fileAbs.Closer = reader
	return fileAbs, nil
}

func (hc *HadoopCache) Put(filePath string, content io.Reader) error {
	if filePath[0] != '/' {
		filePath = "/" + filePath
	}
	client, err := hdfsUtil.NewClient("172.20.0.2:9000")
	if err != nil {
		return err
	}
	return client.WriteFile(filePath, content)

}

func (hc *HadoopCache) Truncate(filepath string) error {
	client, err := hdfsUtil.NewClient("172.20.0.2:9000")
	if err != nil {
		return err
	}

	if filepath[0] != '/' {
		filepath = "/" + filepath
	}

	return client.DeleteFile(filepath)
}

func (hc *HadoopCache) Delete(filepath string) error {
	client, err := hdfsUtil.NewClient("172.20.0.2:9000")
	if err != nil {
		return err
	}

	if filepath[0] != '/' {
		filepath = "/" + filepath
	}

	return client.DeleteFile(filepath)
}
