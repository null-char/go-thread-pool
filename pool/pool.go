package pool

import (
	"fmt"
	"log"
	"sync"
	"time"

	j "github.com/null-char/go-thread-pool/job"
)

// Pool describes the basic information associated with a worker pool.
type Pool struct {
	numWorkers uint
	jobs       chan j.Job
	results    chan j.Result
	done       chan bool
}

// New creates a new pool with the specified number of workers.
func New(numWorkers uint) *Pool {
	p := &Pool{numWorkers: numWorkers}
	// We'll make two channels. One for sending and receiving jobs and another for sending and receiving results.
	// Both channels will be buffered so that we don't process a huge number of jobs simultaneously.
	p.jobs = make(chan j.Job, numWorkers)
	p.results = make(chan j.Result, numWorkers)

	return p
}

// BeginWork fires off goroutines for allocating resources (thus creating jobs), for processing the jobs allocated and then finally
// collecting the results and piping it into the postProcFunc for post-processing.
func (p *Pool) BeginWork(resources []j.Resource, procFunc j.ProcessFunc, postProcFunc j.ProcessResultFunc) {
	startTime := time.Now()
	p.done = make(chan bool)
	go p.allocate(resources)
	go p.processJobs(procFunc)
	go p.collect(postProcFunc)

	// Block until our goroutines notify us that we're done.
	<-p.done
	endTime := time.Now()
	totalTime := endTime.Sub(startTime)
	log.Println("ALL PHASES DONE")
	log.Println(fmt.Sprintf("TOTAL TIME TAKEN: %dms", totalTime.Milliseconds()))
}

func (p *Pool) allocate(resources []j.Resource) {
	for i, r := range resources {
		newJob := j.Job{ID: uint(i), Resource: r}
		p.jobs <- newJob
	}

	close(p.jobs)
	log.Println("ALLOCATION PHASE DONE")
}

func (p *Pool) doWork(wg *sync.WaitGroup, proc j.ProcessFunc) {
	for job := range p.jobs {
		output, err := proc(job.Resource)
		res := j.Result{Job: job, Value: output, Err: err}

		if res.Err != nil {
			log.Println(fmt.Sprintf("Error processing job ID [%d]: %s", job.ID, res.Err.Error()))
		} else {
			log.Println(fmt.Sprintf("Processed job ID [%d] with resource: %s and output: %d", job.ID, job.Resource, res.Value))
		}

		p.results <- res
	}

	wg.Done()
}

func (p *Pool) processJobs(proc j.ProcessFunc) {
	log.Println(fmt.Sprintf("SPAWNING %d WORKERS", p.numWorkers))

	wg := sync.WaitGroup{}

	for i := 0; uint(i) < p.numWorkers; i++ {
		wg.Add(1)
		go p.doWork(&wg, proc)
	}

	wg.Wait()
	close(p.results)
	log.Println("WORK PHASE DONE")
}

func (p *Pool) collect(postProc j.ProcessResultFunc) {
	for r := range p.results {
		newRes := postProc(r)

		if newRes.Err != nil {
			log.Println(fmt.Sprintf("Error post-processing result with job ID [%d]: %s", newRes.Job.ID, newRes.Err.Error()))
		} else {
			log.Println(fmt.Sprintf("Completed job ID [%d] with resource: %s and final output: %d", newRes.Job.ID, newRes.Job.Resource, newRes.Value))
		}
	}

	log.Println("COLLECTION PHASE DONE")
	p.done <- true
}
