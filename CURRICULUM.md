# Go + Linux — course for senior Go (web/gRPC)

Assumes: strong Go, HTTP/gRPC services. Goal: how Go sits on the Linux ABI.

1. Process model — PID, threads, /proc, G vs M vs P
2. Syscalls — syscall vs x/sys/unix, strace a Go binary
3. File descriptors — tables, dup, inheritance, os.File
4. Signals — os/signal, SIGTERM graceful shutdown (your gRPC servers)
5. Processes — os/exec, fork/exec under the hood, wait
6. Sockets — net FD, SO_REUSEPORT, TCP under Dial/Listen
7. Netpoller — epoll, how net/http & gRPC wait without blocking Ms
8. Memory — heap, mmap, /proc/self/maps, GC vs RSS
9. Sync with the kernel — futex, runtime semaphores
10. Namespaces & cgroups — what containers actually are
11. Capabilities & seccomp — least privilege for services
12. Netlink & /sys — when stdlib is not enough
