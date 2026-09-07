package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println("=== Lesson 10: Namespaces & cgroups ===")
	fmt.Println()

	var uts unix.Utsname
	if err := unix.Uname(&uts); err == nil {
		fmt.Printf("uname: sys=%s release=%s machine=%s nodename=%s\n",
			cstring(uts.Sysname[:]), cstring(uts.Release[:]),
			cstring(uts.Machine[:]), cstring(uts.Nodename[:]))
	}

	fmt.Println()
	fmt.Println("/proc/self/ns/* (namespace inodes):")
	entries, err := os.ReadDir("/proc/self/ns")
	if err != nil {
		panic(err)
	}
	for _, e := range entries {
		target, err := os.Readlink(filepath.Join("/proc/self/ns", e.Name()))
		if err != nil {
			fmt.Printf("  %-8s -> <err %v>\n", e.Name(), err)
			continue
		}
		fmt.Printf("  %-8s -> %s\n", e.Name(), target)
	}

	fmt.Println()
	fmt.Println("/proc/self/cgroup:")
	cg, err := os.ReadFile("/proc/self/cgroup")
	if err != nil {
		panic(err)
	}
	for i, line := range strings.Split(strings.TrimSpace(string(cg)), "\n") {
		if i >= 8 {
			fmt.Println("  …")
			break
		}
		fmt.Printf("  %s\n", line)
	}

	fmt.Println()
	fmt.Println("Takeaway: container = namespaces + cgroups; your process already has ns inodes.")
}

func cstring(b []byte) string {
	i := 0
	for i < len(b) && b[i] != 0 {
		i++
	}
	return string(b[:i])
}
