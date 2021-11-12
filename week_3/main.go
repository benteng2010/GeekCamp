package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//定义路由及端口
type Server struct {
	address string
	router string
	handler func(w http.ResponseWriter, r *http.Request)
}

func main() {
	route_1 := &Server{
		address: ":7001",
		router: "/f1",
		handler: f1,
	}

	route_2 := &Server{
		address: ":7002",
		router: "/f2",
		handler: f2,
	}

	servers := []*Server{
		route_1,
		route_2,
	}

	mux := http.NewServeMux()
	serverOut := make(chan struct{})
	g, ctx := errgroup.WithContext(context.Background())
	wg := sync.WaitGroup{}

	//定义信号量
	quitSign := make(chan os.Signal, 1)
	signal.Notify(quitSign, os.Interrupt, syscall.SIGTERM, syscall.SIGTSTP)
	defer signal.Stop(quitSign)

	for _, v := range servers {
		mux.HandleFunc(v.router, v.handler)
		server := &http.Server{Addr: v.address, Handler: mux}
		wg.Add(1)
		g.Go(func() error {
			wg.Done()
			return Start(server)
		})
	}
	//手动关闭http server
	wg.Add(1)
	g.Go(func() error {
		http.HandleFunc("/shutdown", func(writer http.ResponseWriter, request *http.Request) {
			serverOut <- struct{}{}
		})
		wg.Done()
		return http.ListenAndServe(":7003", nil)
	})
	wg.Wait()

	select {
		case <-ctx.Done():
			log.Println("Process finished with exit")
		case <- serverOut:
			log.Println("Process finished with shutdown")
		case <- quitSign:
			log.Println("Process finished with quit")
	}

}

func f1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello 333")
}

func f2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello 444")
}

func Start(s *http.Server) error {
	err := s.ListenAndServe()
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}

