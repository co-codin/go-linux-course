package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("=== Lesson 5: Processes (os/exec) ===")
	fmt.Println()

	out, err := exec.Command("echo", "hello-from-child").Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("exec.Command(...).Output() -> %q\n", string(out))

	cmd := exec.Command("sleep", "0.2")
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	fmt.Printf("child Start() -> Process.Pid=%d (parent pid=%d)\n", cmd.Process.Pid, os.Getpid())
	if err := cmd.Wait(); err != nil {
		panic(err)
	}
	fmt.Printf("Wait() done; ProcessState=%v\n", cmd.ProcessState)

	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	child := exec.Command("bash", "-c", "echo got-fd3-from-parent >&3")
	child.ExtraFiles = []*os.File{w}
	child.Stdout = os.Stdout
	child.Stderr = os.Stderr
	if err := child.Start(); err != nil {
		panic(err)
	}
	w.Close()
	buf := make([]byte, 64)
	n, _ := r.Read(buf)
	r.Close()
	_ = child.Wait()
	fmt.Printf("ExtraFiles pipe from child: %q\n", string(buf[:n]))

	fmt.Println()
	fmt.Println("Takeaway: os/exec is fork+exec+wait (plus FD plumbing via ExtraFiles).")
}
