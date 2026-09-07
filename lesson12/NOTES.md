# Lesson 12 — Netlink & /sys

- `/sys/class/net` exposes iface operstate, MAC, MTU, queues — stable for read-only introspection.
- `/proc/net/dev` and `/proc/net/tcp` are the classic counters; `ss`/`ip` use netlink instead.
- Netlink (`AF_NETLINK`) is the bidirectional control plane: routes, addresses, neigh, qdisc, conntrack.
- Go stdlib (`net.Interfaces`, `net.InterfaceAddrs`) wraps enough for listing; advanced ops need netlink libs (`github.com/vishvananda/netlink`) or raw sockets.
- Changing routes/addresses from Go ⇒ netlink messages, not writes to `/sys` (many attrs are read-only).
- In containers you only see the netns's interfaces — host eth0 may be a veth peer.
- When stdlib is not enough: netlink for routing/policy; `/sys` for device/driver knobs; keep privileges minimal (`CAP_NET_ADMIN`).
