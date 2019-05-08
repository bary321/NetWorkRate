package main

import (
	"encoding/json"
	"github.com/bary321/NetWorkRate"
	"io/ioutil"
	"os"
)

type Config struct {
	Servers  []*Server     `json:"servers"`
	Default  *DefaultInter `json:"default"`
	Interval *Interval     `json:"interval,omitempty"`
}

type Server struct {
	Ip    string               `json:"ip"`
	Wan   string               `json:"wan,omitempty"`
	Lan   string               `json:"lan,omitempty"`
	Rates *NetWorkRate.IORates `json:"rates, omitempty"` //理论上这个不应该出现在配置文件中
}

type DefaultInter struct {
	Wan string `json:"wan"`
	Lan string `json:"lan"`
}

type Interval struct {
	Server int `json:"server"`
	Client int `json:"client"`
}

func (c *Config) Get(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	d, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	err = json.Unmarshal(d, c)
	if err != nil {
		return err
	}
	length := len(c.Servers)
	for i := 0; i < length; i++ {
		if c.Servers[i].Wan == "" {
			c.Servers[i].Wan = c.Default.Wan
		}
		if c.Servers[i].Lan == "" {
			c.Servers[i].Lan = c.Default.Lan
		}
	}

	return nil

}
