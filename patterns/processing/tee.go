package main

// The "Tee" pattern in Go is a concurrency pattern that allows you to take a stream of data from a single channel and split it into two separate channels.
// This pattern is particularly useful when you need to perform multiple operations on the same data, such as filtering, processing, or aggregating.
// By splitting the data into two channels, you can perform these operations concurrently, without blocking or slowing down the original channel.

func tee(done <-chan interface{}, in <-chan interface{}) (_, _ chan interface{}) {
	out1 := make(chan interface{})
	out2 := make(chan interface{})

	go func() {
		defer func() {
			close(out1)
			close(out2)
		}()

		for val := range in {
			var copy1, copy2 = out1, out2

			for i := 0; i < 2; i++ {
				select {
				case <-done:
				case copy1 <- val:
					copy1 = nil
				case copy2 <- val:
					copy2 = nil
				}
			}
		}
	}()
	return out1, out2
}
