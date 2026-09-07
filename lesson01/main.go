package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

func main() {
	fmt.Println("=== Lesson 1: Your Go process on Linux ===")
	fmt.Println()

	pid := os.Getpid()
	fmt.Printf("PID (os.Getpid):           %d\n", pid)
	fmt.Printf("PPID (os.Getppid):         %d\n", os.Getppid())
	fmt.Printf("UID/GID:                   %d / %d\n", os.Getuid(), os.Getgid())
	fmt.Printf("GOMAXPROCS:                %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("NumCPU:                    %d\n", runtime.NumCPU())
	fmt.Printf("NumGoroutine (start):      %d\n", runtime.NumGoroutine())

	// Spawn goroutines — watch OS thread count in /proc, not just goroutine count
	for i := 0; i < 50; i++ {
		go func() { time.Sleep(2 * time.Second) }()
	}
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("NumGoroutine (after 50):   %d\n", runtime.NumGoroutine())

	status, _ := os.ReadFile("/proc/self/status")
	for _, line := range strings.Split(string(status), "\n") {
		if strings.HasPrefix(line, "Name:") ||
			strings.HasPrefix(line, "State:") ||
			strings.HasPrefix(line, "Threads:") ||
			strings.HasPrefix(line, "VmRSS:") ||
			strings.HasPrefix(line, "voluntary_ctxt") ||
			strings.HasPrefix(line, "nonvoluntary_ctxt") {
			fmt.Printf("/proc/self/status:         %s\n", line)
		}
	}

	cmdline, _ := os.ReadFile("/proc/self/cmdline")
	fmt.Printf("/proc/self/cmdline:        %q\n", strings.ReplaceAll(string(cmdline), "\x00", " "))

	exe, _ := os.Readlink("/proc/self/exe")
	fmt.Printf("/proc/self/exe:            %s\n", exe)

	fmt.Println()
	fmt.Println("Takeaway: 50 goroutines ≠ 50 OS threads.")
	fmt.Println("Linux sees one process; Go's scheduler multiplexes Gs onto Ms (threads).")
	fmt.Println()
	fmt.Println("(holding 3s so the terminal stays readable…)")
	time.Sleep(8 * time.Second)
}
