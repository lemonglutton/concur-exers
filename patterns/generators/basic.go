package main

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
