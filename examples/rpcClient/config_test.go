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
	f, err := os.Open("./config.json")
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
