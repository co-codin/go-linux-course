package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("=== Lesson 9: Futex / sync ===")
	fmt.Println()

	var mu sync.Mutex
	mu.Lock()
	fmt.Printf("main holds mutex; NumGoroutine=%d\n", runtime.NumGoroutine())

	const n = 5
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()
			fmt.Printf("  g%d: Lock() … blocking\n", id)
			mu.Lock()
			fmt.Printf("  g%d: acquired\n", id)
			mu.Unlock()
		}(i)
	}

	time.Sleep(200 * time.Millisecond)
	fmt.Printf("after spawn (still holding): NumGoroutine=%d (Gs parked on mutex)\n", runtime.NumGoroutine())
	fmt.Println("main Unlock — waiters proceed via futex wake")
	mu.Unlock()
	wg.Wait()

	fmt.Printf("done; NumGoroutine=%d\n", runtime.NumGoroutine())
	fmt.Println()
	fmt.Println("Takeaway: contended mutex uses futex syscall — uncontended path stays in user space.")
}
