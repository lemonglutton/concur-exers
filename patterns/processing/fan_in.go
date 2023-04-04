package main

// This pattern is useful in situations where we have a collection of channels of a certain type, and we want to multiplex them into a single channel to simplify consumption.
// Fan-in can be used to parallelize some process or stage of the pipeline. For example, we can spawn 10 goroutines to run a certain function, gather all of their channels, and multiplex them into one.
// By doing so, we are saving the consumer the headache of creating a select statement with multiple channels for the same function but different goroutines.
import "sync"

func fanIn(done <-chan interface{}, channels []<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})

	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for val := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- val:
			}

		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}
