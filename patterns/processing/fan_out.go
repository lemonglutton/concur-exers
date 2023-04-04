package main

// The Fan-out pattern involves spawning a collection of channels for a provided function.
// This pattern can be implemented in different ways, but the idea is to execute a function x times within different goroutines and return their channels.
// This pattern can simplify the process of running a function in different goroutines and is often used in combination with the Fan-in pattern.

// fn is a function which is spawning goroutine underneath, so in array we will have channels from newly created goroutines
func fanOut(done <-chan interface{}, stream <-chan interface{}, fn func(done <-chan interface{}, stream <-chan interface{}) <-chan interface{}, workersNum int) []<-chan interface{} {
	workers := make([]<-chan interface{}, workersNum)
	for i := 0; i < workersNum; i++ {
		workers[i] = fn(done, stream)
	}

	return workers
}
