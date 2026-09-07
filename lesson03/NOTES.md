# Lesson 03 — File descriptors

- Every open resource in a Linux process is a number in the FD table (`/proc/self/fd`).
- `os.File.Fd()` and `(*net.TCPListener).File()` expose the underlying integer; the latter dups it — close the File or you leak.
- Inheritance: FDs marked close-on-exec (Go's default for most opens) are not passed across `exec`; unset CLOEXEC when you intentionally hand FDs to children (`ExtraFiles`).
- `dup`/`dup2`/`fcntl(F_DUPFD)` create aliases to the same open-file description (shared offset, flags).
- Soft/hard `RLIMIT_NOFILE` caps how many FDs you can hold; hit it and `accept`/`Dial` fail with `EMFILE`/`ENFILE`.
- gRPC/HTTP servers are just listen FDs + accepted conn FDs multiplexed by the netpoller.
- Debugging: `ls -l /proc/<pid>/fd`, `lsof -p <pid>`, `ss -ltnp`.
