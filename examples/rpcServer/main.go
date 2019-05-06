package main

import (
	"github.com/bary321/NetWorkRate"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	common := new(NetWorkRate.Common)
	if err := rpc.Register(common); err != nil {
		log.Fatal(err)
	}
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}

	go func() {
		log.Println("start sar server")
	}()

	if err := http.Serve(l, nil); err != nil {
		log.Fatal(err)
	}
}
