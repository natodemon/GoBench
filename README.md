# GoBench
## Purpose

The aim of this project is to understand and demonstrate the effects of parallelising HTTP requests on a set of parameters. 

The project also serves as a way to practice programming in Golang and deploying to a containerised environment.

This repository is an implementation of the specification outlined here: [jig/bench/README.md](https://github.com/jig/bench)

## Task 1: Nginx & Apache Benchmark

### Architecture & Setup

The system used for testing contains a 4 core Intel i7 (mobile) chip, 8GB of memory and is running Antergos Linux. The web server is a locally hosted instance of the latest Nginx docker image. The Apache Benchmark image used is the one provided ([jig/docker-ab/README.md](https://github.com/jig/docker-ab)) 

The Nginx server hosts a simple text-only file of around 3kb.

### Findings

Concurrent connections decrease latency and increase TPS (Transactions per second) up to a certain limit. 


| Conns | TPS | Avg Latency | Std Dev | Max Latency |
| :---: | :-: | :---------: | :-----: | :---------: |
| 20    |     |             |         |             |
| 30    |     |             |         |             |
| 50    |     |             |         |             |
| 80    |     |             |         |             |