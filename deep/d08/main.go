package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Deep D8: pprof mutex/block + futex ===")
	fmt.Println()

	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)

	go func() {
		_ = http.ListenAndServe("127.0.0.1:18088", nil) // pprof endpoints
	}()
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("pprof on http://127.0.0.1:18088/debug/pprof/  pid=%d\n", os.Getpid())
	fmt.Println("  mutex: /debug/pprof/mutex   block: /debug/pprof/block")

	var mu sync.Mutex
	var wg sync.WaitGroup
	// create contention
	start := time.Now()
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 200; j++ {
				mu.Lock()
				time.Sleep(50 * time.Microsecond)
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	fmt.Printf("contended work done in %v\n", time.Since(start).Round(time.Millisecond))
	fmt.Println()
	fmt.Println("Under the hood: FUTEX_WAIT / FUTEX_WAKE (strace -e futex -f -p PID during the loop).")
	fmt.Println("pprof mutex/block profiles attribute that wait time back to Lock sites.")
	fmt.Println()
	fmt.Println("Takeaway: hot mutex → futex syscalls → see it in strace AND pprof — fix the lock, not the kernel.")
	time.Sleep(8 * time.Second)
}
