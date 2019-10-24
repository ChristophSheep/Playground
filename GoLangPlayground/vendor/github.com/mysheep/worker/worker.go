package worker

import (
	"fmt"
)

type Any interface{}

type Job struct {
	val    Any
	folder string
}

func NewJob(val Any, folder string) *Job {
	return &Job{val, folder}
}

type JobResult struct {
	val Any
}

// NewJobReselt create a new job result
func NewJobResult(val Any) *JobResult {
	return &JobResult{val: val}
}

// Worker is doing the work
// He wait for a jobs, if one is there he takes it
// do the work and send out the result into the results channel
// see https://gobyexample.com/worker-pools
func Worker(id int, jobs <-chan Job, results chan<- JobResult, doWork func(Job) *JobResult) {
	for job := range jobs { // see https://gobyexample.com/range-over-channels
		fmt.Println("Worker", id, "started job")
		result := doWork(job)
		fmt.Println("Worker", id, "finished job")
		results <- *result
	}
}

// SendJobs send some (5) jobs into the job channel
func SendJobs(n int, jobs chan<- Job) {
	for j := 1; j <= n; j++ {
		jobs <- *NewJob(j, ".")
	}
	close(jobs)
}

// GetResults gets the all results from the channel and print it
func GetResults(n int, results <-chan JobResult, quit chan bool) {
	for j := 1; j <= n; j++ {
		res := <-results
		fmt.Printf("Result %d val: %v\n", j, res.val)
	}
	quit <- true
}
