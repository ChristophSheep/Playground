package main

import (
	"math/rand"
	"time"

	"github.com/mysheep/worker"
)

func doWork(job worker.Job) *worker.JobResult {
	time.Sleep(time.Duration(1+rand.Intn(2)) * time.Second)
	return worker.NewJobResult(1) // TODO
}

func main() {

	const N = 2

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
	go worker.SendJobs(jobs)

	// Receive results of the workers
	//
	go worker.GetResults(results, quit)

	<-quit
}
