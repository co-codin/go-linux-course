# Lesson 05 — Processes

- `os/exec` is a thin wrapper over `fork`/`clone` + `execve` + `waitpid`; Go never exposes raw `fork` safely (runtime threads).
- `Cmd.Start` returns after the child is created; `Wait` reaps it. Leaking unreaped children creates zombies.
- `ExtraFiles` maps to FDs 3,4,… in the child (0/1/2 are Stdin/Stdout/Stderr). Close unused ends in both processes.
- `Cmd.Env` / `Dir` / `SysProcAttr` (credentials, `Setsid`, `Pdeathsig`, `Cloneflags`) control the child's Linux context.
- `CombinedOutput` / `Output` buffer everything — for long-running children stream via `StdoutPipe`.
- Containers: PID namespaces make `Process.Pid` local; host tooling sees a different number.
- Prefer `exec.CommandContext` so cancel/deadline kills the process group cleanly.
