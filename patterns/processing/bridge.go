package main

// This pattern can be useful when working with multiple channels of different types and uncertain quantity, and you want to create a simple interface that aggregates them for the consumer to use.
// While this pattern is helpful, it is important to note that the order in which the channels are passed to the channel of channels matters.
// The primary channel will block until the previously received channel is fully drained before proceeding to the next one.

// assuming chanStream when closed will close our channel too, so we don't need to use orDone pattern
func bridge(done <-chan interface{}, chanStream <-chan chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})

	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case channel, ok := <-chanStream:
				if !ok {
					return
				}
				for val := range channel {
					select {
					case <-done:
						return
					case valStream <- val:
					}
				}
			}
		}
	}()

	return valStream
}

// assuming chanStream when closed will close our channel too, so we don't need to use orDone pattern
func bridgeAlternative(done <-chan interface{}, chanStream <-chan chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})

	go func() {
		defer close(valStream)
		for {
			var stream <-chan interface{}
			select {
			case <-done:
				return
			case potentialChannel, ok := <-chanStream:
				if !ok {
					return
				}
				stream = potentialChannel
			}
			for val := range stream {
				select {
				case <-done:
					return
				case valStream <- val:
				}
			}
		}
	}()

	return valStream
}
