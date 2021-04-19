package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bary321/NetWorkRate"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

var (
	dev, ip string
	special bool
	interval int
	server string
)

func init() {
	flag.StringVar(&dev, "dev", "", "")
	flag.StringVar(&server, "server", "172.20.255.3:13498", "")
	flag.StringVar(&ip, "ip", "39.153.163.66", "")
	flag.BoolVar(&special, "special", true, "")
	flag.IntVar(&interval, "interval", 60, "")
}

type Network struct {
	RecRate      float64   `gorm:"column:recRate;not null;COMMENT:'接受网络速度'" form:"recRate"json:"recRate" validate:"required"`
	SentRate      float64   `gorm:"column:sentRate;not null;COMMENT:'发送网络速度'" form:"sectRate"json:"sentRate" validate:"required"`
	Region    int       `gorm:"column:region;not null;type:tinyint(1);index;DEFAULT:1;COMMENT:'地区（1，鄂尔多斯 2，佛山）'" form:"region" json:"region" validate:"oneof=1 2"`
	IP        string    `gorm:"column:ip;index;DEFAULT:1;COMMENT:'ip'" form:"ip" json:"ip" validate:"ip"`
	Card      string    `gorm:"column:card;not null;COMMENT:'网卡'" form:"card"json:"card" validate:"required"`
}

func networkAdd(server, ip, card string, recRate, sentRate float64, region int) error {
	var myClient = &http.Client{Timeout: 120 * time.Second}

	nw := &Network{
		RecRate:   recRate/NetWorkRate.KB,
		SentRate: sentRate/NetWorkRate.KB,
		Region: region,
		IP:     ip,
		Card:   card,
	}

	args := url.Values{}
	var u = url.URL{
		Scheme:   "http",
		Host:     server,
		Path:     "/network/add",
		RawQuery: args.Encode(),
	}

	data, err := json.Marshal(nw)
	if err != nil {
		return err
	}

	r, err := myClient.Post(u.String(), "", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer r.Body.Close()
	out, err := ioutil.ReadAll(r.Body)
	fmt.Println("return", string(out))
	if err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	devs := make([]string, 0, 1)

	if len(dev) > 0 {
		devs = append(devs, dev)
	}
	length := 15
	logger := NetWorkRate.NewCustomLogger("", false, true, length)
	fmt.Println("devs", devs)
	f1, _ := NetWorkRate.IOCountersByFile(special, devs)
	if len(f1.IOCountersStats) == 0 {
		if special {
			fmt.Println("没有找到匹配的网络接口")
		} else {
			fmt.Println("can't find any net devices")
		}
		os.Exit(1)
	}
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		f2, _ := NetWorkRate.IOCountersByFile(special, devs)
		rates, _ := NetWorkRate.GetRate(f1, f2)
		logger.Println(rates)
		err := networkAdd(server, ip, dev, rates.Rates[0].RecvBytesRate, rates.Rates[0].SentBytesRate, 1)
		if err != nil {
			fmt.Println("add network err", err)
		}
		fmt.Println()
		f1 = f2
	}
}


