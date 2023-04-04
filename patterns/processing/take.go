package main

// Take pattern is very simple pattern, which is about taking only few first values from channel.
// Once for loop finishes goroutine exites
// This pattern is often used in combination with repeat generator pattern

func take(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
	takeStream := make(chan interface{})

	go func() {
		defer close(takeStream)
		for i := 1; i <= num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}
