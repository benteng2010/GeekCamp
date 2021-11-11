package main

import (
	"fmt"
	"net/http"
	"golang.org/x/sync/errgroup"
)

//定义路由及端口
type Server struct {
	port string
	router string
	handler func(w http.ResponseWriter, r *http.Request)
}

func main() {
	route_1 := &Server{
		port: ":7001",
		router: "/f1",
		handler: f1,
	}

	route_2 := &Server{
		port: ":7002",
		router: "/f2",
		handler: f2,
	}

	servers := []*Server{
		route_1,
		route_2,
	}


}

func f1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello f1")
}

func f2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello f2")
}


https://benteng2010:123456@github.com