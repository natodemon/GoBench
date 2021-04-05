# GoBench
## Purpose

The aim of this project is to understand and demonstrate the effects of parallelising HTTP requests on a set of parameters. 

The project also serves as a way to practice programming in Golang and deploying to a containerised environment.

This repository is an implementation of the specification outlined here: [jig/bench/README.md](https://github.com/jig/bench)

## Task 1: Nginx & Apache Benchmark

### Architecture & Setup

The system used for testing contains a 4 core Intel i7 (mobile) chip, 8GB of memory and is running Antergos Linux. The web server is a locally hosted instance of the latest Nginx docker image. The Apache Benchmark image used is the one provided ([jig/docker-ab/README.md](https://github.com/jig/docker-ab)) 

The Nginx server hosts a text-only HTML file of around 130kb.

### Findings

The below table shows the results of running a number of tests with the above setup. All tests were carried out using 100 000 requests (-n 100000) with the only change being to the concurrency level (-c). The CPU usage was also monitored for the duration of each run and will be commented on below. All values except TPS and Conns are in milliseconds.

| Conns | TPS  | Avg Latency | Std Dev | Max Latency |
| :---: | :--: | :---------: | :-----: | :---------: |
| 1     | 3351 | 0.30        | 0.1     | 7           |
| 20    | 6545 | 3.00        | 1.6     | 41          |
| 30    | 6578 | 4.56        | 2.3     | 57          |
| 50    | 6977 | 7.17        | 3.5     | 109         |
| 100   | 7345 | 13.6        | 6.8     | 196         |

As an initial conclusion we can see that more concurrent connections increases the TPS (Transactions per second) value, but negatively impacts the latency as this also increases. Increasing the concurrency has quickly diminishing returns, 20 and 50 simultaneous connections have TPS values within a few hundred transactions but these transactions take on average more than twice as long. The maximum latency also increases greatly with higher concurrency. This could cause user experience or more serious problems depending on the use case.

From the second test (20 connections) CPU usage was already around 95% on 3 cores and slightly lower on the 4th. This utilisation did not increase much with the subsequent tests and it is for this reason that I believe the TPS figures do not increase significatly after 20 connections. However, a failed request was not detected at any point during the tests so I believe I have yet to push the system to it's performance limit. 

Having performed more granular tests, I believe the optimum concurrency (-c) value for this system and request content is 13. Average CPU utilisation is around 90% while returning a TPS of 6250 and an average request latency of 2 ms.

## Task 2: Go Benchmark (goab)

### Implementation

The aim of this task was to replicate the essential functionality of Apache Bench with an equivalent program written in Go. This program would perform the same tests as ab and monitor the same parameters to allow for a comparison of the two solutions.

### Findings & Comparison

