package main

import (
	"fmt"
	"iolearn/pkg/common/chk"
	"iolearn/pkg/common/jsonutil"
	"syscall"
	"time"
)

func main() {

	fmt.Println("検証")

	fmt.Println("sync blocking server")

	// ソケットの作成
	// AF_INET: ip4
	// SOCK_STREAM: tcp
	// 3つ目の引数はprotocol, 1つ目のprotocol familyに順ずるのであれば0でいい
	// しかし、複数のプロトコルが存在してもいい
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM|syscall.O_NONBLOCK, 0)
	chk.SE(err)
	defer func(fd int) {
		err = syscall.Close(fd)
		chk.SE(err)
	}(fd)

	// ソケットの設定
	// Sockaddrを直接は定義させてくれない
	// Sockaddrのinterfaceを満たしているSockaddrInet4を使用する
	addr := &syscall.SockaddrInet4{
		Port: 1234,
		Addr: [4]byte{0, 0, 0, 0},
	}
	syscall.Bind(fd, addr)

	// TCPクライアントからの接続要求をまてる状態にする
	// 2つめの引数は、backlog 保留中の接続キューの最大長をしていする、キューがいっぱいの状態で接続要求がくると
	// Clientがわで「ECONNREFUSED」というエラーを受け取る
	syscall.Listen(fd, 2)

	// 順番に処理していく
	// Readをまったり、書き込みをしている間、カーネルから帰ってこない
	for {

		// non blockしているから、ここでAcceptを待たない感じか
		nfd, sa, err := syscall.Accept(fd)
		chk.SE(err)

		fmt.Println("client is ", jsonutil.Marshal(sa))
		chk.SE(err)

		// clientからのメッセージを待つ
		// ここでは1024byteにしている、clientからのメッセージがこれ以上の場合は、次のReadで続きを受け取る
		// 今回は、ループさせずに取得して終わる
		var buf [1024]byte
		msgLen, err := syscall.Read(nfd, buf[:])
		chk.SE(err)
		fmt.Println("msg is ", string(buf[:msgLen]))

		// clientにメッセージを送る
		_, err = syscall.Write(nfd, buf[:msgLen])
		chk.SE(err)

		time.Sleep(time.Second)

		// 接続の終了
		syscall.Close(nfd)

	}
}
