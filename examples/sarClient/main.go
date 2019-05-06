package main

import (
	"fmt"
	"github.com/bary321/NetWorkRate"
	"github.com/gin-gonic/gin"
	"sync"
	"time"
)

var (
	rateNow = new(NetWorkRate.IORates)
	mutex   = new(sync.RWMutex)
	one     = new(NetWorkRate.IOConterFileStat)
	next    = new(NetWorkRate.IOConterFileStat)
)

func init() {
	one, _ = NetWorkRate.IOCountersByFile(false, nil)

	time.Sleep(time.Second)

	next, _ = NetWorkRate.IOCountersByFile(false, nil)
	mutex.Lock()
	rateNow, _ = NetWorkRate.GetRate(one, next)
	defer mutex.Unlock()
}

func GetRate(c *gin.Context) {
	mutex.RLock()
	defer mutex.RUnlock()
	c.JSON(200, rateNow)
}

func main() {
	r := gin.Default()
	r.GET("/", GetRate)
	go func() {
		one = next
		time.Sleep(time.Second)
		next, _ = NetWorkRate.IOCountersByFile(false, nil)
		mutex.Lock()
		rateNow, _ = NetWorkRate.GetRate(one, next)
		defer mutex.Unlock()
	}()
	if err := r.Run(); err != nil {
		fmt.Println("gin run err")
	}
}
