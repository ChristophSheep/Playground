package main

import (
	"fmt"
	"time"
)

type Job struct {
	val    int
	folder string
}

type JobResult struct {
	val int
}

// Worker is doing the work
// see https://gobyexample.com/worker-pools
func Worker(id int, jobs <-chan Job, results chan<- JobResult, doWork func(Job) JobResult) {
	for job := range jobs { // see https://gobyexample.com/range-over-channels
		fmt.Println("worker", id, "started job")
		result := doWork(job)
		fmt.Println("worker", id, "finished job")
		results <- result
	}
}

func doWork(job Job) JobResult {
	time.Sleep(1 * time.Second)
	return JobResult{val: job.val * 1}
}

// SendJobs send some job into the job channel
func SendJobs(jobs chan<- Job) {
	for j := 1; j <= 5; j++ {
		jobs <- Job{val: j, folder: ""}
	}
	close(jobs)
}

// GetResults results from the channel and print it
func GetResults(results <-chan JobResult, quit chan bool) {
	for j := 1; j <= 5; j++ {
		res := <-results
		fmt.Println("result", res.val)
	}
	quit <- true
}

func main() {

	const N = 2

	// Create channels
	//
	quit := make(chan bool)
	jobs := make(chan Job, N)
	results := make(chan JobResult, N)

	// Create workers
	//
	const W = 3
	for i := 0; i < W; i++ {
		go Worker(i, jobs, results, doWork)
	}

	// Create some jobs for the workers
	//
	go SendJobs(jobs)

	// Receive results of the workers
	//
	go GetResults(results, quit)

	<-quit
}
