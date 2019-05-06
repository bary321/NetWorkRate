package main

import (
	"encoding/json"
	"fmt"
	"github.com/bary321/NetWorkRate"
	"log"
	"net/rpc"
	"sync"
)

var (
	Servers = []string{"192.168.2.90", "192.168.2.91"}
)

func GetRates(client *rpc.Client, wg *sync.WaitGroup, rate *NetWorkRate.IORates) {
	defer wg.Done()
	args := &NetWorkRate.Args{1}

	r := new(NetWorkRate.IORates)
	divCall := client.Go("Common.GetRate", args, &r, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	// check errors, print, etc.
	if replyCall.Error != nil {
		log.Fatal(replyCall.Error)
	}
	d, _ := json.Marshal(r)
	// fmt.Println(r)
	if err := json.Unmarshal(d, &rate); err != nil {
		log.Println(err)
	}
}

func main() {
	length := len(Servers)
	clients := make([]*rpc.Client, 0)
	rates := make([]NetWorkRate.IORates, length)

	for i := 0; i < length; i++ {
		client, err := rpc.DialHTTP("tcp", Servers[i]+":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}
		clients = append(clients, client)
	}

	for {
		wg := new(sync.WaitGroup)
		for i := 0; i < length; i++ {
			wg.Add(1)
			go GetRates(clients[i], wg, &rates[i])
		}
		wg.Wait()
		for i := 0; i < length; i++ {
			fmt.Println(&rates[i])
		}
	}
}
