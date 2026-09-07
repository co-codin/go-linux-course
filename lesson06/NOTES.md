# Lesson 06 — Sockets

- `net.Listen` → `socket` + `bind` + `listen`; `Accept` → `accept4`; `Dial` → `socket` + `connect` (possibly non-blocking + epoll).
- `(*TCPConn).File()` / `(*TCPListener).File()` dup the FD and put the copy in blocking mode — use sparingly; close the File.
- `SO_REUSEADDR` is on by default for Listen in Go; `SO_REUSEPORT` needs `ListenConfig.Control` or `syscall`.
- Deadlines (`SetDeadline`) arm the netpoller's timer; they are not kernel `SO_RCVTIMEO` by default on all paths.
- gRPC and net/http share this stack: each RPC stream rides a TCP (or HTTP/2) conn FD.
- Inspect with `ss -tinp`, `tcpdump`, `/proc/net/tcp`.
- Prefer `ListenConfig` when you need socket options before bind (e.g. reuseport, mark, bind-to-device).
