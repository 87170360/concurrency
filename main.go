package main

import (
	"context"
	"fmt"
	"github.com/kratos/pkg/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func serve(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello from %v", addr)
	})
	s := http.Server{Addr: addr, Handler: mux}

	go func() {
		<-ctx.Done()
		err := s.Shutdown(ctx)
		fmt.Printf("shut down serve :%v, err:%v\n", addr, err)
	}()
	return s.ListenAndServe()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	sinChan := make(chan os.Signal, 1)
	signal.Notify(sinChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		signal := <-sinChan
		fmt.Printf("receive signal :%v\n", signal)
		cancel()
	}()

	g := errgroup.WithCancel(ctx)

	var addrs = []string{
		"127.0.0.1:8001",
		"127.0.0.1:8002",
		"127.0.0.1:8003",
	}

	for _, addr := range addrs {
		addr := addr
		g.Go(func(ctx context.Context) error {
			return serve(ctx, addr)
		})
	}

	if err := g.Wait(); err != nil {
		fmt.Println(err)
	}


	fmt.Println("main quit")

}
