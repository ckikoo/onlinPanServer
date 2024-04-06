package main

import (
	"fmt"
	"math/rand"
	"onlineCLoud/pkg/timer"
	"sync"
	"sync/atomic"
	"time"
)

func main() {

	timer, clenaFunc := timer.NewTimerManager()
	defer clenaFunc()

	var count int32
	var wg sync.WaitGroup
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		// go func() {

		timer.Add(i, time.Now().Add(time.Second), func() {
			atomic.AddInt32(&count, 1)
			wg.Done()
		})

		// }()
	}
	for i := 0; i < 10; i++ {
		timer.Del(rand.Int31() % 10000)
	}
	wg.Wait()

	fmt.Println(count)

}
