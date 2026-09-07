package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

// A few well-known capability bits (linux/capability.h)
var capNames = map[uint]string{
	0:  "CAP_CHOWN",
	1:  "CAP_DAC_OVERRIDE",
	5:  "CAP_KILL",
	6:  "CAP_SETGID",
	7:  "CAP_SETUID",
	12: "CAP_NET_ADMIN",
	13: "CAP_NET_RAW",
	21: "CAP_SYS_ADMIN",
}

func main() {
	fmt.Println("=== Lesson 11: Capabilities ===")
	fmt.Println()

	status, err := os.ReadFile("/proc/self/status")
	if err != nil {
		panic(err)
	}
	var capEff, capPrm, capBnd string
	for _, line := range strings.Split(string(status), "\n") {
		switch {
		case strings.HasPrefix(line, "CapEff:"):
			capEff = strings.TrimSpace(strings.TrimPrefix(line, "CapEff:"))
			fmt.Println(line)
		case strings.HasPrefix(line, "CapPrm:"):
			capPrm = strings.TrimSpace(strings.TrimPrefix(line, "CapPrm:"))
			fmt.Println(line)
		case strings.HasPrefix(line, "CapBnd:"):
			capBnd = strings.TrimSpace(strings.TrimPrefix(line, "CapBnd:"))
			fmt.Println(line)
		case strings.HasPrefix(line, "Uid:") || strings.HasPrefix(line, "Gid:"):
			fmt.Println(line)
		}
	}

	fmt.Println()
	fmt.Println("decoded CapEff bits (known names):")
	decodeCaps(capEff)
	_ = capPrm
	_ = capBnd

	err = unix.Access("/etc/hosts", unix.R_OK)
	fmt.Printf("unix.Access(/etc/hosts, R_OK) -> %v\n", err)

	fmt.Printf("euid=%d egid=%d (setuid binaries elevate; modern services drop caps instead)\n",
		os.Geteuid(), os.Getegid())

	fmt.Println()
	fmt.Println("Takeaway: modern services drop caps instead of running as root.")
}

func decodeCaps(hexStr string) {
	v, err := strconv.ParseUint(hexStr, 16, 64)
	if err != nil {
		fmt.Printf("  parse error: %v\n", err)
		return
	}
	if v == 0 {
		fmt.Println("  (none)")
		return
	}
	for bit, name := range capNames {
		if v&(1<<bit) != 0 {
			fmt.Printf("  bit %d: %s\n", bit, name)
		}
	}
	fmt.Printf("  CapEff=0x%x\n", v)
}
