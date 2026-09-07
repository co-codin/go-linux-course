package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println("=== Deep D6: sendfile / zero-copy ===")
	fmt.Println()

	dir := os.TempDir()
	srcPath := filepath.Join(dir, "d06-src.bin")
	dstPath := filepath.Join(dir, "d06-dst.bin")
	payload := make([]byte, 1<<20) // 1 MiB
	for i := range payload {
		payload[i] = byte(i)
	}
	if err := os.WriteFile(srcPath, payload, 0o644); err != nil {
		panic(err)
	}
	defer os.Remove(srcPath)
	defer os.Remove(dstPath)

	src, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	defer src.Close()
	dst, err := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	// sendfile(out, in, offset, count) — kernel copies pages, no user buffer
	var off int64
	n, err := unix.Sendfile(int(dst.Fd()), int(src.Fd()), &off, len(payload))
	if err != nil {
		panic(err)
	}
	fmt.Printf("unix.Sendfile: copied %d bytes  off_now=%d\n", n, off)

	fi, _ := dst.Stat()
	fmt.Printf("dst size=%d (expect %d)\n", fi.Size(), len(payload))

	// contrast: classic userspace copy path name
	fmt.Println()
	fmt.Println("Userspace path: read(src) → []byte → write(dst)  (2 copies into/out of user RAM)")
	fmt.Println("sendfile path:  pages stay in kernel  (what http.ServeContent may use on Linux)")
	fmt.Println()
	fmt.Println("Takeaway: prefer sendfile/splice for large static payloads; gRPC protobuf still buffers in user space.")
	time.Sleep(8 * time.Second)
}
