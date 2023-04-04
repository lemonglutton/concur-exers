package main

// The Repeat-fn pattern is a simple generator function that is not much different from the basic generator pattern. The idea here is to provide a function that will be invoked in a forever loop.
// This pattern can be used with some periodic setting to, for example, check every some time interval if something was done, or generate some value.

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)

		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}
