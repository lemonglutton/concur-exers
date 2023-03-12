package main

import (
	"log"
	"sync"
	"time"
)

func RunContextExample() {
	b := Background()
	parent, cancelParent := WithTimeout(b, time.Duration(13*time.Second))
	child, _ := WithCancel(parent)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		time.Sleep(14 * time.Second)
		log.Printf("Cancel parent")
		cancelParent()
	}()

	go func(ctx Context) {
		defer wg.Done()
		<-ctx.Done()
	}(parent)

	go func(ctx Context) {
		defer wg.Done()
		<-child.Done()

	}(child)
	wg.Wait()

	// Parent cancels child via cancel function
	// Parent cancels child via timeout/deadline
	// Child cancels itself via timeout/deadline
	// Child cancels itself via cancel

}
