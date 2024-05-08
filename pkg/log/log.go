package logger

import (
	"fmt"
	"onlineCLoud/internel/app/config"
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
		fmt.Printf("config.C.LOGGER: %v\n", config.C.LOGGER)
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
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

		// 配置日志级别
		atomicLevel := zap.NewAtomicLevel()
		atomicLevel.SetLevel(zap.InfoLevel)

		// 创建日志核心
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderCfg),
			zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(createLogFile(logFilePath))),
			atomicLevel,
		)

		// 创建 Logger 实例
		globalLogger = zap.New(core)
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
func Log(level string, message ...string) {

	str := ""

	for _, v := range message {
		str += "[" + v + "]"
	}

	switch level {
	case "DEBUG":
		globalLogger.Debug(str)
	case "INFO":
		globalLogger.Info(str)
	case "WARN":
		globalLogger.Warn(str)
	case "ERROR":
		globalLogger.Error(str)
	}
}

// Close 关闭日志记录器
func Close() {
	if globalLogger != nil {
		globalLogger.Sync()
	}

}
