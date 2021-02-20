package selectt

import (
	"fmt"
	"iolearn/pkg/common/chk"
	"iolearn/pkg/common/jsonutil"
	"log"
	"syscall"
)

// 同期ノンブロッキング
// プロセスがカーネルにシステムコールする
// プロセスはカーネルの返事を「待たない」
// プロセスは任意のタイミングでIOの状態をカーネルに問い合わせる(polling)

// netcat -u localhost 1235
// netcat -u localhost 1236

// Server .
func Server() {

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
