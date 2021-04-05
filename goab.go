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
	// Makes simple http request using http.get
	// Measures time for request (find this method) & records in result channel struct
	// Also record if error occured or not

	for curReq := range requestsChan {
		//for len(resultChan) <= totalReqs {

		var resRecord ReqResult

		tempStart := time.Now()
		res, reqErr := http.Get(reqUrl)
		if reqErr != nil {
			resRecord.errorOccured = true
			wg.Done()
		}

		body, _ := io.ReadAll(res.Body)
		if body == nil {
			println("This is pointless")
		}
		res.Body.Close()
		timeEnd := time.Now()

		resRecord.requestTime = timeEnd.Sub(tempStart)
		resRecord.reqIndex = curReq

		//fmt.Printf("WorkerID: %d | Request No: %d \n", id, curReq)
		//fmt.Println("Request number:", curReq)
		//fmt.Println("Req time:", resRecord.requestTime)

		resultChan <- resRecord

	}
	wg.Done()
}

func parseResults(done chan bool, showErrors bool) {
	var errCount int = 0
	var latencySum time.Duration

	for reqInf := range resultChan {
		latencySum += reqInf.requestTime
		if reqInf.errorOccured {
			errCount++
		}
	}

	avgLatency := latencySum / time.Duration(totalReqs)

	fmt.Println("Avg latency (ms) over", totalReqs, ":", avgLatency)
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
	//println("Allocation complete")
	close(requestsChan)
} // Creates channel of 'requests' that httpworkers consume

func main() {
	concurrentCons := flag.Int("c", 1, "concurrent requests")
	noOfReqs := flag.Int("n", 1, "requests to make")
	showErrsPtr := flag.Bool("k", false, "show errors count")

	flag.Parse()

	posArgs := os.Args[len(os.Args)-1]

	//fmt.Println(posArgs)

	if len(posArgs) < 1 {
		log.Fatal(errors.New("no URL input"))
	}

	reqUrl = posArgs
	totalReqs = *noOfReqs

	//println("Allocation started")
	go allocReqs(*noOfReqs)

	// concRequests := make(chan struct{}, *noOfReqs)
	// for i := 0; i < *noOfReqs; i++ {

	// }

	done := make(chan bool)
	var wg sync.WaitGroup

	go parseResults(done, *showErrsPtr)

	mainStart := time.Now()

	for i := 0; i < *concurrentCons; i++ {
		wg.Add(1)
		//fmt.Printf("Worker %d initiated\n", i)
		go httpWorker(&wg, i)
	}
	wg.Wait()
	//fmt.Println("Closing channel")
	close(resultChan)

	//go parseResults(done)
	<-done
	endTime := time.Since(mainStart)

	fmt.Println("Total time:", endTime)
	fmt.Printf("TPS: %.2f \n", (float64(totalReqs) / endTime.Seconds()))

	// Create a waitgroup with (-c) number of httpworkers
	// Will need form of counting total requests, use basic struct or array/slice
	// httpworkers will continue to make requests until total is exhausted
	// Record total time, calculate average latency and TPS
	// Display this to terminal

}
