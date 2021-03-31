package main

import (
	"fmt"
	"github.com/amoghe/go-crypt"
	"log"
	"os"
	"strings"
	"sync"
)

func dataAverage(nThread int, dataArray []string) [][]string {
	resultArray := make([][]string, nThread)
	for i := 0; i < nThread; i++ {
		resultArray[i] = dataArray[(len(dataArray) / nThread) * i: (len(dataArray) / nThread) * (i + 1)]
		// Assign the left data to the last thread
		if i == nThread - 1 {
			resultArray[i] = dataArray[(len(dataArray) / nThread) * i:]
		}
	}
	return resultArray
}

func testPass(encryptedPass string, passInDict string) bool {
	// Get salt
	lastIndex := strings.LastIndex(encryptedPass, "$")
	salt := encryptedPass[0:lastIndex]
	// Get sha512
	sha512, err := crypt.Crypt(passInDict, salt)
	if err != nil {
		log.Fatal(err)
	}

	if sha512 == encryptedPass {
		return true
	}
	return false
}

func crackShadowConcurrency(shadowFile string, dictFile string, nThread int) {
	wg := sync.WaitGroup{}

	shadowFileFd, err := os.Open(shadowFile)
	if err != nil {
		log.Fatal(err)
	}
	defer shadowFileFd.Close()

	// We want get size of this file, so that we can create buf of that size to hold the data
	shadowFileInfo, err := shadowFileFd.Stat()
	if err != nil {
		log.Fatal(err)
	}
	shadowfileSize := shadowFileInfo.Size()

	dictFileFd, err := os.Open(dictFile)
	if err != nil {
		log.Fatal(err)
	}
	defer dictFileFd.Close()

	dictFileInfo, err := dictFileFd.Stat()
	if err != nil {
		log.Fatal(err)
	}
	dictFileSize := dictFileInfo.Size()

	shadowBuf := make([]byte, shadowfileSize)
	_, err = shadowFileFd.Read(shadowBuf)
	if err != nil {
		log.Fatal(err)
	}

	dictBuf := make([]byte, dictFileSize)
	_, err = dictFileFd.Read(dictBuf)
	if err != nil {
		log.Fatal(err)
	}

	// We need to split shadow file
	shadowArray := strings.Split(string(shadowBuf), "\n")
	dictArray := strings.Split(string(dictBuf), "\n")

	averageShadowArray := dataAverage(nThread, shadowArray)
	averageDictArray := dataAverage(nThread, dictArray)

	// Concurrency
	for i :=0; i < nThread; i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			// Repoint to specified data
			shadowArray = averageShadowArray[i]
			dictArray = averageDictArray[i]

			// try every account
			for _, shadow := range shadowArray {
				// Handle null string, skip it
				if shadow == "" {
					continue
				}

				// If no encrypted password found, encrypted password must contain '$'
				if !strings.Contains(shadow, "$") {
					continue
				}

				// Data before the first ':' is username, otherwise is encrypted password
				data := strings.Split(shadow, ":")
				username := data[0]
				encryptedPass := data[1]

				// try every password
				for _, dict := range dictArray {
					if dict == "" {
						continue
					}
					if data != nil {
						if testPass(string(encryptedPass), dict) {
							fmt.Println("[+] Found PASSWORD")
							fmt.Println("Username:" + string(username))
							fmt.Println("Password:" + string(dict))
						}
					}
				}
			}
		}()
	}
	wg.Wait()
}
