package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 原子计数器
func main() {
	var ops atomic.Uint64 //原子整数表示（正）计数器
	var wg sync.WaitGroup //等待所有g完成

	for i := 0; i < 50; i++ { //启动50个
		wg.Add(1)

		go func() {
			for c := 0; c < 1000; c++ {
				ops.Add(1)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println("ops:", ops.Load())
}
