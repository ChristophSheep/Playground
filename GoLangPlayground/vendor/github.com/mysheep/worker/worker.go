package worker

import (
	"fmt"
)

type Any interface{}

type Job struct {
	val    Any
	folder string
}

type JobResult struct {
	val Any
}

func NewJobResult(val Any) *JobResult {
	return &JobResult{val: val}
}

// Worker is doing the work
// see https://gobyexample.com/worker-pools
func Worker(id int, jobs <-chan Job, results chan<- JobResult, doWork func(Job) *JobResult) {
	for job := range jobs { // see https://gobyexample.com/range-over-channels
		fmt.Println("Worker", id, "started job")
		result := doWork(job)
		fmt.Println("Worker", id, "finished job")
		results <- *result
	}
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
		fmt.Printf("Result %d val: %v\n", j, res.val)
	}
	quit <- true
}
