package main

import (
	"testing"
	"time"
)

const cancelContextErrMsg = "Context was cancelled"
const timeoutContextErrMsg = "Context timeout exceeded"
const deadlineContextErrMsg = "Context didn't meet deadline"

// // Some naive tests for the Context package implementation
func TestWithCancel(t *testing.T) {

	t.Run("Parent should close it's and child's channels when cancelation signal comes", func(t *testing.T) {
		t.Parallel()

		b := Background()
		parent, cancelParent := WithCancel(&b)
		child, _ := WithCancel(parent)

		go func() {
			<-time.After(15 * time.Second)
			t.Errorf("Test timeout..")
		}()
		cancelParent()

		parentCopy := parent.Done()
		childCopy := child.Done()
		for i := 0; i < 2; i++ {
			select {
			case <-parentCopy:
				parentCopy = nil

			case <-childCopy:
				childCopy = nil
			}
		}

		if parent.Err().Error() != cancelContextErrMsg || child.Err().Error() != cancelContextErrMsg {
			t.Errorf("Expected: parentErrMsg = \"%v\", childErrMsg = \"%v\", Actual: parentErrMsg = \"%v\" childErrMsg = \"%v\"\n", cancelContextErrMsg, cancelContextErrMsg, parent.Err().Error(), child.Err().Error())
		}
	})

	t.Run("Child shouldn't close parent's channel when its cancelation signal comes, only channel which belongs to him", func(t *testing.T) {
		t.Parallel()

		b := Background()
		parent, cancelParent := WithCancel(&b)
		defer cancelParent()

		child, cancelChild := WithCancel(parent)
		go func() {
			<-time.After(5 * time.Second)
			t.Errorf("Test timeout..")
		}()
		cancelChild()

		if parent.Err() != nil && child.Err().Error() != cancelContextErrMsg {
			t.Errorf("Expected: parentErrMsg = \"%v\", childErrMsg = \"%v\", Actual: parentErrMsg = \"%v\"0 childErrMsg = \"%v\"n", cancelContextErrMsg, cancelContextErrMsg, parent.Err().Error(), child.Err().Error())
		}
	})
}

func TestWithTimeout(t *testing.T) {
	t.Run("Writing to cancel channel after it being closed due to timeout expiration shouldn't return panic", func(t *testing.T) {
		t.Parallel()

		b := Background()
		parent, cancelParent := WithTimeout(&b, time.Duration(2*time.Second))
		child, _ := WithCancel(parent)

		go func() {
			<-time.After(5 * time.Second)
			t.Errorf("Test timeout..")
		}()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Expected: Code should not panic, Actual: Code panicked")
			}
		}()

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
	})

	t.Run("Parent should close it's and child's channels when timeout expires", func(t *testing.T) {
		t.Parallel()

		b := Background()
		parent, _ := WithTimeout(&b, time.Duration(2*time.Second))
		child, _ := WithCancel(parent)

		go func() {
			<-time.After(5 * time.Second)
			t.Errorf("Test timeout..")
		}()

		parentCopy := parent.Done()
		childCopy := child.Done()
		for i := 0; i < 2; i++ {
			select {
			case <-parentCopy:
				parentCopy = nil

			case <-childCopy:
				childCopy = nil
			}
		}

		if parent.Err().Error() != timeoutContextErrMsg || child.Err().Error() != cancelContextErrMsg {
			t.Errorf("Expected: parentErrMsg = \"%v\", childErrMsg = \"%v\", Actual: parentErrMsg = \"%v\" childErrMsg = \"%v\"\n", timeoutContextErrMsg, cancelContextErrMsg, parent.Err().Error(), child.Err().Error())
		}
	})

	t.Run("Child shouldn't close parent's channel when its timeout expires, only channel which belongs to him", func(t *testing.T) {
		t.Parallel()

		b := Background()
		parent, _ := WithCancel(&b)
		child, _ := WithTimeout(parent, time.Duration(1*time.Nanosecond))

		go func() {
			<-time.After(5 * time.Second)
			t.Errorf("Test timeout..")
		}()
		<-child.Done()

		if parent.Err() != nil && child.Err().Error() != cancelContextErrMsg {
			t.Errorf("Expected: parentErrMsg = \"%v\", childErrMsg = \"%v\", Actual: parentErrMsg = \"%v\"0 childErrMsg = \"%v\"n", cancelContextErrMsg, cancelContextErrMsg, parent.Err().Error(), child.Err().Error())
		}
	})
}

func TestWithDeadline(t *testing.T) {
	t.Run("Writing to cancel channel after it being closed due to deadline expiration shouldn't return panic", func(t *testing.T) {
		t.Parallel()

		b := Background()
		parent, cancelParent := WithDeadline(&b, time.Now().Local(), time.Now().Local().Add(time.Duration(time.Second*2)))
		child, _ := WithCancel(parent)

		go func() {
			<-time.After(5 * time.Second)
			t.Errorf("Test timeout..")
		}()

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Expected: Code should not panic, Actual: Code panicked")
			}
		}()

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
	})

	t.Run("Parent should close it's and child's channels when parent's deadline comes", func(t *testing.T) {
		t.Parallel()

		b := Background()
		parent, _ := WithDeadline(&b, time.Now().Local(), time.Now().Local().Add(time.Duration(time.Second*2)))
		child, _ := WithCancel(parent)

		go func() {
			<-time.After(5 * time.Second)
			t.Errorf("Test timeout..")
		}()

		parentCopy := parent.Done()
		childCopy := child.Done()
		for i := 0; i < 2; i++ {
			select {
			case <-parentCopy:
				parentCopy = nil

			case <-childCopy:
				childCopy = nil
			}
		}

		if parent.Err().Error() != deadlineContextErrMsg || child.Err().Error() != cancelContextErrMsg {
			t.Errorf("Expected: parentErrMsg = \"%v\", childErrMsg = \"%v\", Actual: parentErrMsg = \"%v\" childErrMsg = \"%v\"\n", deadlineContextErrMsg, cancelContextErrMsg, parent.Err().Error(), child.Err().Error())
		}
	})

	t.Run("Child shouldn't close parent's channel when deadline comes, only channel which belongs to him", func(t *testing.T) {
		t.Parallel()

		b := Background()
		parent, _ := WithCancel(&b)
		child, _ := WithDeadline(&b, time.Now().Local(), time.Now().Local().Add(time.Duration(time.Nanosecond)))

		go func() {
			<-time.After(5 * time.Second)
			t.Errorf("Test timeout..")
		}()
		<-child.Done()

		if parent.Err() != nil && child.Err().Error() != deadlineContextErrMsg {
			t.Errorf("Expected: parentErrMsg = \"%v\", childErrMsg = \"%v\", Actual: parentErrMsg = \"%v\"0 childErrMsg = \"%v\"n", cancelContextErrMsg, cancelContextErrMsg, parent.Err().Error(), child.Err().Error())
		}
	})

	t.Run("Deadline context should return Empty, not initialized context, when providing Deadline which has already passed", func(t *testing.T) {
		t.Parallel()

		b := Background()
		ctx, cancelCtx := WithDeadline(&b, time.Now().Local(), time.Date(2021, 8, 15, 14, 2, 3, 445, time.Local))

		if *ctx != (Context{}) && cancelCtx != nil {
			t.Errorf("Expected: ctx = Context{}, cancelCtx = nil, Actual: ctx = %v, cancelCtx = %v\n", ctx, cancelCtx)
		}
	})
}
