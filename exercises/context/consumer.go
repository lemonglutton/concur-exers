package main

import "fmt"

func RunContextExample() {
	b := Background()
	parent, cancelParent := WithCancel(&b)
	child, cancelChild := WithCancel(parent)

	parentCopy := parent.Done()
	childCopy := child.Done()
	cancelChild()
	cancelParent()
	for i := 0; i < 2; i++ {
		select {
		case <-parentCopy:
			parentCopy = nil

		case <-childCopy:
			childCopy = nil
		}
	}
	fmt.Println(parent.Err(), child.Err())

}
