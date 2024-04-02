package processutil

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// ExecuteCommand 执行外部命令并返回输出结果和可能的错误。
func ExecuteCommand(cmd *exec.Cmd, logOutput bool) (string, error) {
	// 创建一个命令对象

	// 如果需要记录命令输出，则创建一个缓冲区来捕获输出
	var stdout, stderr bytes.Buffer
	if logOutput {
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
	}

	// 执行命令
	err := cmd.Run()

	// 如果执行出错，返回错误信息
	if err != nil {
		return "", fmt.Errorf("执行命令时出错: %s\n%s", err.Error(), stderr.String())
	}

	// 获取命令的标准输出
	output := stdout.String()

	// 返回输出结果
	return strings.TrimSpace(output), nil
}
