package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Minimal RTM_GETROUTE dump via netlink
func dumpRoutes() error {
	fd, err := unix.Socket(unix.AF_NETLINK, unix.SOCK_RAW, unix.NETLINK_ROUTE)
	if err != nil {
		return err
	}
	defer unix.Close(fd)

	sa := &unix.SockaddrNetlink{Family: unix.AF_NETLINK}
	if err := unix.Bind(fd, sa); err != nil {
		return err
	}

	// nlmsg: NLMSG_HDRLEN + rtmsg
	type nlMsgHdr struct {
		Len   uint32
		Type  uint16
		Flags uint16
		Seq   uint32
		Pid   uint32
	}
	const (
		NLMSG_HDRLEN = 16
		RTM_GETROUTE = 26
		NLM_F_REQUEST = 0x01
		NLM_F_DUMP    = 0x300
	)
	msgLen := uint32(NLMSG_HDRLEN + unix.SizeofRtMsg)
	buf := make([]byte, msgLen)
	hdr := nlMsgHdr{Len: msgLen, Type: RTM_GETROUTE, Flags: NLM_F_REQUEST | NLM_F_DUMP, Seq: 1, Pid: 0}
	*(*nlMsgHdr)(unsafe.Pointer(&buf[0])) = hdr
	// rtmsg zeroed = all tables/families dump

	if err := unix.Sendto(fd, buf, 0, sa); err != nil {
		return err
	}

	rb := make([]byte, 65536)
	n, _, err := unix.Recvfrom(fd, rb, 0)
	if err != nil {
		return err
	}
	fmt.Printf("netlink RTM_GETROUTE dump: %d bytes first recv\n", n)

	// walk messages briefly
	off := 0
	routes := 0
	for off+16 <= n {
		l := binary.LittleEndian.Uint32(rb[off:])
		t := binary.LittleEndian.Uint16(rb[off+4:])
		if l < 16 || int(off)+int(l) > n {
			break
		}
		if t == unix.NLMSG_DONE {
			break
		}
		if t == unix.NLMSG_ERROR {
			fmt.Println("NLMSG_ERROR in dump")
			break
		}
		if t == RTM_NEWROUTE || t == 24 { // RTM_NEWROUTE=24
			routes++
		}
		off += int(l)
		// align
		off = (off + 3) & ^3
	}
	fmt.Printf("parsed RTM_NEWROUTE msgs in first buffer: %d (may need multi-recv for full table)\n", routes)
	return nil
}

const RTM_NEWROUTE = 24

func main() {
	fmt.Println("=== Deep D7: netlink routes ===")
	fmt.Println()

	ifaces, _ := net.Interfaces()
	fmt.Println("stdlib net.Interfaces (no netlink needed):")
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		fmt.Printf("  %s: %v\n", iface.Name, addrs)
	}

	fmt.Println()
	fmt.Println("/proc/net/route (legacy view kernel still exposes):")
	b, _ := os.ReadFile("/proc/net/route")
	lines := strings.Split(string(b), "\n")
	for i, line := range lines {
		if i >= 5 || line == "" {
			break
		}
		fmt.Println(" ", line)
	}

	fmt.Println()
	if err := dumpRoutes(); err != nil {
		fmt.Println("netlink dump:", err)
	} else {
		fmt.Println("AF_NETLINK socket worked — same channel as `ip route`")
	}

	fmt.Println()
	fmt.Println("Takeaway: read with stdlib;/proc; change routing/addrs via netlink (+ CAP_NET_ADMIN).")
	_ = syscall.AF_NETLINK
	time.Sleep(8 * time.Second)
}
