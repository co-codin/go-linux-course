package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	fmt.Println("=== Lesson 6: Sockets ===")
	fmt.Println()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	fmt.Printf("listening on %s\n", ln.Addr())

	if tl, ok := ln.(*net.TCPListener); ok {
		f, err := tl.File()
		if err == nil {
			fmt.Printf("listen socket FD (via File()): %d\n", f.Fd())
			f.Close()
		}
	}

	done := make(chan struct{})
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		line, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Printf("server read: %q\n", line)
		fmt.Fprintf(conn, "pong\n")
		close(done)
	}()

	conn, err := net.Dial("tcp", ln.Addr().String())
	if err != nil {
		panic(err)
	}
	if tc, ok := conn.(*net.TCPConn); ok {
		f, err := tc.File()
		if err == nil {
			fmt.Printf("client conn FD (via File()): %d\n", f.Fd())
			f.Close() // File() dups; original conn still usable
		}
	}
	fmt.Fprintf(conn, "ping\n")
	reply, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf("client read: %q\n", reply)
	conn.Close()
	<-done

	fmt.Println()
	fmt.Println("Takeaway: net.Conn is a socket FD — Dial/Listen wrap socket(2)/bind/listen/accept.")
}
