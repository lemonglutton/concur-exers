package main

import (
	"time"
)

func RunContextExample() {
	b := Background()
	parent, cancelParent := WithTimeout(b, time.Duration(2*time.Second))
	child, _ := WithCancel(parent)

	parentCopy := parent.Done()
	childCopy := child.Done()
	for i := 0; i < 2; i++ {
		select {
		case <-parentCopy:
			cancelParent()
			parentCopy = nil

		case <-childCopy:
			childCopy = nil
		}
	}
}
