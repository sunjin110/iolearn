package main

import (
	"fmt"
	"iolearn/pkg/common/chk"
	"syscall"
)

func main() {
	fmt.Println("sync blocking server")

	// ソケットの作成
	// AF_INET: ip4
	// SOCK_STREAM: tcp
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	chk.SE(err)
	// defer syscall.Close(fd)
	defer func(fd int) {
		err = syscall.Close(fd)
		chk.SE(err)
	}(fd)

}
