package main

import (
	"flag"
	"fmt"
	"github.com/bary321/NetWorkRate"
	"log"
	"net/rpc"
	"sync"
	"time"
)

var (
	configFile = flag.String("c", "./config.json", "配置文件的位置")
)

func GetRates(client *rpc.Client, wg *sync.WaitGroup, rate **NetWorkRate.IORates, interval int) {
	defer wg.Done()
	args := &NetWorkRate.Args{interval}

	divCall := client.Go("Common.GetRate", args, rate, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	// check errors, print, etc.
	if replyCall.Error != nil {
		log.Fatal(replyCall.Error)
	}
	return
}

func GiveSwitch(s string) bool {
	if s == "" {
		return false
	} else {
		return true
	}
}

func main() {
	flag.Parse()

	c := new(Config)
	err := c.Get(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	l := NetWorkRate.NewCustomLogger(c.LogFile, GiveSwitch(c.LogFile), c.Console, 15)

	length := len(c.Servers)
	clients := make([]*rpc.Client, 0)

	rates := new(NetWorkRate.IORates)

	for i := 0; i < length; i++ {
		client, err := rpc.DialHTTP("tcp", c.Servers[i].Ip+":1234")
		log.Println("connect to ", c.Servers[i].Ip)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		clients = append(clients, client)
	}

	for {
		time.Sleep(time.Duration(c.Interval.Client) * time.Second)

		wg := new(sync.WaitGroup)
		for i := 0; i < length; i++ {
			wg.Add(1)
			go GetRates(clients[i], wg, &c.Servers[i].Rates, c.Interval.Server)
		}
		wg.Wait()

		rs := make([]*NetWorkRate.IORate, 0)
		wanRate := new(NetWorkRate.IORate)
		wanRate.Name = "wan"
		rs = append(rs, wanRate)
		lanRate := new(NetWorkRate.IORate)
		lanRate.Name = "lan"
		rs = append(rs, lanRate)
		rates.Rates = rs
		for i := 0; i < length; i++ {
			l := len(c.Servers[i].Rates.Rates)
			for j := 0; j < l; j++ {
				if c.Servers[i].Rates.Rates[j].Name == c.Servers[i].Wan {
					rs[0].Add(c.Servers[i].Rates.Rates[j])
				}
				if c.Servers[i].Rates.Rates[j].Name == c.Servers[i].Lan {
					rs[1].Add(c.Servers[i].Rates.Rates[j])
				}
			}
		}
		fmt.Println()
		l.Println(rates)
	}
}
