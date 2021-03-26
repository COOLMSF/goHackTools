package main

import (
    "flag"
    "fmt"
    "os"
    "sync"
)

var allUrls = make([]urlInfo, 0x1)

func main() {
    urlPtr := flag.String("url", "", "url to claw")
    nThreadPtr := flag.Int("nThread", 1, "N threads to claw")
    nDepthPtr := flag.Int("nDepth", 1, "max depth to claw")
    flag.Parse()

    url := *urlPtr
    nDepth := *nDepthPtr
    nThread := *nThreadPtr

    // Url must be specified
    if url == "" {
        flag.Usage()
        os.Exit(-1)
    }

    allUrls[0].url = url
    // Scan all urls, use n thread
    var wg sync.WaitGroup
    for n := 0; n < nThread; n++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for i := 1; i <= nDepth; i++ {
                for _, url := range allUrls {
                    readUrl(url)
                }
                // Filter unnecessary data
                allUrls = urlFilter()
            }
        }()
    }
    wg.Wait()

    for _, url := range allUrls {
        fmt.Println(url.url)
    }
}
