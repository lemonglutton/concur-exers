package main

import (
	"fmt"
	"time"
)

type example struct {
	data      interface{}
	timestamp int64
}

// m := NewInMemoryCache(&fifo{}, 10, nil)

// cars := []Car{
// 	Car{vinNumber: "1"},
// 	Car{vinNumber: "2"},
// 	Car{vinNumber: "3"},
// 	Car{vinNumber: "4"},
// 	Car{vinNumber: "5"},
// 	Car{vinNumber: "6"},
// 	Car{vinNumber: "7"},
// 	Car{vinNumber: "8"},
// 	Car{vinNumber: "9"},
// 	Car{vinNumber: "10"},
// 	Car{vinNumber: "11"},
// 	Car{vinNumber: "12"}}

// for _, val := range cars {
// 	go func () {

// 		m.Update(val)

// 	}
// }

// You can edit this code!
// Click here and start typing.

func main() {

	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() {
			defer close(orDone)

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}

			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):

				}
			}
		}()
		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()

		return c
	}

	start := time.Now()
	<-or(
		sig(5*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("Done after %v", time.Since(start))

}
