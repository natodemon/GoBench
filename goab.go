package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Println("Hello World!")

	connsPtr := flag.Int("c", 1, "concurrent requests")
	requestsPtr := flag.Int("n", 1, "requests to make")
	showErrsPtr := flag.Bool("k", false, "show errors count")

	flag.Parse()

}
