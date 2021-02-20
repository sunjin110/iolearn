package main

import (
	"fmt"
	"iolearn/pkg/syncnonblock"
)

func main() {

	fmt.Println("sync non blocking")

	syncnonblock.Server()

}
