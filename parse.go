package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// printFileContent 递归地遍历目录，并打印每个文件的内容
func printFileContent(dir string) error {
	// 打开目录
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	// 遍历目录中的每个条目
	for _, entry := range dirEntries {
		entryPath := filepath.Join(dir, entry.Name())

		// 如果是目录，则递归地处理子目录
		if entry.IsDir() {
			err := printFileContent(entryPath)
			if err != nil {
				return err
			}
		} else {
			// 如果是文件，则打印文件内容
			err := printFile(entryPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func printFile(path string) error {
	if !strings.HasSuffix(path, ".go") {
		return nil
	}
	key := "fmt.Printf"
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lineNumber := 0
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		lineNumber++

		if strings.Contains(line, key) {
			fmt.Printf("%v:%v:%v", path, lineNumber, line)
		}
	}

	return nil
}

func main() {
	// 指定要递归遍历的目录
	dir := "."

	// 调用 printFileContent 函数进行递归遍历并打印文件内容
	err := printFileContent(dir)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
