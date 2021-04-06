package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type ReqResult struct {
	requestTime  time.Duration
	errorOccured bool
	reqIndex     int
}

var resultChan = make(chan ReqResult, 200)
var requestsChan = make(chan int, 200)
var reqUrl string
var totalReqs int

func httpWorker(wg *sync.WaitGroup, id int) {
	for curReq := range requestsChan {
		var resRecord ReqResult

		tempStart := time.Now()
		res, reqErr := http.Get(reqUrl)
		if reqErr != nil {
			resRecord.errorOccured = true
			wg.Done()
		}

		body, _ := io.ReadAll(res.Body)
		if body == nil {
			println("")
		}
		res.Body.Close()

		timeEnd := time.Now()
		resRecord.requestTime = timeEnd.Sub(tempStart)
		resRecord.reqIndex = curReq
		resultChan <- resRecord
	}
	wg.Done()
}

func parseResults(done chan bool, showErrors bool) {
	var errCount int
	var maxLatency time.Duration
	var latencySum time.Duration

	for reqInf := range resultChan {
		latencySum += reqInf.requestTime
		if reqInf.errorOccured {
			errCount++
		}
		if reqInf.requestTime > maxLatency {
			maxLatency = reqInf.requestTime
		}
	}

	avgLatency := latencySum / time.Duration(totalReqs)

	fmt.Printf("Avg latency over %v requests: %v \nMax latency: %v \n", totalReqs, avgLatency, maxLatency)
	if showErrors {
		fmt.Println("Total error count:", errCount)
	}

	done <- true
}

func allocReqs(reqs int) {
	for i := 0; i < reqs; i++ {
		reqIndex := i
		requestsChan <- reqIndex
	}
	close(requestsChan)
}

func main() {
	concurrentCons := flag.Int("c", 1, "concurrent requests")
	noOfReqs := flag.Int("n", 1, "requests to make")
	showErrsPtr := flag.Bool("k", false, "show errors count")

	flag.Parse()

	positionalArg := os.Args[len(os.Args)-1]

	if len(positionalArg) < 1 {
		log.Fatal(errors.New("no URL input"))
	}

	reqUrl = positionalArg
	totalReqs = *noOfReqs

	go allocReqs(*noOfReqs)

	done := make(chan bool)
	go parseResults(done, *showErrsPtr)

	var wg sync.WaitGroup
	mainStart := time.Now()
	for i := 0; i < *concurrentCons; i++ {
		wg.Add(1)
		go httpWorker(&wg, i)
	}
	wg.Wait()
	close(resultChan)

	<-done
	endTime := time.Since(mainStart)

	fmt.Println("Total time:", endTime)
	fmt.Printf("TPS: %.2f \n", (float64(totalReqs) / endTime.Seconds()))
}
