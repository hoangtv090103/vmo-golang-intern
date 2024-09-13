package solutions

import (
	"fmt"
	"sync"
)

var (
	rwlCounter int
	rwMutex    sync.RWMutex
)

func RWIncrement(wg *sync.WaitGroup) {
	defer wg.Done()

	// Write lock is used because we are modifying the counter
	rwMutex.Lock()
	defer rwMutex.Unlock()

	for i := 0; i < 1000; i++ {
		rwlCounter += 23 // Now the access is synchronized
	}
}

func UseRWLock() {
	fmt.Println("Solving RC using Read/Write Lock")
	var wg sync.WaitGroup
	wg.Add(2)

	// Simulate concurrent increments
	go RWIncrement(&wg)
	go RWIncrement(&wg)

	wg.Wait()

	// Final result will now be correct and consistent
	fmt.Println("Final Counter:", rwlCounter)
}
