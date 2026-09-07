# Lesson 09 — Futex / sync

- Uncontended `sync.Mutex` is an atomic CAS in user space — no syscall.
- On contention the waiter calls `FUTEX_WAIT`; unlock does `FUTEX_WAKE`. That is the Linux sync primitive under Go mutexes/condvars.
- Parked Gs do not pin an OS thread while waiting; the M runs other work (same idea as netpoller).
- `runtime.Semacquire` / `Semrelease` are the internal futex wrappers used by mutex, RWMutex, WaitGroup notes, channels.
- Channel send/recv under load also ends in sleep/wakeup on runtime semaphores — profile with `mutex` and `block` pprof.
- Spurious wakes and thundering herds are handled inside the runtime; don't reimplement futex unless you must.
- strace `-e futex` on a contended microbenchmark shows the wait/wake pairs clearly.
