package main

import (
	"iolearn/pkg/syncblock"
	"time"
)

func main() {
	go syncblock.Server()

	time.Sleep(time.Second)

	syncblock.Client("client1")
}
