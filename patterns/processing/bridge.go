package main

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
