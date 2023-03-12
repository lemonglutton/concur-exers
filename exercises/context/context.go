package main

import (
	"time"
)

type Contexter interface {
	Done() <-chan struct{}
}

// This is a simplified version of context struct
type Context struct {
	finish chan interface{}
}

type Cancelfunc func()

func Background() Context {
	return Context{make(chan interface{})}
}

func WithCancel(parent Context) (Context, Cancelfunc) {
	cancel := make(chan interface{})
	ctx := Context{cancel}

	cancelFunc := func() {
		cancel <- struct{}{}
	}

	go func() {
		defer close(cancel)

		select {
		case <-parent.Done():
			return
		case <-cancel:
			return
		}
	}()

	return ctx, cancelFunc
}

func WithDeadline(parent Context, d time.Time) (Context, Cancelfunc) {
	deadline := make(chan interface{})
	cancel := make(chan struct{})

	ctx := Context{deadline}
	now := time.Now().UTC()
	diff := d.Sub(now)

	sendDeadline := func() {
		deadline <- struct{}{}
	}

	cancelFunc := func() {
		cancel <- struct{}{}
	}

	go func() {
		defer close(deadline)
		defer close(cancel)

		select {
		case <-parent.Done():
			return
		case <-time.After(diff):
			sendDeadline()
		case <-cancel:
			return
		}
	}()
	return ctx, cancelFunc
}

func (ctx *Context) Done() <-chan interface{} {
	return ctx.finish
}

func WithTimeout(parent Context, d time.Duration) (Context, Cancelfunc) {
	timeout := make(chan interface{})
	cancel := make(chan interface{})
	ctx := Context{timeout}

	cancelFunc := func() {
		cancel <- struct{}{}
	}

	go func() {
		defer close(timeout)
		defer close(cancel)

		select {
		case <-parent.Done():
			return
		case <-time.After(d):
			return
		case <-cancel:
			return
		}
	}()

	return ctx, cancelFunc
}
