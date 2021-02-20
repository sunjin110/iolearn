package main

import (
	"fmt"
	"iolearn/pkg/selectt"
)

func main() {

	fmt.Println("selectt")

	// netcat -u localhost 1235
	// netcat -u localhost 1236
	selectt.UdpServer()

}
