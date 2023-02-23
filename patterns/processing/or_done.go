package main

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
