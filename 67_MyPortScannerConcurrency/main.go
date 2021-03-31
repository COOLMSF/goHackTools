package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

func main() {
	var portArray []string
	var startPort, stopPort int
	wg := sync.WaitGroup{}

	n, err := fmt.Sscanf(portRange, "%d-%d", &startPort, &stopPort)
	if n != 2 || err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "port error")
		log.Fatal(err)
	}

	// Generate port array, we need convert integer into ascii, because DataAverage needs asscii
	for i := startPort; i < stopPort; i++ {
		portArray = append(portArray, strconv.Itoa(i))
	}

	// Average data, so that every worker does its own job
	portAverage := DataAverage(nThread, portArray)

	timeNow := time.Now()
	for i := 0; i < nThread; i++ {
		wg.Add(1)
		go portScanWorker(host, portAverage[i], &wg)
	}
	wg.Wait()
	timeEscape := time.Since(timeNow)

	sort.Strings(openPortArray)
	fmt.Println("Open ports:")
	fmt.Println(openPortArray)
	fmt.Printf("Time:%v\n", timeEscape)
}
