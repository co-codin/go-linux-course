package main

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println("=== Lesson 2: Syscalls (x/sys/unix) ===")
	fmt.Println()

	path := "/etc/hosts"
	fd, err := unix.Open(path, unix.O_RDONLY, 0)
	if err != nil {
		fmt.Println("open:", err)
		os.Exit(1)
	}
	defer unix.Close(fd)

	buf := make([]byte, 64)
	n, err := unix.Read(fd, buf)
	if err != nil {
		fmt.Println("read:", err)
		os.Exit(1)
	}

	fmt.Printf("unix.Open(%q) -> fd=%d\n", path, fd)
	fmt.Printf("unix.Read(fd,%d) -> n=%d data=%q\n", len(buf), n, string(buf[:n]))
	fmt.Println()
	fmt.Println("These are direct Linux syscalls (openat/read/close).")
	fmt.Println("os.Open wraps the same thing with buffering + *os.File.")
	fmt.Println()
	fmt.Println("Next in the terminal: strace -e openat,read,close ./lesson02")
}
