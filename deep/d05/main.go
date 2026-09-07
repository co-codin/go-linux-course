package main

import (
	"fmt"
	"net"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println("=== Deep D5: TCP_INFO / socket options ===")
	fmt.Println()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	go func() {
		c, err := ln.Accept()
		if err != nil {
			return
		}
		defer c.Close()
		buf := make([]byte, 16)
		_, _ = c.Read(buf)
		_, _ = c.Write([]byte("ok"))
		time.Sleep(300 * time.Millisecond)
	}()

	c, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		panic(err)
	}
	defer c.Close()
	tc := c.(*net.TCPConn)
	raw, err := tc.SyscallConn()
	if err != nil {
		panic(err)
	}

	var rcv, snd int
	var info unix.TCPInfo
	err = raw.Control(func(fd uintptr) {
		rcv, _ = unix.GetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_RCVBUF)
		snd, _ = unix.GetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_SNDBUF)
		// TCP_INFO
		size := unsafe.Sizeof(info)
		_, _, errno := syscall.Syscall6(
			syscall.SYS_GETSOCKOPT,
			fd,
			uintptr(unix.IPPROTO_TCP),
			uintptr(unix.TCP_INFO),
			uintptr(unsafe.Pointer(&info)),
			uintptr(unsafe.Pointer(&size)),
			0,
		)
		if errno != 0 {
			fmt.Println("TCP_INFO errno:", errno)
		}
	})
	if err != nil {
		panic(err)
	}

	_, _ = c.Write([]byte("hi"))
	time.Sleep(50 * time.Millisecond)

	fmt.Printf("conn %s ↔ %s\n", c.LocalAddr(), c.RemoteAddr())
	fmt.Printf("SO_RCVBUF=%d  SO_SNDBUF=%d\n", rcv, snd)
	fmt.Printf("TCP_INFO: state=%d rtt=%dus rttvar=%dus snd_cwnd=%d total_retrans=%d\n",
		info.State, info.Rtt, info.Rttvar, info.Snd_cwnd, info.Total_retrans)
	fmt.Println()
	fmt.Println("Takeaway: SyscallConn().Control is how you read TCP_INFO / set keepalive from Go.")
	fmt.Println("ss -ti shows the same kernel struct for your gRPC sockets.")
	time.Sleep(8 * time.Second)
}
