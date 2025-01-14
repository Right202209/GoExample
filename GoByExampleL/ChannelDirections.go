package main

import "fmt"

func ping(pings chan<- string, msg string) {
	pings <- msg
}

//只接受发送通道

func pong(pings <-chan string, pongs chan<- string) {
	msg := <-pings
	pongs <- msg
}

//接受 一个接受通道 和 一个发送通道

func main() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)

	ping(pings, "passed message")
	pong(pings, pongs)

	fmt.Println(<-pongs)
}
