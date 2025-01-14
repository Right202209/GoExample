package main

// 通道同步

import (
	"fmt"
	"time"
)

func worker(done chan bool) {
	fmt.Print("working...")
	time.Sleep(time.Second)
	fmt.Println("done")

	done <- true
}

func main() {
	done := make(chan bool, 1)

	go worker(done)
	<-done
}

// 如果从此程序中删除 <- done 行，则程序将在 worker 启动之前退出。
