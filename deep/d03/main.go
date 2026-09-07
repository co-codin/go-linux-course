package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	fmt.Println("=== Deep D3: Graceful drain vs SIGTERM ===")
	fmt.Println()

	var inFlight atomic.Int64
	mux := http.NewServeMux()
	mux.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		inFlight.Add(1)
		defer inFlight.Add(-1)
		// slow handler — stands in for an RPC
		time.Sleep(800 * time.Millisecond)
		fmt.Fprintln(w, "done")
	})

	srv := &http.Server{Addr: "127.0.0.1:18082", Handler: mux}
	go func() {
		fmt.Printf("serving on %s pid=%d\n", srv.Addr, os.Getpid())
		_ = srv.ListenAndServe()
	}()
	time.Sleep(150 * time.Millisecond)

	// start a few slow requests
	for i := 0; i < 3; i++ {
		go func() {
			resp, err := http.Get("http://127.0.0.1:18082/work")
			if err != nil {
				fmt.Println("client err:", err)
				return
			}
			resp.Body.Close()
			fmt.Println("client: got response")
		}()
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("in-flight before SIGTERM: %d\n", inFlight.Load())

	// simulate k8s SIGTERM
	p, _ := os.FindProcess(os.Getpid())
	_ = p.Signal(syscall.SIGTERM)

	ctxSig := make(chan os.Signal, 1)
	signal.Notify(ctxSig, syscall.SIGTERM, syscall.SIGINT)
	<-ctxSig

	fmt.Println("SIGTERM → Shutdown with 2s deadline (GracefulStop analogue)")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	start := time.Now()
	err := srv.Shutdown(ctx)
	fmt.Printf("Shutdown returned in %v err=%v in-flight-now=%d\n",
		time.Since(start).Round(time.Millisecond), err, inFlight.Load())
	fmt.Println()
	fmt.Println("Takeaway: stop accepting, drain in-flight, then exit — beat terminationGracePeriodSeconds.")
	time.Sleep(6 * time.Second)
}
