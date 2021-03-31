package main

import (
	"flag"
	"net"
	"os"
	"runtime"
	"sync"
	"time"
)

var nThread int
var host, portRange string
var mutex sync.Mutex
// This variable should be protected
var openPortArray []string

func init() {
	flag.IntVar(&nThread, "n", runtime.NumCPU(), "N threads")
	flag.StringVar(&host, "h", "", "target host")
	flag.StringVar(&portRange, "p", "", "port range")
	flag.Parse()

	if host == "" || portRange == "" {
		flag.Usage()
		os.Exit(-1)
	}
}

func DataAverage(nThread int, dataArray []string) [][]string {
	openPortArray := make([][]string, nThread)
	for i := 0; i < nThread; i++ {
		openPortArray[i] = dataArray[(len(dataArray) / nThread) * i: (len(dataArray) / nThread) * (i + 1)]
		// Assign the left data to the last thread
		if i == nThread - 1 {
			openPortArray[i] = dataArray[(len(dataArray) / nThread) * i:]
		}
	}
	return openPortArray
}

func portScanWorker(host string, portRange []string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, _ := range portRange {
		_, err := net.DialTimeout("tcp", host + ":" + portRange[i], time.Second)
		if err == nil {
			mutex.Lock()
			openPortArray = append(openPortArray, portRange[i])
			mutex.Unlock()
		}
	}
}
