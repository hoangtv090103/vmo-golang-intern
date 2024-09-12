package solutions

import (
    "fmt"
    "sync"
)

var channelCounter int


func channelIncrement(c chan int, wg *sync.WaitGroup) {
    for i := 0; i < 1000; i++ {
        c <- 23  // Send increment value through the channel
    }
    wg.Done()
}



func UseChannel() {
    fmt.Println("Solving RC using Channel")
    var wg sync.WaitGroup
    c := make(chan int)


    // Goroutine to listen on the channel and update the counter
	go func() {
		for val := range c {
			channelCounter += val
		}
	}()

    wg.Add(2)
    go channelIncrement(c, &wg)
    go channelIncrement(c, &wg)


    wg.Wait()
    close(c)  // Close the channel after all increments are sent
    
    fmt.Println("Final Counter:", channelCounter)  // Consistent result
}
