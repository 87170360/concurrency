package main

import (
	"context"
	"fmt"
	"github.com/kratos/pkg/sync/errgroup"
	"net/http"
)

func serve(ctx context.Context, addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Hello from %v", addr)
	})
	s := http.Server{Addr: addr, Handler: mux}

	go func() {
		<-ctx.Done()

	}()
	return s.ListenAndServe()
}

func main() {
	g := new(errgroup.Group)
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
}
