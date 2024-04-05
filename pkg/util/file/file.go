package fileUtil

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"strings"
)

func NewDir(path string) (string, error) {
	err := os.MkdirAll(path, 0744)
	if err != nil {
		return "", err
	}
	return path, nil
}

func FileCreate(filename string, mode int) (*os.File, error) {
	_, err := NewDir(path.Dir(filename))
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filename, mode, 0755)
	return file, err
}

func FileMerge(tempDir, dest string) error {

	files, err := os.ReadDir(tempDir)
	if err != nil {
		return err
	}
	finalFile, err := FileCreate(dest, os.O_WRONLY|os.O_CREATE)
	if err != nil {
		fmt.Printf("dest: %v\n", dest)
		fmt.Printf("err: %v\n", err)
		return err
	}
	defer finalFile.Close()
	for i := 0; i < len(files); i++ {
		blockFile, err := os.Open(path.Join(tempDir, fmt.Sprintf("%d", i)))
		if err != nil {
			return err
		}
		defer blockFile.Close()
		_, err = io.Copy(finalFile, blockFile)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetFileExt(filename string) string {
	ext := path.Ext(filename)
	if ext != "" {
		ext = ext[1:]
		ext = strings.ToLower(ext) // 将扩展名转换为小写
	}
	return ext
}

func Rename(filename string) string {
	index := strings.LastIndex(filename, ".")
	if index == -1 {
		return filename + fmt.Sprintf("%02d", rand.Int31n(100))
	}

	str := fmt.Sprintf("%v%v%v", filename[:index], fmt.Sprintf("%v", rand.Int31n(100)), filename[index:])

	return str
}
