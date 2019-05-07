package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestStruct(t *testing.T) {
	c := new(Config)
	f, err := os.Open("./example.json")
	if err != nil {
		t.Error(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%s\n", b)
	fmt.Println(c)
}

func TestGet(t *testing.T) {
	c := new(Config)
	err := c.Get("./example.json")
	if err != nil {
		t.Error(err)
	}

	length := len(c.Servers)
	for i := 0; i < length; i++ {
		if c.Servers[i].Ip == "192.168.2.90" {
			if c.Servers[i].Wan != "eth1" || c.Servers[i].Lan != "eth0" {
				t.Error("get function fail")
				t.Error(c.Servers[i])
			}
		}
		if c.Servers[i].Ip == "192.168.2.91" {
			if c.Servers[i].Wan != "ens32" {
				t.Error("get function fail")
				t.Error(c.Servers[i])
			}
		}
	}

}
