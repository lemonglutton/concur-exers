package main

// The Repeat pattern is similar to Repeat-Fn, but instead of invoking a function, it generates a sequence of values.
// This pattern can generate values periodically or continuously until a "done" signal is received, as shown in the example below.
// This pattern is often used in combination with the Take pattern, which allows you to take a subset of the generated sequence each time.

func repeat(done <-chan interface{}, vals []int) <-chan interface{} {
	valueStream := make(chan interface{})

	go func() {
		defer close(valueStream)

		for {
			for _, v := range vals {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}
