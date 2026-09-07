# Lesson 10 — Namespaces & cgroups

- Namespaces isolate views: pid, mnt, net, uts, ipc, user, cgroup, time. Each is an inode under `/proc/self/ns/`.
- Same inode number ⇒ same namespace. Different net ns ⇒ different loopback, routes, iptables.
- Cgroups (v2 almost everywhere now) limit/account CPU, memory, IO, pids — path in `/proc/self/cgroup`.
- A "container" is just a process tree with new namespaces + a cgroup; no hypervisor required.
- Go services usually inherit the container's ns/cgroup; you rarely call `unshare`/`clone` yourself.
- Debugging "works on host, fails in pod": compare `ns` inodes, cgroup memory.max, and net ns (`ip link`).
- Orchestrators (k8s) set these before `exec` of your binary — your PID 1 is already namespaced.
