# GoBench
## Purpose

The aim of this project is to understand and demonstrate the effects of parallelising HTTP requests on a set of parameters. 

The project also serves as a way to practice programming in Golang and deploying to a containerised environment.

This repository is an implementation of the specification outlined here: [jig/bench/README.md](https://github.com/jig/bench)

## Usage

Usage is very similar to Apache Bench, however see below for a few example executions.

    go run goab.go [-k -c=connections -n=requests] url
    ./goab [-k -c=connections -n=requests] url

Example:

    go run goab.go -k -c=20 -n=50000 http://127.17.0.1:8080


## Task 1: Nginx & Apache Benchmark

### Architecture & Setup

The system used for testing contains a 4 core Intel i7 (mobile) chip, 8GB of memory and is running Antergos Linux. The web server is a locally hosted instance of the latest Nginx docker image. The Apache Benchmark image used is the one provided ([jig/docker-ab/README.md](https://github.com/jig/docker-ab)) 

The Nginx server hosts a simple text-only HTML file of around 3.5kb.

### Findings

The below table shows the results of running a number of tests with the above setup. All tests were carried out using 50 000 requests (-n 50000) with the only change being to the concurrency level (-c). The CPU usage was also monitored for the duration of each run and will be commented on below. All values except TPS and Conns are in milliseconds.

| Conns | TPS   | Avg Latency | Std Dev | Max Latency |
| :---: | :---: | :---------: | :-----: | :---------: |
| 1     | 4638  | 0.22        | 0.1     | 4           |
| 20    | 10640 | 1.88        | 1.1     | 29          |
| 30    | 10491 | 2.86        | 1.5     | 32          |
| 50    | 10938 | 4.57        | 2.3     | 38          |
| 100   | 10802 | 9.26        | 5.4     | 131         |

As an initial conclusion we can see that more concurrent connections increases the TPS (Transactions per second) value, but negatively impacts the latency as this also increases. Increasing the concurrency has quickly diminishing returns, 20 and 50 simultaneous connections have TPS values within a few hundred transactions but these transactions take on average more than twice as long. The maximum latency also increases greatly with higher concurrency. This could cause user experience or more serious problems depending on the use case.

From the second test (20 connections) CPU usage was already around 95% on 3 cores and slightly lower on the 4th. This utilisation did not increase much with the subsequent tests and it is for this reason that I believe the TPS figures do not increase significatly after 20 connections. However, a failed request was not detected at any point during the tests so I believe I have yet to push the system to it's performance limit. 

Having performed more granular tests, I believe the optimum concurrency (-c) value for this system and request content is 13. Average CPU utilisation is around 90% while returning a TPS of 6250 and an average request latency of 2 ms.

## Task 2: Go Benchmark (goab)

### Implementation

The aim of this task was to replicate the essential functionality of Apache Bench with an equivalent program written in Go. This program would perform the same tests as ab and monitor the same parameters to allow for a comparison of the two solutions.

Given golang's architecture, namely it's *goroutines* and *channels*, it is relatively easy to make functions asynchronous and execute them concurrently. The approach taken with goab was to implement the http request within a goroutine and thus have the ability to run multiple requests concurrently. Each goroutine runs in it's own thread and runs until the entered number of requests is reached.

### Observations & Comparison

As an initial analysis, the program runs almost on par, if slighly less performant than Apache Bench. The performance follows a similar curve to ab, with the TPS almost doubling as concurrent connections are introduced but quickly flattening out as concurrency is increased. The latency figures are also similar with average latency increasing with concurrency.

| Conns | TPS   | Avg Latency | Max Latency |
| :---: | :---: | :---------: | :---------: |
| 1     | 5024  | 0.20        |    5.9      |
| 20    | 9785  | 2.04        |    24       |
| 30    | 9491  | 3.15        |    28       |
| 50    | 8443  | 5.91        |    54       |
| 100   | 7053  | 14.2        |    86       | * Run with 25 000 requests

As can be seen above, the TPS figures begin to drop when using a concurrency level of 50 or more, I believe this is due to reaching the performance limit of the test machine but also potentially inefficiencies in the implementation. Using -c 100 I regularly ran into stability issues and the program often crashed, hence the reduced number of requests. Given that each connection runs as a separate thread, I attempted to test the program on a machine with more cores but was unsuccessful in setting up a test environment. 

Although I did run into stability issues under certain circumstances, I did monitor for failed requests and for all successful tests the number of errored responses was always zero.

## Task 3: Go HTTP Server

## Implementation

Please see [natodemon/GoHTTP](https://github.com/natodemon/GoHTTP) for the http server implementation and usage.

## Observations & Comparison

Running tests in the same environment and under the same conditions as the previous tests yields some interesting results. The simple golang http server returns much higher TPS values and lower latency figures than the nginx server previously tested. Predictably, the apache bench program returns better results, however both testing programs have markedly better performance. The below table shows the figures for a couple of test runs of 50 000 requests:

| Conns | TPS goab | TPS ab | Latency goab | Latency ab  |
| :---: | :------: | :----: | :----------: | :---------: |
| 1     | 6249     |  8864  | 0.16         |  0.11       |
| 20    | 14135    | 22415  | 1.41         |  0.89       |
| 50    | 12564    | 22248  | 3.97         |  2.25       |
| 100   | 9760     | 21961  | 10.2         |  4.55       |

The overall performance follows a similar curve to the results using nginx, TPS increases with concurrent requests but quickly reaches a limit. I believe the improved performance is owed to the lower overhead of the Go httpserver, given it's simplicity and single use-case but also largely due to use of Goroutines. Internally, each http request is served within a goroutine and thus is performed on an individual thread. This allows for improved performance, especially on systems with high core counts.

## Conclusion

Overall, this project has been a success in it's goal as an exercise to better understand and gain some experience coding Golang. While the Go alternative to Apache Bench, goab did not prove to be as performant as originally hoped, the go httpserver most definitely proved the advantages of Go's architecture. Given more time, I would go back and attempt to reduce a lot of the overhead involved in the management of the goroutines in goab. Having seen the efficiency of the httpserver, I believe I could eventually improve it's performance and in turn, it's stability issues.