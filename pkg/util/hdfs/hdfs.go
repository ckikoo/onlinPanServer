package hdfsUtil

import (
	"io"
	fileUtil "onlineCLoud/pkg/util/file"
	"os"
	"path"

	"github.com/colinmarc/hdfs/v2"
)

type HdfsClient struct {
	client *hdfs.Client
}

func NewClient(address string) (*HdfsClient, error) {
	client, err := hdfs.New(address)
	if err != nil {
		return nil, err
	}
	return &HdfsClient{
		client: client,
	}, nil
}

func (hdfs *HdfsClient) Exists(path string) (bool, error) {
	_, err := hdfs.client.Stat(path)
	if err == nil {
		// 文件或目录存在
		return true, nil
	} else if os.IsNotExist(err) {
		// 文件或目录不存在
		return false, nil
	} else {
		// 其他错误
		return false, err
	}
}

func (hdfs *HdfsClient) OpenFile(path string) (*hdfs.FileReader, error) {
	reader, err := hdfs.client.Open(path)
	if err != nil {
		return nil, err
	}
	return reader, nil
}
func (c *HdfsClient) WriteFile(path string, content io.Reader) error {
	// 检查文件是否存在
	exists, err := c.Exists(path)
	if err != nil {
		return err
	}

	var writer *hdfs.FileWriter
	if !exists {
		// 文件不存在，创建新文件
		writer, err = c.client.Create(path)
		if err != nil {
			return err
		}
	} else {
		// 文件存在，追加写入
		writer, err = c.client.Append(path)
		if err != nil {
			return err
		}
	}

	// 将内容写入文件
	_, err = io.Copy(writer, content)
	if err != nil {
		return err
	}

	// 关闭文件
	err = writer.Close()
	if err != nil {
		return err
	}

	return nil
}

func (hdfs *HdfsClient) NewDir(path string) error {
	err := hdfs.client.MkdirAll(path, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (hdfs *HdfsClient) CreateFile(filename string) (*hdfs.FileWriter, error) {
	err := hdfs.NewDir(path.Dir(filename))
	if err != nil {
		return nil, err
	}
	writer, err := hdfs.client.Create(filename)
	if err != nil {
		return nil, err
	}
	return writer, nil
}

func (hdfs *HdfsClient) CopyDirFromLocal(srcDir, destDir string) error {
	err := hdfs.NewDir(destDir)
	if err != nil {
		return err
	}

	fileInfos, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, file := range fileInfos {
		if file.IsDir() {
			err := hdfs.CopyDirFromLocal(path.Join(srcDir, file.Name()), path.Join(destDir, file.Name()))
			if err != nil {
				return err
			}
		} else {
			out, err := hdfs.CreateFile(path.Join(destDir, file.Name()))
			if err != nil {
				return err
			}
			defer out.Close()

			srcFile, err := os.Open(path.Join(srcDir, file.Name()))
			if err != nil {
				return err
			}
			defer srcFile.Close()

			_, err = io.Copy(out, srcFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (hdfs *HdfsClient) CopyFileFromRemote(src, dest string) error {
	file, err := fileUtil.FileCreate(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC)
	if err != nil {
		return err
	}
	defer file.Close()

	read, err := hdfs.client.Open(src)
	if err != nil {
		return err
	}
	defer read.Close()

	_, err = io.Copy(file, read)
	return err
}

func (hdfs *HdfsClient) AppendFile(filename string, content string) error {
	appentWriter, err := hdfs.client.Append(filename)
	if err != nil {
		return err
	}
	defer appentWriter.Close()

	_, err = appentWriter.Write([]byte(content))
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (hdfs *HdfsClient) DeleteFile(path string) error {
	err := hdfs.client.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}
func (hdfs *HdfsClient) GetFileInfo(path string) (os.FileInfo, error) {
	info, err := hdfs.client.Stat(path)
	if err != nil {
		return info, err
	}

	return info, nil
}

func (hdfs *HdfsClient) ReadDir(path string) ([]os.FileInfo, error) {
	files, err := hdfs.client.ReadDir(path)
	if err != nil {
		return nil, err
	}

	return files, nil
}
