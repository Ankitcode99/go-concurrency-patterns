package main

import (
	"fmt"
	"sync"
	"time"
)

type Job func(workerID int)

type ThreadPool struct {
	wg       sync.WaitGroup
	jobQueue chan Job
}

func NewThreadPool(size int) *ThreadPool {
	pool := &ThreadPool{
		jobQueue: make(chan Job, size),
	}

	for i := 0; i < size; i++ {
		workerID := i
		pool.wg.Add(1)
		go func() {
			defer pool.wg.Done()
			for job := range pool.jobQueue {
				job(workerID)
			}
		}()
	}
	return pool
}

func (pool *ThreadPool) AddJob(job Job) {
	pool.jobQueue <- job
}

func (p *ThreadPool) Wait() {
	close(p.jobQueue)
	p.wg.Wait()
}

func main() {
	pool := NewThreadPool(30)
	start := time.Now()

	for i := 0; i < 30; i++ {
		jobID := i
		job := func(workerID int) {
			time.Sleep(1 * time.Second)
			fmt.Printf("Worker %d completed job %d\n", workerID, jobID)
		}
		pool.AddJob(job)
	}

	pool.Wait()

	fmt.Printf("Time taken %v\n", time.Since(start))
}
