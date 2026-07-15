package main

import (
	"fmt"
	"time"

	"github.com/Austinhamilton1/qlip/qlipboard"
)

func main() {
	qlipboard.Init()

	quit := make(chan bool)
	update := qlipboard.Watch(quit)

	go func() {
		time.Sleep(10 * time.Second)
		close(quit)
	}()

	for data := range update {
		fmt.Println(data)
	}

	fmt.Println("Watcher exited")
}
