package main

import (
	"iolearn/pkg/syncblock"
	"time"
)

func main() {
	go syncblock.Server()

	// time.Sleep(time.Second)

	go syncblock.Client("client1")
	go syncblock.Client("client2")

	time.Sleep(3 * time.Second)
}
