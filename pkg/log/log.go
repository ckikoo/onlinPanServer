package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

// InitLogger 初始化全局日志记录器
func InitLogger() {
	once.Do(func() {
		dir := "./log"
		// 创建日志目录
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			panic(err)
		}

		// 构建日志文件路径
		logFilePath := filepath.Join(dir, "app.log")

		// 配置日志编码器
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.TimeKey = "ts" // 保留时间键，避免与 zapcore.DefaultEncoderConfig 冲突
		encoderCfg.LevelKey = "level"
		encoderCfg.NameKey = "logger"
		encoderCfg.CallerKey = "caller" // 添加 caller 键，用于记录调用位置（文件名+行号）
		encoderCfg.MessageKey = "msg"
		encoderCfg.StacktraceKey = "stacktrace"

		atomicLevel := zap.NewAtomicLevel()
		atomicLevel.SetLevel(zap.InfoLevel)
		// 创建日志核心
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()), // 使用 ConsoleEncoder
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(createLogFile(logFilePath))),
			atomicLevel, // 直接传递 atomicLevel，它实现了 LevelEnabler 接口
		)

		// 创建 Logger 实例
		globalLogger = zap.New(core, zap.AddCaller())
	})
}

// createLogFile 创建日志文件
func createLogFile(logFilePath string) *os.File {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return file
}

// Log 记录日志
func Log(level string, message ...interface{}) {

	fmt.Printf("globalLogger: %v\n", globalLogger)
	str := ""

	for _, v := range message {
		str += fmt.Sprintf("[ %v ]", v)
	}
	fmt.Printf("str: %v\n", str)
	switch level {
	case "DEBUG":
		globalLogger.Debug(str)
	case "INFO":
		globalLogger.Info(str)
	case "WARN":
		globalLogger.Warn(str)
	case "ERROR":
		globalLogger.Error(str)
	default:
		globalLogger.Info(str)
	}
}

// Close 关闭日志记录器
func Close() {
	if globalLogger != nil {
		globalLogger.Sync()
	}

}
