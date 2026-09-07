# Lesson 11 — Capabilities

- Root is split into ~40 capabilities (`CapEff`/`CapPrm`/`CapBnd`/`CapInh` in `/proc/self/status`).
- Bounding set (`CapBnd`) caps what you can ever raise; drop it early in the entrypoint.
- Non-root + `CAP_NET_BIND_SERVICE` lets you bind :443 without full root — prefer this over setuid.
- Kubernetes: `securityContext.capabilities.drop: ["ALL"]` then add back only what you need.
- Go binaries are often static and don't need ambient caps; if you do, set them on the file via `setcap` or the runtime.
- `unix.Access` / open failures with EACCES are the day-to-day signal you dropped too hard — or correctly.
- Pair with seccomp (`Restricted` profiles) so even gained caps can't call dangerous syscalls.
