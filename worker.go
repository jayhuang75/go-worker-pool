package worker

import (
	"fmt"
	"sync"
	"time"
)

// ProcessorFunc signature that defines the dependency injection to process "Jobs"
// Your custom function need to be done with workers
type ProcessorFunc func(resource interface{}) error

// ResultProcessorFunc signature that defines the dependency injection to process "Results"
// We capture the result - error and display independently
type ResultProcessorFunc func(result Result) error

// Job Structure that wraps Jobs information
type Job struct {
	id       int
	resource interface{}
}

// Result holds the main structure for worker processed job results.
type Result struct {
	Job Job
	Err error
}

// Pool Manager generic struct that keeps all the logic to manage the queues
type Pool struct {
	numRoutines int
	jobs        chan Job
	results     chan Result
	done        chan bool
	completed   bool
}

// NewWorkerPool NewManager returns a new manager structure ready to be used.
func NewWorkerPool(numRoutines int) *Pool {
	fmt.Printf("[Worker Pool] Creating a new Pool")
	r := &Pool{numRoutines: numRoutines}
	r.jobs = make(chan Job, numRoutines)
	r.results = make(chan Result, numRoutines)
	return r
}

// Start func
func (p *Pool) Start(resources []interface{}, procFunc ProcessorFunc, resFunc ResultProcessorFunc) {
	startTime := time.Now()
	fmt.Printf("[Worker Pool] Starting at: %s\n", startTime)
	go p.allocate(resources)
	p.done = make(chan bool)
	go p.collectResult(resFunc)
	go p.workerPool(procFunc)
	<-p.done // Stop all the job
	fmt.Printf("[Worker Pool] Total time taken: [%f] seconds\n", time.Since(startTime).Seconds())
}

// allocate allocates jobs based on an array of resources to be processed by the worker pool
func (p *Pool) allocate(jobs []interface{}) {
	defer close(p.jobs)
	fmt.Printf("[Worker Pool] Allocating [%d] resources\n", len(jobs))
	for i, v := range jobs {
		job := Job{id: i, resource: v}
		p.jobs <- job
	}
	fmt.Printf("[Worker Pool] Done Allocation.\n")
}

// work performs the actual work by calling the processor and passing in the Job as reference obtained
// from iterating over the "Jobs" channel
func (p *Pool) doWork(i int, wg *sync.WaitGroup, processor ProcessorFunc) {
	defer wg.Done()
	fmt.Printf("[Worker Pool] Worker [%d] starting\n ", i)
	for job := range p.jobs {
		fmt.Printf("[Worker Pool] Worker [%d] working on Job ID [%d]\n", i, job.id)
		output := Result{job, processor(job.resource)}
		p.results <- output
		fmt.Printf("[Worker Pool] Worker [%d] done with Job ID [%d]\n", i, job.id)
	}
	fmt.Printf("[Worker Pool] Worker [%d] done.\n", i)
}

// workerPool creates or spawns new "work" goRoutines to process the "Jobs" channel
func (p *Pool) workerPool(processor ProcessorFunc) {
	defer close(p.results)
	fmt.Printf("[Worker Pool] Worker Pool spawning new GoRoutine, total: [%d]\n", p.numRoutines)
	var wg sync.WaitGroup
	for i := 0; i < p.numRoutines; i++ {
		wg.Add(1)
		go p.doWork(i, &wg, processor)
		fmt.Printf("[Worker Pool] Spawned work GoRoutine [%d]\n", i)
	}
	fmt.Printf("[Worker Pool] Done spawning work total [%d] GoRoutine\n", p.numRoutines)
	wg.Wait()
	fmt.Printf("[Worker Pool] All worker GoRoutine done processing\n")

}

// Collect post processes the channel "Results" and calls the ResultProcessorFunc passed in as reference
// for further processing.
func (p *Pool) collectResult(proc ResultProcessorFunc) {
	fmt.Printf("[Worker Pool] GoRoutine collect starting\n")
	for result := range p.results {
		outcome := proc(result)
		fmt.Printf("[Worker Pool] Job with id: [%d] completed, outcome: %s\n", result.Job.id, outcome)
	}
	fmt.Printf("[Worker Pool] GoRoutine collect done, setting channel done as completed\n")
	p.done <- true
	p.completed = true
}

// IsCompleted utility method to check if all work has done from an outside caller.
func (p *Pool) IsCompleted() bool {
	return p.completed
}
