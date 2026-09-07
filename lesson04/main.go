package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("=== Lesson 4: Signals (graceful shutdown) ===")
	fmt.Println()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	fmt.Printf("PID %d — pretend gRPC server serving…\n", os.Getpid())
	fmt.Println("(sending SIGTERM to self in 400ms)")

	go func() {
		time.Sleep(400 * time.Millisecond)
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(syscall.SIGTERM)
	}()

	sig := <-ch
	fmt.Printf("\ngot %v — GracefulStop / http.Server.Shutdown path\n", sig)
	fmt.Println("drain in-flight RPCs, then exit 0")
	fmt.Println()
	fmt.Println("Takeaway: Kubernetes sends SIGTERM; handle it — don't let the process die mid-RPC.")
	time.Sleep(6 * time.Second)
}
