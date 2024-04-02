package test

import (
	"fmt"
	"os/exec"
	"testing"
)

// func TestFileUTil(t *testing.T) {
// 	cmd := "ls -al"

// 	str, err := processutil.ExecuteCommand(cmd, false)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	fmt.Printf("str: %v\n", str)
// }

func TestOSBash(t *testing.T) {
	cmd := exec.Command("ls")

	// 执行命令并等待完成
	output, err := cmd.CombinedOutput()

	if err != nil {
		// 命令执行出错
		fmt.Println("命令执行出错:", err)
		return
	}

	// 输出命令的标准输出
	fmt.Println("命令输出:")
	fmt.Println(string(output))
}
