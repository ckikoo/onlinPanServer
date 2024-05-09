package test

import (
	logger "onlineCLoud/pkg/log"
	"testing"
)

func TestLog(t *testing.T) {
	logger.InitLogger()
	logger.Log("INFO", "111")
}
