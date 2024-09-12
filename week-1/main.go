package main

import (
	"fmt"
	"sync"
	"week-1/solutions"
	// "week-1/solutions"
)

var counter int

func increment(wg *sync.WaitGroup) {
	for i := 0; i < 1000; i++ {
		counter += 23 // This is where the race condition occurs
	}
	wg.Done()
}

func RaceCondition() {
	var wg sync.WaitGroup
	wg.Add(2)

	go increment(&wg)
	go increment(&wg)

	wg.Wait()
	fmt.Println("Final Counter:", counter) // Result will be unpredictable
}

func main() {
    RaceCondition()
	solutions.UseChannel()
	solutions.UseMutex()
	solutions.UseRWLock()
	solutions.UseAtomicOperation()
	solutions.UseImmutability()
}
