package main

import (
	"fmt"
	"time"
)

func main() {
	requsets := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requsets <- i
	}
	close(requsets)

	//Limiter channel will receive a value every 200 milliseconds, A regulator for rate limiter
	limiter := time.Tick(200 * time.Millisecond)

	for req := range requsets {
		<-limiter
		fmt.Println("requests", req, time.Now())
	}
	burstyLimiter := make(chan time.Time, 3)

	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	burstyRequsets := make(chan int, 6)
	for i := 1; i <= 6; i++ {
		burstyRequsets <- i
	}
	close(burstyRequsets)

	for req := range burstyRequsets {
		<-burstyLimiter
		fmt.Println("requests", req, time.Now())
	}
}
