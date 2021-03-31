# GoBench
## Purpose

The aim of this project is to understand and demonstrate the effects of parallelising HTTP requests on a set of parameters. 

The project also serves as a way to practice programming in Golang and deploying to a containerised environment.

This repository is an implementation of the specification outlined here: [jig/bench/README.md](https://github.com/jig/bench)

## Task 1: 

### Task 1: Architecture & Setup

The system used for testing contains a 4 core Intel i7 (mobile) chip and is Antergos Linux. The web server is a locally hosted instance of the latest Nginx docker image. The Apache Benchmark image used is the one provided ([jig/docker-ab/README.md](https://github.com/jig/docker-ab)) 

The Nginx server hosts a simple text-only file of around 3kb.

### Task 1: Findings

