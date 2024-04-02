package random

import (
	"fmt"
	"testing"
)

func TestRandom(t *testing.T) {
	for i := 0; i < 100000; i++ {
		fmt.Printf("GetRandom(4): %v\n", GetRandom(6))
	}
}
