package ossUtil

import (
	"errors"
	"fmt"
	"io"
	"onlineCLoud/internel/app/config"
	"os"
	"path"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssClient struct {
	bucket *oss.Bucket
	client *oss.Client
}

func NewClient() (*OssClient, error) {
	C := config.C.Oss
	client, err := oss.New(C.Host, C.OssAccessKeyID, C.OssAccessKeySecret)
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket(C.Bucket)
	if err != nil {
		return nil, err
	}
	return &OssClient{
		bucket: bucket,
		client: client,
	}, nil
}

func (oss *OssClient) CopyDirFromLocal(srcDir, destDir string) error {

	fileInfos, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, file := range fileInfos {
		if file.IsDir() {
			err := oss.CopyDirFromLocal(path.Join(srcDir, file.Name()), path.Join(destDir, file.Name()))
			if err != nil {
				return err
			}
		} else {

			srcFile, err := os.Open(path.Join(srcDir, file.Name()))
			if err != nil {
				return err
			}
			defer srcFile.Close()

			err = oss.bucket.PutObject(path.Join(destDir, file.Name()), srcFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (oss *OssClient) Get(filepath string) (io.ReadCloser, error) {
	rc, err := oss.bucket.GetObject(filepath)
	if err != nil {
		return nil, err
	}

	return rc, nil
}

// DeleteDir 删除指定目录或文件
func (client *OssClient) DeleteDir(dirOrFile string) error {
	isFile := false

	// 检查是文件还是目录
	lor, err := client.bucket.ListObjects(oss.Prefix(dirOrFile), oss.Delimiter("/"))
	if err != nil {
		return err
	}
	if len(lor.Objects) == 0 && len(lor.CommonPrefixes) == 0 {
		// 没有列出任何对象，可能是文件
		isFile = true
		lor, err = client.bucket.ListObjects(oss.Prefix(dirOrFile))
		if err != nil {
			return err
		}
		if len(lor.Objects) == 0 {
			return errors.New("文件或目录不存在")
		}
	}

	if isFile {
		// 删除文件
		err := deleteObjects(client, []string{dirOrFile})
		if err != nil {
			return err
		}
		fmt.Println("文件删除成功")
	} else {
		// 递归删除目录
		err := deleteDirectory(client, dirOrFile)
		if err != nil {
			return err
		}
		fmt.Println("目录及其内容删除成功")
	}

	return nil
}

// deleteObjects 批量删除对象
func deleteObjects(client *OssClient, keys []string) error {
	delRes, err := client.bucket.DeleteObjects(keys)
	if err != nil {
		return err
	}
	if len(delRes.DeletedObjects) > 0 {
		fmt.Println("这些对象未能成功删除:", delRes.DeletedObjects)
		return errors.New("文件删除不完全")
	}
	return nil
}

// deleteDirectory 递归删除目录
func deleteDirectory(client *OssClient, dir string) error {
	marker := oss.Marker("")
	for {
		lor, err := client.bucket.ListObjects(marker, oss.Prefix(dir))
		if err != nil {
			return err
		}

		keys := []string{}
		for _, object := range lor.Objects {
			keys = append(keys, object.Key)
		}

		err = deleteObjects(client, keys)
		if err != nil {
			return err
		}

		if !lor.IsTruncated {
			break
		}
		marker = oss.Marker(lor.NextMarker)
	}
	return nil
}
