package solutions

import (
    "fmt"
    "sync"
)


var mutexCounter int
var mu sync.Mutex


func mutexIncrement(wg *sync.WaitGroup) {
    for i := 0; i < 1000; i++ {
        mu.Lock()
        mutexCounter += 23  // Protected by the mutex
        mu.Unlock()
    }
    wg.Done()
}


func UseMutex() {
    fmt.Println("Solving RC using Mutex")
    var wg sync.WaitGroup
    wg.Add(2)


    go mutexIncrement(&wg)
    go mutexIncrement(&wg)


    wg.Wait()
    fmt.Println("Final Counter:", mutexCounter)  // Consistent result
}
