package main

import (
	"flag"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func factorial(n int) int {
	if n <= 1 {
		return 1
	}

	return n * factorial(n-1)
}

// Worker gets jobs from jobs channel and processes them
func worker(id int, jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	for n := range jobs {
		fmt.Printf("Worker %d: Received job %d\n", id, n)
		fmt.Printf("Worker %d: %d! = %d\n", id, n, factorial(n))
	}
}

// Run this program with the following command:
// go run worker_pool.go -jobs 10 -workers 3
func main() {

	// Parse command line arguments
	numJobs := flag.Int("jobs", 10, "Number of jobs to process")
	numWorkers := flag.Int("workers", 3, "Number of workers to create")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	// Task Queue
	// Create a channel with buffer size 10
	jobs := make(chan int, 10)

	// Worker Pool
	// Create 3 workers
	var wg sync.WaitGroup
	for w := 1; w <= *numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	// Manager Goroutine
	// Send jobs to workers
	for i := 0; i < *numJobs; i++ {
		jobs <- rand.Intn(10)
	}

	// Close the jobs channel
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All workers finished.")
}
