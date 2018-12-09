package main

import (
	"fmt"
	"os"
	"time"
)

func countdown(count chan struct{}) {
	tick := time.Tick(1*time.Second)
	for i:=10; i>0; i-- {
		fmt.Printf("countdown: %d\n", i)
		<- tick
	}
	count <- struct{}{}
	return
}

// 中断发射
func abort(ach chan struct{}) {
	os.Stdin.Read(make([]byte, 1))
	ach<- struct{}{}
	return
}

func main() {
	var count = make(chan struct{})
	var ach = make(chan struct{})
	go countdown(count)
	go abort(ach)
	select {

	case <-ach:
		fmt.Println("abort launch...")
	}
	return
}