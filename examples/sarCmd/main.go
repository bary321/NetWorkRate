package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bary321/NetWorkRate"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

//test
var (
	setFirst   = false
	setLast    = false
	first      = new(NetWorkRate.IOConterFileStat)
	last       = new(NetWorkRate.IOConterFileStat)
	timeFormat = "06-01-02 15:04:05"

	devs     = make([]string, 4)
	interval = flag.Int("n", 1, "刷新的时间间隔")
)

func PrintFirstWithPrefix(average bool, length int, t time.Time) {
	printPrefix(average, t)
	NetWorkRate.FirstLine(length)
}

func PrintlineWithPrefix(average bool, length int, t time.Time, rate *NetWorkRate.IORate) {
	printPrefix(average, t)
	NetWorkRate.LinesPrint(length, rate)
}

func printPrefix(average bool, t time.Time) {
	if average {
		fmt.Print("Average: ")
	}
	fmt.Print(t.Format(timeFormat), " ")
}

type sliceValue []string

func newSliceValue(vals []string, p *[]string) *sliceValue {
	*p = vals
	return (*sliceValue)(p)
}

func (s *sliceValue) Set(val string) error {
	*s = sliceValue(strings.Split(val, ","))
	return nil
}

func (s *sliceValue) Get() interface{} { return []string(*s) }

func (s *sliceValue) String() string { return strings.Join([]string(*s), ",") }

func main() {

	flag.Var(newSliceValue([]string{}, &devs), "d", "显示哪几个网络设备，以“,”号隔开")

	flag.Parse()

	special := false
	if len(devs) != 0 {
		special = true
	}
	length := 15

	sigs := make(chan os.Signal, 1)
	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// This goroutine executes a blocking receive for
	// signals. When it gets one it'll print it out
	// and then notify the program that it can finish.
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println("sig:", sig)
		if !setFirst || !setLast || first == nil || last == nil {
			os.Exit(0)
		}
		rates, _ := NetWorkRate.GetRate(first, last)
		t := time.Now()
		PrintFirstWithPrefix(true, length, t)
		for i := 0; i < len(rates.Rates); i++ {
			PrintlineWithPrefix(true, length, t, rates.Rates[i])
		}
		os.Exit(0)
	}()

	f1, _ := NetWorkRate.IOCountersByFile(special, devs)
	if len(f1.IOCountersStats) == 0 {
		if special {
			fmt.Println("没有找到匹配的网络接口")
		} else {
			fmt.Println("can't find any net devices")
		}
		os.Exit(1)
	}
	if setFirst == false {
		tmp, _ := json.Marshal(f1)
		_ = json.Unmarshal(tmp, first)
		setFirst = true
	}
	for {
		time.Sleep(time.Duration(*interval) * time.Second)
		f2, _ := NetWorkRate.IOCountersByFile(special, devs)
		last = f2
		if !setLast {
			setLast = true
		}
		rates, _ := NetWorkRate.GetRate(f1, f2)
		t := time.Now()
		PrintFirstWithPrefix(false, length, t)
		for i := 0; i < len(rates.Rates); i++ {
			PrintlineWithPrefix(false, length, t, rates.Rates[i])
		}

		fmt.Println()
		f1 = f2
	}
}
