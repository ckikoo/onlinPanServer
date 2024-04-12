package ossUtil

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssClient struct {
	bucket *oss.Bucket
	client *oss.Client
}

func NewClient() (*OssClient, error) {
	client, err := oss.New("oss-cn-fuzhou.aliyuncs.com", "LTAI5t8ApTY8CGGRSkhazaSb", "gOo3v7a7H5W6CmmoAo3UxisBdwK5LK")
	if err != nil {
		return nil, err
	}
	bucket, err := client.Bucket("online-pan")
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

	fmt.Printf("srcDir: %v\n", srcDir)
	fmt.Printf("destDir: %v\n", destDir)

	for _, file := range fileInfos {
		fmt.Printf("file: %v\n", file)
		if file.IsDir() {
			err := oss.CopyDirFromLocal(path.Join(srcDir, file.Name()), path.Join(destDir, file.Name()))
			if err != nil {
				return err
			}
		} else {

			srcFile, err := os.Open(path.Join(srcDir, file.Name()))
			if err != nil {
				panic(err)
				return err
			}
			defer srcFile.Close()

			err = oss.bucket.PutObject(path.Join(destDir, file.Name()), srcFile)
			if err != nil {
				panic(err)
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
