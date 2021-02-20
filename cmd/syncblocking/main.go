package main

import (
	"iolearn/pkg/syncblock"
	"time"
)

func main() {

	// 任意のclientを送り終えたタイミングでServerを正常終了させるchannelの仕組みをいれる

	go syncblock.Server()

	time.Sleep(time.Second) // serverが立ち上がるのをまつ

	go syncblock.Client("client1")
	go syncblock.Client("client2")

	time.Sleep(3 * time.Second)
}
