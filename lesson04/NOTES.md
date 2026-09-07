# Lesson 04 — Signals

- `signal.Notify` routes POSIX signals into a channel; without it, SIGTERM kills the process immediately.
- Always `signal.Stop` and drain when done — otherwise the runtime keeps a notifier goroutine alive.
- Default SIGPIPE/SIGURG handling differs in Go; don't casually `Notify` everything.
- gRPC: on SIGTERM call `GracefulStop` (or `Stop` with a deadline) so in-flight RPCs finish.
- `unix.Kill(pid, sig)` / `syscall.Kill` are the same ABI as `kill(2)`; useful for self-tests and process managers.
- PID 1 in a container receives unreaped children's signals differently — always install a handler in the entrypoint binary.
- Prefer catching SIGTERM+SIGINT; leave SIGKILL alone (it is not deliverable to userspace).
