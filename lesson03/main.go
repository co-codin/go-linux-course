package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	fmt.Println("=== Lesson 3: File descriptors ===")
	fmt.Println()

	f, err := os.Open("/etc/hosts")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	// *os.File and TCP listener are both FDs in the same table
	fmt.Printf("os.Open(/etc/hosts)  fd=%d  (*os.File)\n", f.Fd())
	if lf, err := ln.(*net.TCPListener).File(); err == nil {
		fmt.Printf("net.Listen(tcp)      fd=%d  (socket)  addr=%s\n", lf.Fd(), ln.Addr())
		lf.Close() // dup; safe to close the copy
	}

	fmt.Println()
	fmt.Println("/proc/self/fd →")
	entries, _ := os.ReadDir("/proc/self/fd")
	for _, e := range entries {
		link, err := os.Readlink(filepath.Join("/proc/self/fd", e.Name()))
		if err != nil {
			continue
		}
		// skip noisy memfd/internal unless useful
		if strings.Contains(link, "anon_inode:[io_uring]") {
			continue
		}
		fmt.Printf("  fd %s -> %s\n", e.Name(), link)
	}

	fmt.Println()
	fmt.Println("Takeaway: files, sockets, and listeners share one FD table.")
	fmt.Println("ulimit -n caps your gRPC fan-in before Go ever complains.")
	time.Sleep(8 * time.Second)
}
