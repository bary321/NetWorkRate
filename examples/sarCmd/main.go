package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bary321/NetWorkRate"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	setFirst   = false
	setLast    = false
	first      = new(NetWorkRate.IOConterFileStat)
	last       = new(NetWorkRate.IOConterFileStat)
	timeFormat = "06-01-02 15:04:05"

	interval = flag.Int("i", 1, "刷新的时间间隔")
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

func main() {

	flag.Parse()

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

	f1, _ := NetWorkRate.IOCountersByFile(false, "")
	if setFirst == false {
		tmp, _ := json.Marshal(f1)
		_ = json.Unmarshal(tmp, first)
		setFirst = true
	}
	for {
		time.Sleep(time.Duration(*interval) * time.Second)
		f2, _ := NetWorkRate.IOCountersByFile(false, "")
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
