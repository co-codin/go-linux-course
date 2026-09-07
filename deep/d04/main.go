package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

func rssKB() string {
	b, _ := os.ReadFile("/proc/self/status")
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "VmRSS:") {
			return strings.TrimSpace(line)
		}
	}
	return "VmRSS: ?"
}

func memstats(label string) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Printf("%s HeapAlloc=%d MiB Sys=%d MiB | %s\n",
		label, ms.HeapAlloc>>20, ms.Sys>>20, rssKB())
}

func main() {
	fmt.Println("=== Deep D4: GOMEMLIMIT vs VmRSS ===")
	fmt.Println()

	limit := int64(64 << 20) // 64 MiB soft limit for GC
	debug.SetMemoryLimit(limit)
	fmt.Printf("debug.SetMemoryLimit(%d MiB)  (env GOMEMLIMIT works the same)\n", limit>>20)
	memstats("start")

	// allocate past the soft limit in chunks, touching pages
	var hold [][]byte
	for i := 0; i < 8; i++ {
		chunk := make([]byte, 16<<20)
		for j := 0; j < len(chunk); j += 4096 {
			chunk[j] = 1
		}
		hold = append(hold, chunk)
		runtime.GC()
		memstats(fmt.Sprintf("after chunk %d (holding ~%d MiB)", i+1, (i+1)*16))
	}

	fmt.Println()
	fmt.Println("Go will GC harder near GOMEMLIMIT; cgroup OOM still keys off RSS/working set.")
	fmt.Printf("cgroup path: ")
	cg, _ := os.ReadFile("/proc/self/cgroup")
	fmt.Print(string(cg))
	fmt.Println()
	fmt.Println("Takeaway: set GOMEMLIMIT below cgroup memory.max — leave headroom for stacks/caches.")
	runtime.KeepAlive(hold)
	time.Sleep(8 * time.Second)
}
