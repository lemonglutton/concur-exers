package main

// This is basic generator pattern. In this scenario it takes a collection of integers, but it can accept any time and converts this collection into stream of values.
// Generator is a pattern which is responsible for generating values, so for writing values to channel created by itself.
// That's why it's recommended to keep ownership of the channel in generator's funcion, to avoid problems with closing, writing of the channel/to the channel

func basic(done <-chan interface{}, vals []int) <-chan int {
	stream := make(chan int)

	go func() {
		defer close(stream)
		for val := range vals {
			select {
			case <-done:
				return
			case stream <- val:
			}
		}
	}()
	return stream
}
