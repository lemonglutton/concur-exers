package main

// fn is a function which is spawning goroutine underneath, so in array we will have channels from newly created goroutines
func fanOut(done <-chan interface{}, stream <-chan interface{}, fn func(done <-chan interface{}, stream <-chan interface{}) <-chan interface{}, workersNum int) []<-chan interface{} {
	workers := make([]<-chan interface{}, workersNum)
	for i := 0; i < workersNum; i++ {
		workers[i] = fn(done, stream)
	}

	return workers
}
