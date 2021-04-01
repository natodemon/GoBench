package main

import (
	"flag"
	"os"
	"sync"
)

type ReqResult struct {
	requestTime  float32
	errorOccured bool
}

var testResults = make(chan ReqResult)

func httpWorker(wg *sync.WaitGroup, url string) {
	// Makes simple http request using http.get
	// Measures time for request (find this method) & records in result channel struct
	// Also record if error occured or not
}

func main() {
	concurrentCons := flag.Int("c", 1, "concurrent requests")
	noOfReqs := flag.Int("n", 1, "requests to make")
	showErrsPtr := flag.Bool("k", false, "show errors count")

	posArgs := os.Args[1:]

	flag.Parse()
	if len(posArgs) > 1 {
		// Throw an error
	}

	var hostUrl string = posArgs[0]

	var wg sync.WaitGroup

	// Create a waitgroup with (-c) number of httpworkers
	// Will need form of counting total requests, use basic struct or array/slice
	// httpworkers will continue to make requests until total is exhausted
	// Record total time, calculate average latency and TPS
	// Display this to terminal

}
