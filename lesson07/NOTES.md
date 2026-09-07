# Lesson 07 — Netpoller / epoll

- Linux readiness is `epoll_create1` + `epoll_ctl` + `epoll_wait`; Go's runtime embeds this as the netpoller.
- A blocked `Read`/`Accept` parks the G; the M is free to run other Gs until epoll reports the FD ready.
- `eventfd`/`timerfd` are the same class of pollable FDs as sockets — useful for demos and custom wakeups.
- Level vs edge trigger (`EPOLLET`): Go uses edge-triggered semantics carefully; user code usually stays level via stdlib.
- Thousands of idle conns cost FDs + tiny netpoller state, not one OS thread each — that is why Go servers scale.
- `runtime.Netpoll` / `internal/poll` are the source of truth if you need to dig.
- strace: watch `epoll_ctl` on Dial/Listen and `epoll_wait` under load.
