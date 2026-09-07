package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	mode := "server"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "server":
		http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "pong\n")
		})
		addr := "127.0.0.1:18080"
		fmt.Println("=== Deep D1: strace a live Go server ===")
		fmt.Printf("listening on http://%s/ping  pid=%d\n", addr, os.Getpid())
		fmt.Println("in another terminal: strace -f -e accept4,read,write,epoll_wait,epoll_ctl,close -p PID")
		if err := http.ListenAndServe(addr, nil); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	case "load":
		fmt.Println("firing 40 GETs…")
		ok := 0
		for i := 0; i < 40; i++ {
			resp, err := http.Get("http://127.0.0.1:18080/ping")
			if err != nil {
				fmt.Println("err:", err)
				continue
			}
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			ok++
		}
		fmt.Printf("done ok=%d\n", ok)
	default:
		fmt.Println("usage: d01 server|load")
	}
	_ = time.Second
}
