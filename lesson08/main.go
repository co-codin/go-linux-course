package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"unsafe"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println("=== Lesson 8: Memory ===")
	fmt.Println()

	printMem := func(label string) {
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		fmt.Printf("%s MemStats: Alloc=%d MiB HeapAlloc=%d MiB Sys=%d MiB\n",
			label, ms.Alloc/1024/1024, ms.HeapAlloc/1024/1024, ms.Sys/1024/1024)
		status, _ := os.ReadFile("/proc/self/status")
		for _, line := range strings.Split(string(status), "\n") {
			if strings.HasPrefix(line, "VmRSS:") || strings.HasPrefix(line, "VmSize:") {
				fmt.Printf("%s /proc/self/status: %s\n", label, line)
			}
		}
	}

	printMem("before")

	big := make([]byte, 32<<20) // 32 MiB
	for i := 0; i < len(big); i += 4096 {
		big[i] = 1 // touch pages so RSS grows
	}
	runtime.KeepAlive(big)
	printMem("after 32MiB []byte")

	page, err := unix.Mmap(-1, 0, unix.Getpagesize(),
		unix.PROT_READ|unix.PROT_WRITE,
		unix.MAP_PRIVATE|unix.MAP_ANON)
	if err != nil {
		panic(err)
	}
	page[0] = 0xAB
	fmt.Printf("unix.Mmap anon page: len=%d firstByte=0x%02x ptr=%p\n", len(page), page[0], unsafe.Pointer(&page[0]))
	_ = unix.Munmap(page)

	fmt.Println()
	fmt.Println("Takeaway: Go heap ≠ RSS exactly; mmap regions show up in /proc/self/maps separately.")
}
