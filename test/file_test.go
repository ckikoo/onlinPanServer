package test

import (
	"fmt"
	"io"
	fileUtil "onlineCLoud/pkg/util/file"
	"os"
	"path"
	"path/filepath"
	"testing"
	"time"
)

func TestCreateFile(t *testing.T) {

	f, err := fileUtil.FileCreate("a/b/c/d/e/a.go", 0)

	if err != nil {
		t.Error(err)
	}
	defer f.Close()
}

func TestMergeFile(t *testing.T) {
	tempDir := fmt.Sprintf("../upload/%v/%v/%v", time.Now().Month(), 2, "26debbea-e358-4eba-9d27-981491a3af7b")
	fs, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprintf(path.Join(tempDir, "total"))
	finalFile, err := os.Create(path)

	if err != nil {
		t.Fatal(err)
	}
	defer finalFile.Close()

	for _, f := range fs {
		blockFileName := filepath.Join(tempDir, f.Name())
		blockFile, err := os.Open(blockFileName)
		if err != nil {
			t.Fatal(err)
		}
		defer blockFile.Close()
		_, err = io.Copy(finalFile, blockFile)
		if err != nil {
			t.Fatal(err)
		}
	}

}

func TestFileExt(t *testing.T) {
	filename := "本专科生国家助学金申请表（新表-注意表格格式）.pdf"
	str := fileUtil.GetFileExt(filename)
	fmt.Printf("str: %v\n", str)
}
