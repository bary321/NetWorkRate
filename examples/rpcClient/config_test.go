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

	if c.Servers[0].Wan != "eth1" || c.Servers[0].Lan != "eth0" {
		t.Error("get function fail")
		t.Error(c.Servers[0])
	}
}
