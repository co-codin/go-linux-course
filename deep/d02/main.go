package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

func listenReuse(addr string) (net.Listener, error) {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var opErr error
			err := c.Control(func(fd uintptr) {
				opErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
				if opErr != nil {
					return
				}
				// SO_REUSEPORT = 15 on linux
				opErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, 0xf, 1)
			})
			if err != nil {
				return err
			}
			return opErr
		},
	}
	return lc.Listen(context.Background(), "tcp", addr)
}

func main() {
	fmt.Println("=== Deep D2: ListenConfig + SO_REUSEPORT ===")
	fmt.Println()

	addr := "127.0.0.1:18081"
	ln1, err := listenReuse(addr)
	if err != nil {
		panic(err)
	}
	defer ln1.Close()
	ln2, err := listenReuse(addr)
	if err != nil {
		fmt.Println("second Listen FAILED (reuseport unsupported?):", err)
		os.Exit(1)
	}
	defer ln2.Close()

	fmt.Printf("two listeners on SAME addr %s\n", addr)
	fmt.Println("kernel will distribute accepts across both (multi-process fan-in pattern)")

	counts := [2]int{}
	serve := func(id int, ln net.Listener) {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			counts[id]++
			c.Close()
		}
	}
	go serve(0, ln1)
	go serve(1, ln2)

	// dial many times
	for i := 0; i < 100; i++ {
		c, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println("dial:", err)
			continue
		}
		c.Close()
	}
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("accepts: listener0=%d listener1=%d (sum=%d)\n", counts[0], counts[1], counts[0]+counts[1])
	fmt.Println()
	fmt.Println("Takeaway: SO_REUSEPORT lets N processes bind one port — classic Go/nginx scale-out.")
	fmt.Println("Wire it via net.ListenConfig.Control before bind.")
}
