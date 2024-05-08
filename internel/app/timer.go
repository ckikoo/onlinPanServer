package app

import (
	"fmt"
	"onlineCLoud/internel/app/define"
	"onlineCLoud/pkg/cache/file/hadoop"
	"onlineCLoud/pkg/cache/file/local"
	"onlineCLoud/pkg/timer"
	"os"
	"path/filepath"
	"time"
)

func cleanLocalExpCache() {
	timer := timer.GetInstance()
	lc := &local.LocalCache{}

	infos, err := lc.ReadDir("./upload")
	if err != nil {
		fmt.Printf("failed to read directory: %v\n", err)
		return
	}

	for _, file := range infos {
		info, err := file.Info()
		if err != nil {
			fmt.Printf("failed to get file info: %v\n", err)
			continue
		}

		exp := info.ModTime().Add(time.Hour * 24)
		timer.Add(fmt.Sprintf(define.LocalCacheTimerKey, info.Name()), exp, func() {
			if err := os.RemoveAll(filepath.Join("./upload", info.Name())); err != nil {
				fmt.Printf("timer:del_localfile:%v err: %v\n", info.Name(), err)
			}
		})

	}
}
func cleanHadoopExpCache() {
	timer := timer.GetInstance()
	lc := &hadoop.HadoopCache{}

	infos, err := lc.ReadDir("/")
	if err != nil {
		fmt.Printf("failed to read directory: %v\n", err)
		return
	}

	for _, file := range infos {
		exp := file.ModTime().Add(time.Hour * 24 * 7)

		timer.Add(fmt.Sprintf(define.HadoopCacheTimerKey, file.Name()), exp, func() {
			fmt.Println("111111")
			lc.Delete(file.Name())
		})

	}
}

// InitTimer 函数初始化定时器，用于定期清理本地和Hadoop的过期缓存。
// 该函数不接受参数，也不返回任何值。
func InitTimer() {
	// 清理本地经验缓存
	cleanLocalExpCache()
	// 清理Hadoop经验缓存
	cleanHadoopExpCache()
}
