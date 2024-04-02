package test

import (
	"fmt"
	"onlineCLoud/internel/app/service"
	"testing"
)

func TestFFmpeg(t *testing.T) {
	fileID := "1111"
	filePath := "/home/go/src/onlinePanServer/test/video/videoplayback.mp4"

	t.Log(fileID, filePath)
	errL := service.CutFile4Video(fileID, filePath)
	if errL != nil {
		t.Fatal(errL)
	}
}

func TestCut(t *testing.T) {
	inputFilePath := "video/videoplayback.mp4"
	err := service.CutFile4Video("0", inputFilePath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCover(t *testing.T) {
	inputFilePath := "video/videoplayback.mp4"
	err := service.CreateCover4Video(inputFilePath, 150, "video/videoplayback/video.png")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		t.Fatal(err)
	}
}

// func TestNoUsed(t *testing.T) {
// 	inputFilePath := "/home/go/src/onlinePanServer/test/video/videoplayback.mp4"
// 	outputDirPath := "/home/go/src/onlinePanServer/test/video/videoplayback"
// 	outputPath := "/home/go/src/onlinePanServer/test/video/videoplayback/index.ts"

// 	// 检查输出目录是否存在，如果不存在则创建
// 	if _, err := os.Stat(outputDirPath); os.IsNotExist(err) {
// 		err := os.MkdirAll(outputDirPath, os.ModePerm)
// 		if err != nil {
// 			fmt.Printf("无法创建输出目录：%s\n", err)
// 			return
// 		}
// 	}

// 	// 使用 FFmpeg 将视频文件转换为输出文件
// 	// cmd := exec.Command("ffmpeg", "-y", "-i", inputFilePath, "-vcodec", "copy", "-acodec", "copy", "-bsf:v", "h264_mp4toannexb", outputPath)

// 	// processutil.ExecuteCommand(cmd, true)
// 	// 执行命令并等待完成

// 	err := cmd.Run()

// 	if err != nil {
// 		// 命令执行出错
// 		fmt.Printf("FFmpeg 命令执行出错：%s\n", err)
// 		return
// 	}

// 	fmt.Println("转换完成！")
// }

func TestNoUsed(t *testing.T) {
	inputFilePath := "/home/go/src/onlinePanServer/test/video/videoplayback.mp4"

	service.CutFile4Video("11111", inputFilePath)
}
