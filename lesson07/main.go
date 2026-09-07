package main

import (
	"fmt"
	"os"
	"unsafe"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println("=== Lesson 7: Netpoller / epoll ===")
	fmt.Println()

	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		panic(err)
	}
	defer unix.Close(epfd)
	fmt.Printf("EpollCreate1 -> epfd=%d\n", epfd)

	// eventfd is a simple fd the kernel can signal; Go's netpoller uses epoll on sockets the same way.
	efd, err := unix.Eventfd(0, unix.EFD_NONBLOCK|unix.EFD_CLOEXEC)
	if err != nil {
		panic(err)
	}
	defer unix.Close(efd)
	fmt.Printf("Eventfd -> efd=%d\n", efd)

	ev := unix.EpollEvent{Events: unix.EPOLLIN, Fd: int32(efd)}
	if err := unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, efd, &ev); err != nil {
		panic(err)
	}
	fmt.Println("EpollCtl ADD efd with EPOLLIN")

	// First wait: nothing ready → timeout
	events := make([]unix.EpollEvent, 4)
	n, err := unix.EpollWait(epfd, events, 100) // 100ms
	if err != nil {
		panic(err)
	}
	fmt.Printf("EpollWait (before signal) -> n=%d (timeout expected)\n", n)

	// Signal eventfd (write uint64 counter)
	var one uint64 = 1
	buf := (*(*[8]byte)(unsafe.Pointer(&one)))[:]
	if _, err := unix.Write(efd, buf); err != nil {
		panic(err)
	}
	fmt.Println("wrote 1 to eventfd")

	n, err = unix.EpollWait(epfd, events, 1000)
	if err != nil {
		panic(err)
	}
	fmt.Printf("EpollWait (after signal) -> n=%d events[0].Fd=%d Events=0x%x\n", n, events[0].Fd, events[0].Events)

	fmt.Println()
	fmt.Println("Go runtime netpoller: same epoll_wait loop wakes Ms when sockets are ready.")
	fmt.Printf("(pid=%d — strace -e epoll_wait,epoll_ctl ./lesson07)\n", os.Getpid())
	fmt.Println()
	fmt.Println("Takeaway: netpoller is epoll under the hood — goroutines park without blocking OS threads.")
}
