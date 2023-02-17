package main

import (
	"fmt"
	"testing"
)

func TestWaitGroup(t *testing.T) {
	goroutineNum := 10
	w1 := NewWaitGroup()
	w2 := NewWaitGroup()

	w1.Add(goroutineNum)
	w2.Add(goroutineNum)
	exited := make(chan bool)

	for i := 0; i < goroutineNum; i++ {
		go func() {
			w1.Done()
			w2.Wait()
			exited <- true
		}()
	}
	w1.Wait()
	fmt.Println("wg1 finished")

	for j := 0; j < goroutineNum; j++ {
		select {
		case <-exited:
			t.Fatal("WaitGroup released group too soon")

		default:
		}
		w2.Done()
	}
}
