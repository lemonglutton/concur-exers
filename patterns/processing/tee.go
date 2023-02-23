package main

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
