package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func filterAndPrintFilenames(directory, keyword string) error {
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil // 如果是目录则跳过
		}

		// 如果文件名以.git开头，则将优先级调高
		if strings.HasPrefix(path, ".git") {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNumber := 1
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, keyword) {
				fmt.Printf("[%s][%s][%d][%s]\n", filepath.Dir(path)+"/"+filepath.Base(path), filepath.Base(path), lineNumber, line)
			}
			lineNumber++
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	directory := "."
	keyword := "fmt.Printf"
	if err := filterAndPrintFilenames(directory, keyword); err != nil {
		fmt.Println("Error:", err)
	}
}
