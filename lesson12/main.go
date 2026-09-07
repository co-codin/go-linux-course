package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("=== Lesson 12: Netlink & /sys ===")
	fmt.Println()

	fmt.Println("net.Interfaces():")
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range ifaces {
		fmt.Printf("  %-8s mtu=%d flags=%v\n", iface.Name, iface.MTU, iface.Flags)
	}

	fmt.Println()
	fmt.Println("/sys/class/net/*/operstate & address_assign_type:")
	sysNet := "/sys/class/net"
	entries, err := os.ReadDir(sysNet)
	if err != nil {
		panic(err)
	}
	shown := 0
	for _, e := range entries {
		name := e.Name()
		if name != "lo" && shown >= 2 {
			continue
		}
		oper, _ := os.ReadFile(filepath.Join(sysNet, name, "operstate"))
		aat, _ := os.ReadFile(filepath.Join(sysNet, name, "address_assign_type"))
		addr, _ := os.ReadFile(filepath.Join(sysNet, name, "address"))
		fmt.Printf("  %s: operstate=%s address_assign_type=%s address=%s\n",
			name,
			strings.TrimSpace(string(oper)),
			strings.TrimSpace(string(aat)),
			strings.TrimSpace(string(addr)))
		if name != "lo" {
			shown++
		}
	}

	fmt.Println()
	fmt.Println("/proc/net/dev (first lines):")
	dev, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(dev), "\n")
	for i, line := range lines {
		if i >= 6 || line == "" {
			break
		}
		fmt.Printf("  %s\n", line)
	}

	fmt.Println()
	fmt.Println("Takeaway: kernel control plane is /sys + netlink; stdlib covers common cases.")
}
