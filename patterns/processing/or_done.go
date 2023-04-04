package main

// The "or-done" pattern is a useful technique in Go programming for managing the lifetime of goroutines that communicate with each other through channels.
// In situations where we are unsure how a goroutine from which we are reading will handle closing a channel, this pattern provides a reliable solution.
// Without the use of a "done" channel, there is a risk of introducing a goroutine leak into our program.
// By wrapping an ambiguous channel with the "or-done" pattern, we can ensure that our goroutine will be properly closed when it has completed its work.
// This pattern also allows consumers to use a for/range statement over the channel, which can simplify the process of reading values from the channel.

func orDone(done <-chan interface{}, stream <-chan interface{}) <-chan interface{} {
	orDoneStream := make(chan interface{})

	go func() {
		defer close(orDoneStream)

		for {
			select {
			case <-done:
				return
			case v, ok := <-stream:
				if !ok {
					return
				}
				select {
				case <-done:
				case orDoneStream <- v:
				}
			}
		}
	}()
	return orDoneStream
}
