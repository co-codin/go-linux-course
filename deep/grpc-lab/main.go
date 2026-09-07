package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

const addr = "127.0.0.1:18090"

type greeter struct {
	pb.UnimplementedGreeterServer
}

func (g *greeter) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	// slow RPC — stands in for real work during drain
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(600 * time.Millisecond):
	}
	return &pb.HelloReply{Message: "hello " + in.GetName()}, nil
}

func listenReusePort(address string) (net.Listener, error) {
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var opErr error
			if err := c.Control(func(fd uintptr) {
				opErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
				if opErr != nil {
					return
				}
				opErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, 0xf /* SO_REUSEPORT */, 1)
			}); err != nil {
				return err
			}
			return opErr
		},
	}
	return lc.Listen(context.Background(), "tcp", address)
}

func runServer() {
	fmt.Println("=== gRPC lab: ListenConfig + GracefulStop checklist ===")
	fmt.Println()

	ln, err := listenReusePort(addr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("1. ListenConfig+SO_REUSEPORT bound %s  fd ok  pid=%d\n", addr, os.Getpid())

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &greeter{})

	go func() {
		fmt.Println("2. Serve(ln) — custom listener, not net.Listen inside grpc")
		if err := s.Serve(ln); err != nil {
			fmt.Println("Serve ended:", err)
		}
	}()
	time.Sleep(150 * time.Millisecond)

	// fire in-flight RPCs
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	for i := 0; i < 3; i++ {
		go func(i int) {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: fmt.Sprintf("rpc-%d", i)})
			if err != nil {
				fmt.Printf("   client rpc-%d err: %v\n", i, err)
				return
			}
			fmt.Printf("   client got: %s\n", resp.GetMessage())
		}(i)
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Println("3. 3 in-flight SayHello (600ms each)")

	// checklist on SIGTERM
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		time.Sleep(200 * time.Millisecond)
		fmt.Println("4. SIGTERM (k8s) → begin checklist")
		p, _ := os.FindProcess(os.Getpid())
		_ = p.Signal(syscall.SIGTERM)
	}()

	<-sigCh
	fmt.Println("5. stop accepting new RPCs: GracefulStop()")
	done := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		fmt.Println("6. drained in-flight — GracefulStop returned")
	case <-time.After(2 * time.Second):
		fmt.Println("6. deadline hit — Stop() hard kill remaining")
		s.Stop()
		<-done
	}

	fmt.Println()
	fmt.Println("Checklist:")
	fmt.Println("  [x] bind via ListenConfig (SO_REUSEPORT)")
	fmt.Println("  [x] grpc.Serve(customListener)")
	fmt.Println("  [x] signal.Notify(SIGTERM)")
	fmt.Println("  [x] GracefulStop + hard-Stop deadline")
	fmt.Println("  [ ] preStop sleep / LB deny (cluster-side — not in-process)")
	fmt.Println("  [ ] terminationGracePeriodSeconds > drain budget")
	fmt.Println()
	fmt.Println("Takeaway: wire the listener yourself; GracefulStop is Shutdown for gRPC.")
	time.Sleep(6 * time.Second)
}

func main() {
	runServer()
}
