package test

import (
	"fmt"
	fileUtil "onlineCLoud/pkg/util/file"
	"testing"
)

func TestSr(t *testing.T) {
	str := "596bbed0-474d-40e5-8929-3a250dc7248f"
	fmt.Printf("len(str): %v\n", len(str))
}
func TestStr2(t *testing.T) {
	str := ".txt"
	fmt.Printf("fileUtil.Rename(str): %v\n", fileUtil.Rename(str))

}

func TestStr3(t *testing.T) {
	fmt.Println(1024 * 1024 * 1024 * 100)

}
