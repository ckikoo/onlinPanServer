package hadoop

import (
	"fmt"
	"io"
	"onlineCLoud/pkg/file"
	hdfsUtil "onlineCLoud/pkg/util/hdfs"
)

// HadoopCache Hadoop 缓存
type HadoopCache struct{}

func (hc *HadoopCache) Get(filePath string) (*file.AbstractFile, error) {
	fileAbs := new(file.AbstractFile)

	client, err := hdfsUtil.NewClient("172.20.0.2:9000")
	if err != nil {
		return nil, err
	}
	fmt.Printf("filePath: %v\n", filePath)
	reader, err := client.OpenFile("/" + filePath)
	if err != nil {
		return nil, err
	}
	fmt.Printf("reader: %v\n", reader)
	fileAbs.Closer = reader
	fileAbs.Reader = reader
	fileAbs.Closer = reader
	return fileAbs, nil
}

func (hc *HadoopCache) Put(filePath string, content io.Reader) error {
	client, err := hdfsUtil.NewClient("172.20.0.2:9000")
	if err != nil {
		return err
	}
	fmt.Printf("filePath: %v\n", filePath)
	return client.WriteFile(filePath, content)

}

func (hc *HadoopCache) Truncate(filepath string) error {
	client, err := hdfsUtil.NewClient("172.20.0.2:9000")
	if err != nil {
		return err
	}
	return client.DeleteFile(filepath)
}
