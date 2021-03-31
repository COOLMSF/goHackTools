package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	var nThread int
	var dictFile, shadowFile string

	flag.StringVar(&dictFile, "d", "", "dictionary file")
	flag.StringVar(&shadowFile, "f", "", "shadow file")
	flag.IntVar(&nThread, "n", runtime.NumCPU(), "n threads")
	flag.Parse()

	if dictFile == "" || shadowFile == "" {
		flag.Usage()
		os.Exit(-1)
	}

	timeStart := time.Now()
	crackShadowConcurrency(shadowFile, dictFile, nThread)
	timeSince := time.Since(timeStart)

	fmt.Printf("Time:%v", timeSince)
}
