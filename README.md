# Go + Linux Course

Senior Go/gRPC track: how Go sits on the Linux ABI — from process model and syscalls through epoll, cgroups, and netlink.

## What's inside

- **Lessons 1–12** — core course (`lesson01` … `lesson12`)
- **Deep labs D1–D8** — hands-on labs after the core lessons (`deep/d01` … `deep/d08`)
- **gRPC lab** — ListenConfig + GracefulStop checklist (`deep/grpc-lab`)

## How to run

```bash
cd lesson01 && go run .
```

Same pattern for any lesson or deep lab:

```bash
cd lessonNN && go run .
cd deep/d0N && go run .
cd deep/grpc-lab && go run .
```

## Structure

| Path | Contents |
|------|----------|
| `lesson01/` … `lesson12/` | Core lessons (`main.go`, `go.mod`, optional `NOTES.md`, `run-for-shot.sh`) |
| `deep/d01/` … `deep/d08/` | Deep labs |
| `deep/grpc-lab/` | gRPC ListenConfig + GracefulStop lab |
| `CURRICULUM.md` | Lesson outline |
| `DEEP_CURRICULUM.md` | Deep-lab outline |

See [CURRICULUM.md](CURRICULUM.md) and [DEEP_CURRICULUM.md](DEEP_CURRICULUM.md) for the full outline.
