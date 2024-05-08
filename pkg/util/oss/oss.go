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

// 删除当前目录,以及子目录
func (cleint *OssClient) DeleteDir(dir string) error {

	marker := oss.Marker("")
	objectsDeleted := 0
	for {
		lor, err := cleint.bucket.ListObjects(marker, oss.Prefix(dir))
		if err != nil {
			return err
		}

		keys := []string{}
		for _, object := range lor.Objects {
			keys = append(keys, object.Key)
		}

		delRes, err := cleint.bucket.DeleteObjects(keys)
		if err != nil {
			return err
		}

		if len(delRes.DeletedObjects) > 0 {
			fmt.Println("These objects were not deleted successfully:", delRes.DeletedObjects)
			return errors.New("文件删除不完全")
		}

		objectsDeleted += len(keys)

		marker = oss.Marker(lor.NextMarker)
		if !lor.IsTruncated {
			break
		}
	}

	return nil
}
