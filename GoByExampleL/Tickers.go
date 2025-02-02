package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(400 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for true {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("tick at", t)
			}
		}
	}()
	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}
