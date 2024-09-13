package solutions

import (
	"fmt"
	"sync"
)

func immutableIncrement(wg *sync.WaitGroup, c chan int) {
	defer wg.Done()   // Ensure Done is called even if there's a panic
	localCounter := 0 // Each goroutine has its own counter (immutable)
	for i := 0; i < 1000; i++ {
		localCounter += 23 // Safe, no shared state
	}
	c <- localCounter // Send result to the channel
}

func UseImmutability() {
	fmt.Println("Solving RC using Immutability")
	var wg sync.WaitGroup
	c := make(chan int, 2)

	wg.Add(2)

	go immutableIncrement(&wg, c)
	go immutableIncrement(&wg, c)

	wg.Wait()
	close(c)

	finalCounter := 0
	for result := range c {
		finalCounter += result // Sum the results from both goroutines
	}

	fmt.Println("Final Counter:", finalCounter) // Consistent result
}
