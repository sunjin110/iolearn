package selectt

import (
	"fmt"
	"iolearn/pkg/common/chk"
	"iolearn/pkg/common/jsonutil"
	"log"
	"syscall"
)

// TcpServer .
// selectで2つのfdを同時に1threadで管理することができる
func TcpServer() {

	fd1, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	chk.SE(err)
	fd2, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	chk.SE(err)

	// fd closeは省略、めんどくさい

	addr1 := &syscall.SockaddrInet4{
		Port: 1237,
		Addr: [4]byte{0, 0, 0, 0},
	}

	addr2 := &syscall.SockaddrInet4{
		Port: 1238,
		Addr: [4]byte{0, 0, 0, 0},
	}

	syscall.Bind(fd1, addr1)
	syscall.Bind(fd2, addr2)
	syscall.Listen(fd1, 10)
	syscall.Listen(fd2, 10)

	readFds := &syscall.FdSet{}

	// selectでwaitするsocketとしてfdたちを登録する
	FD_SET(readFds, fd1)
	FD_SET(readFds, fd2)

	// maxFds
	maxFd := fd1
	if maxFd < fd2 {
		maxFd = fd2
	}

	fds := &syscall.FdSet{}
	for {

		// fdの初期化
		copy(fds.Bits[:], readFds.Bits[:])

		// この timeout 時間はシステムクロックの粒度に切り上げられ、 カーネルのスケジューリング遅延により少しだけ長くなる可能性がある点に注意すること。 timeval 構造体の両方のフィールドが 0 の場合、 select() はすぐに復帰する (この機能はポーリング (polling) を行うのに便利である)。 timeout に NULL (タイムアウトなし) が指定されると、 select() は無期限に停止 (block) する。
		// 書き込みfdも別で監視できるが、echoの場合はそのまま返せばいいから問題ない?
		log.Println("ここでとまってる?")
		log.Println("before fds is ", jsonutil.Marshal(fds))
		_, err = syscall.Select(maxFd+1, fds, nil, nil, nil)
		chk.SE(err)

		log.Println("fds is ", jsonutil.Marshal(fds))

		if FD_ISSET(fds, fd1) {

			// accept
			nfd1, sa, err := syscall.Accept(fd1)
			chk.SE(err)
			log.Println("sa1 is ", jsonutil.Marshal(sa))

			var buf [1024]byte
			// read
			i, err := syscall.Read(nfd1, buf[:])
			chk.SE(err)
			fmt.Println("fd1 msg is ", string(buf[:i]))
			// write
			syscall.Write(nfd1, buf[:i])
		}

		if FD_ISSET(fds, fd2) {

			// accept
			syscall.Accept(fd2)

			var buf [1024]byte
			// read
			i, err := syscall.Read(fd2, buf[:])
			chk.SE(err)
			fmt.Println("fd2 msg is ", string(buf[:i]))
			// write
			syscall.Write(fd2, buf[:i])
		}

	}

}

// netcat -u localhost 1235
// netcat -u localhost 1236

// UdpServer .
func UdpServer() {

	// TODO まず普通にselect system callを使用したi/oを実装してみる

	// TODO その後、O_NONBLCOKを指定して同じことをしてみる

	// 受信側

	// ソケットの作成
	fd1, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	chk.SE(err)
	fd2, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	chk.SE(err)

	defer func(fd1, fd2 int) {
		err = syscall.Close(fd1)
		chk.SE(err)
		err = syscall.Close(fd2)
		chk.SE(err)
	}(fd1, fd2)

	// ソケットの設定
	addr1 := &syscall.SockaddrInet4{
		Port: 1235,
		Addr: [4]byte{0, 0, 0, 0},
	}

	addr2 := &syscall.SockaddrInet4{
		Port: 1236,
		Addr: [4]byte{0, 0, 0, 0},
	}

	syscall.Bind(fd1, addr1)
	syscall.Bind(fd2, addr2)

	readFds := &syscall.FdSet{}

	// selectでまつ読み込みソケットとしてfdたちを登録
	FD_SET(readFds, fd1)
	FD_SET(readFds, fd2)

	log.Println("fd1", fd1, "fd2", fd2)
	log.Println("readFds is ", jsonutil.Marshal(readFds))
	// syscall.FD_ZERO()

	// fdSet := syscall.FdSet{}
	// syscall.Select(nfd int, r *syscall.FdSet, w *syscall.FdSet, e *syscall.FdSet, timeout *syscall.Timeval)

	// selectで監視するfdの最大値を取得する
	maxFd := fd1
	if maxFd < fd2 {
		maxFd = fd2
	}

	fds := &syscall.FdSet{}
	for {

		// 読み込み用fd_setの初期化
		// selectが毎回内容を上書きしてしまうため
		copy(fds.Bits[:], readFds.Bits[:])

		log.Println("fds is ", jsonutil.Marshal(fds))

		// fdsに設定されたソケットが読み込み可能になるまでまちます
		// 1つ目の引数はfdの最大値+1にします
		_, err = syscall.Select(maxFd+1, fds, nil, nil, nil)
		chk.SE(err)

		log.Println("after fds is ", jsonutil.Marshal(fds))

		// 全てに対して検証するから、非効率なのか...
		// 読み込み可能かどうかも、都度APで検証するみた

		// fd1に読み込み可能データがある場合
		if FD_ISSET(fds, fd1) {

			var buf [1024]byte
			i, err := syscall.Read(fd1, buf[:])
			chk.SE(err)
			fmt.Println("fd1 msg is ", string(buf[:i]))
		}

		// fd2に書き込み可能データがある場合
		if FD_ISSET(fds, fd2) {
			var buf [1024]byte
			i, err := syscall.Read(fd2, buf[:])
			chk.SE(err)
			fmt.Println("fd2 msg is ", string(buf[:i]))
		}

	}

}

func FD_SET(p *syscall.FdSet, i int) {
	p.Bits[i/64] |= 1 << (uint(i) % 64)
}

func FD_ISSET(p *syscall.FdSet, i int) bool {
	return (p.Bits[i/64] & (1 << (uint(i) % 64))) != 0
}
