package main

import (
	"fmt"
	"time"
)

func main() {
	timer1 := time.NewTimer(3 * time.Second)

	<-timer1.C
	fmt.Println("Timer1 1 fired")

	timer2 := time.NewTimer(8 * time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 fired")
	}()
	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stop")
	}
	time.Sleep(2 * time.Second)
}
