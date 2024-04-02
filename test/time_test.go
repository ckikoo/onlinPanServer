package test

import (
	"fmt"
	"onlineCLoud/internel/app/define"
	"testing"
	"time"
)

func TestT(t *testing.T) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	now := time.Now()
	fmt.Printf("now.Add(time.Hour): %v\n", now.Add(time.Hour*24*define.FileShare7Day))
}
