package main

import (
	"fmt"
	"time"
)

const loopNumber = 5

// This is naive type of waitGroup implementation using channel and empty struct as a signal.
// This approach is not efficient due to one goroutine consuming most of cpu time for no good reason.
// This can be observed when running code. "Checking condition" is Prinln instruction which shows how often this goroutine is checking the condition
func main() {
	wg := NewWaitGroup()

	wg.Add(1)
	go func(cnt int) {
		defer wg.Done()

		for i := 0; i < cnt; i++ {
			time.Sleep(3 * time.Second)
			fmt.Printf("%v...\n", i)
		}
	}(loopNumber)

	wg.Wait()
}
