package solutions

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var atomicCounter int32

func atomicIncrement(wg *sync.WaitGroup) {
	for i := 0; i < 1000; i++ {
		atomic.AddInt32(&atomicCounter, 23) // Atomic operation
	}
	wg.Done()
}

func UseAtomicOperation() {
	fmt.Println("Solving RC using Atomic Operation")
	var wg sync.WaitGroup
	wg.Add(2)

	go atomicIncrement(&wg)
	go atomicIncrement(&wg)

	wg.Wait()
	fmt.Println("Final Counter:", atomicCounter) // Consistent result
}
