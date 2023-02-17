package main

import (
	"testing"
)

func TestWaitGroup(t *testing.T) {
	n := 1
	wg1 := NewWaitGroup()
	wg2 := NewWaitGroup()

	wg1.Add(n)
	wg2.Add(n)
	exited := make(chan bool)
	for i := 0; i < n; i++ {
		go func() {
			wg1.Done()
			wg2.Wait()
			exited <- true
		}()
	}
	wg1.Wait()

	for j := 0; j < n; j++ {
		select {
		case <-exited:
			t.Fatal("WaitGroup released group too soon")

		default:
		}
		wg2.Done()
	}
}
