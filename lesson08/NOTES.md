# Lesson 08 — Memory

- `runtime.MemStats.HeapAlloc` is live Go heap objects; `Sys` is memory obtained from the OS; RSS includes heap + stacks + mmap + allocator free lists.
- Touching pages (`big[i]=1`) forces physical allocation; untouched `make([]byte, N)` may only grow virtual size (VmSize).
- `unix.Mmap` with `MAP_ANON` is how the runtime and many CGO libs get large arenas; visible in `/proc/self/maps`.
- GC frees heap logically; the scavenger returns pages to the OS (`MADV_DONTNEED` / `MADV_FREE`) asynchronously — RSS lags.
- Ballast / `GOMEMLIMIT` / `debug.SetMemoryLimit` matter more than manual free for services.
- Container OOM kills use cgroup memory, not Go's view — export metrics of both HeapAlloc and cgroup usage.
- For leaks: `pprof` heap profiles + compare `/proc/self/maps` growth.
