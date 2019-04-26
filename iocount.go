package NetWorkRate

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	NanoToSecond = 1000000000
	devFile      = "/proc/net/dev"
)

//函数的灵感来自于`github.com/shirou/gopsutil`
type IOConterFileStat struct {
	IOCountersStats []IOCountersStat `json:"IoCountersStat"`
	Modtime         time.Time        `json:"modtime"`
}

type IOCountersStat struct {
	Name        string `json:"name"`        // interface name
	BytesSent   uint64 `json:"bytesSent"`   // number of bytes sent
	BytesRecv   uint64 `json:"bytesRecv"`   // number of bytes received
	PacketsSent uint64 `json:"packetsSent"` // number of packets sent
	PacketsRecv uint64 `json:"packetsRecv"` // number of packets received
	Errin       uint64 `json:"errin"`       // total number of errors while receiving
	Errout      uint64 `json:"errout"`      // total number of errors while sending
	Dropin      uint64 `json:"dropin"`      // total number of incoming packets which were dropped
	Dropout     uint64 `json:"dropout"`     // total number of outgoing packets which were dropped (always 0 on OSX and BSD)
	Fifoin      uint64 `json:"fifoin"`      // total number of FIFO buffers errors while receiving
	Fifoout     uint64 `json:"fifoout"`     // total number of FIFO buffers errors while sending
}

type IORates struct {
	Rates []*IORate
}

type IORate struct {
	Name            string `json:"name"`
	SentBytesRate   float64
	RecvBytesRate   float64
	SentPacketsRate float64
	RecvPacketsRate float64
	ErrinRate       float64
	ErroutRate      float64
	DropinRate      float64
	DropoutRate     float64
	FifoinRate      float64
	FifoOutRate     float64
}

func ReadLines(f *os.File) ([]string, error) {
	return ReadLinesOffsetN(f, 0, -1)
}

func ReadLinesOffsetN(f *os.File, offset uint, n int) ([]string, error) {
	var ret []string

	r := bufio.NewReader(f)
	for i := 0; i < n+int(offset) || n < 0; i++ {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		if i < int(offset) {
			continue
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}

	return ret, nil
}

func IOCountersByFile(special bool, devs []string) (*IOConterFileStat, error) {
	fret := new(IOConterFileStat)

	f, err := os.Open(devFile)
	if err != nil {
		fmt.Println("open file err", err)
		return nil, err
	}
	defer f.Close()
	fileinfo, err := f.Stat()
	fret.Modtime = fileinfo.ModTime()
	lines, err := ReadLines(f)
	if err != nil {
		fmt.Println("readlines", err)
		return nil, err
	}

	parts := make([]string, 2)

	for _, line := range lines[2:] {
		separatorPos := strings.LastIndex(line, ":")
		if separatorPos == -1 {
			continue
		}
		parts[0] = line[0:separatorPos]
		parts[1] = line[separatorPos+1:]

		interfaceName := strings.TrimSpace(parts[0])
		if interfaceName == "" || special && !InArray(interfaceName, devs) {
			continue
		}

		fields := strings.Fields(strings.TrimSpace(parts[1]))
		bytesRecv, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return fret, err
		}
		packetsRecv, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return fret, err
		}
		errIn, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			return fret, err
		}
		dropIn, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return fret, err
		}
		fifoIn, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			return fret, err
		}
		bytesSent, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			return fret, err
		}
		packetsSent, err := strconv.ParseUint(fields[9], 10, 64)
		if err != nil {
			return fret, err
		}
		errOut, err := strconv.ParseUint(fields[10], 10, 64)
		if err != nil {
			return fret, err
		}
		dropOut, err := strconv.ParseUint(fields[11], 10, 64)
		if err != nil {
			return fret, err
		}
		fifoOut, err := strconv.ParseUint(fields[12], 10, 64)
		if err != nil {
			return fret, err
		}

		nic := IOCountersStat{
			Name:        interfaceName,
			BytesRecv:   bytesRecv,
			PacketsRecv: packetsRecv,
			Errin:       errIn,
			Dropin:      dropIn,
			Fifoin:      fifoIn,
			BytesSent:   bytesSent,
			PacketsSent: packetsSent,
			Errout:      errOut,
			Dropout:     dropOut,
			Fifoout:     fifoOut,
		}
		fret.IOCountersStats = append(fret.IOCountersStats, nic)
	}
	return fret, nil
}

func Count(data1 uint64, data2 uint64, times float64) (float64, error) {
	tmp := float64(data2 - data1)
	if tmp < 0 {
		return 0, errors.New("data may incorrect")
	} else if tmp == 0 {
		return 0, nil
	}
	if times <= 0 {
		return 0, errors.New("times error")
	}
	return tmp / times, nil
}

func GetRate(stat1 *IOConterFileStat, stat2 *IOConterFileStat) (*IORates, error) {
	var err error = nil
	rates := new(IORates)

	times := float64(stat2.Modtime.UnixNano()-stat1.Modtime.UnixNano()) / NanoToSecond

	for i := 0; i < len(stat1.IOCountersStats); i++ {
		if stat1.IOCountersStats[i].Name != stat2.IOCountersStats[i].Name {
			return nil, errors.New("stat1 not match stat2")
		}

		rate := new(IORate)
		rate.Name = stat1.IOCountersStats[i].Name
		rate.RecvBytesRate, _ = Count(stat1.IOCountersStats[i].BytesRecv, stat2.IOCountersStats[i].BytesRecv, times)
		rate.SentBytesRate, _ = Count(stat1.IOCountersStats[i].BytesSent, stat2.IOCountersStats[i].BytesSent, times)
		rate.SentPacketsRate, _ = Count(stat1.IOCountersStats[i].PacketsSent, stat2.IOCountersStats[i].PacketsSent, times)
		rate.RecvPacketsRate, _ = Count(stat1.IOCountersStats[i].PacketsRecv, stat2.IOCountersStats[i].PacketsRecv, times)
		rate.ErrinRate, _ = Count(stat1.IOCountersStats[i].Errin, stat2.IOCountersStats[i].Errin, times)
		rate.ErroutRate, _ = Count(stat1.IOCountersStats[i].Errout, stat2.IOCountersStats[i].Errout, times)
		rate.DropinRate, _ = Count(stat1.IOCountersStats[i].Dropin, stat2.IOCountersStats[i].Dropin, times)
		rate.DropoutRate, _ = Count(stat1.IOCountersStats[i].Dropout, stat2.IOCountersStats[i].Dropout, times)
		rate.FifoinRate, _ = Count(stat1.IOCountersStats[i].Fifoin, stat2.IOCountersStats[i].Fifoin, times)
		rate.FifoOutRate, _ = Count(stat1.IOCountersStats[i].Fifoout, stat2.IOCountersStats[i].Fifoout, times)
		rates.Rates = append(rates.Rates, rate)
	}

	return rates, err
}

func InArray(tmp string, temp []string) bool {
	for i := 0; i < len(temp); i++ {
		if tmp == temp[i] {
			return true
		}
	}
	return false
}
