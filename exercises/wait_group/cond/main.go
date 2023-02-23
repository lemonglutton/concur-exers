package main

import (
	"fmt"
	"time"
)

const loopNumber = 5

func main() {
	w1 := NewWaitGroup()

	w1.Add(loopNumber)
	for j := 0; j < loopNumber; j++ {
		go func(num int) {
			defer w1.Done()

			for i := 0; i < 5; i++ {
				fmt.Println("Start...", i)
				time.Sleep(1 * time.Second)
			}

		}(loopNumber)
	}

	w1.Wait()
	fmt.Println("Finished..!")

}
