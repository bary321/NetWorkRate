package main

import (
	"encoding/json"
	"fmt"
	"github.com/bary321/NetWorkRate"
	"net/http"
	"time"
)

func main() {
	url := "http://127.0.0.1:8080"
	rates := new(NetWorkRate.IORates)
	tmp := make([]byte, 2096)
	for {
		res, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		if _, err = res.Body.Read(tmp); err != nil {
			fmt.Println("read err")
		}
		fmt.Println(string(tmp))
		if err = json.Unmarshal(tmp, rates); err != nil {
			fmt.Println("unmarshal err")
		}

		break
		time.Sleep(time.Second)
	}
}
