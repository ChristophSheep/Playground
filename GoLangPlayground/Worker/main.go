package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mysheep/worker"
)

func doWork(job worker.Job) *worker.JobResult {
	// Simulate some work
	time.Sleep(time.Duration(1+rand.Intn(1)) * time.Second)
	//
	return worker.NewJobResult(1) // TODO
}

func main() {

	const N = 2
	const J = 5

	// Create channels
	//
	quit := make(chan bool)
	jobs := make(chan worker.Job, N)
	results := make(chan worker.JobResult, N)

	// Create workers
	//
	const W = 3
	for i := 0; i < W; i++ {
		go worker.Worker(i, jobs, results, doWork)
	}

	// Create some jobs for the workers
	//
	go worker.SendJobs(J, jobs)

	// Receive results of the workers
	//
	go worker.GetResults(J, results, quit)

	<-quit
	fmt.Println("The workers finished all the work .. Bye")
}
