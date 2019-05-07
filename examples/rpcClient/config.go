package main

type Config struct {
	Servers []*Server     `json:"servers"`
	Default *DefaultInter `json:"default"`
}

type Server struct {
	Ip string `json:"ip"`
}

type DefaultInter struct {
	Wan string `json:"wan"`
	Lan string `json:"lan"`
}
