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
	return Context{}
}

func WithCancel(parent Context) (Context, Cancelfunc) {
	cancel := make(chan interface{})
	ctx := Context{cancel}

	cancelFunc := func() {
		close(cancel)
	}

	go func() {
		select {
		case <-parent.Done():
			return
		case <-cancel:
			return
		}
	}()

	return ctx, cancelFunc
}

func WithDeadline(parent Context, now time.Time, d time.Time) (Context, Cancelfunc) {
	deadline := make(chan interface{})
	cancel := make(chan struct{})

	if now.After(d) {
		return Context{}, nil
	}
	diff := d.Sub(now)
	ctx := Context{deadline}

	cancelFunc := func() {
		close(cancel)
	}

	go func() {
		defer close(deadline)

		select {
		case <-parent.Done():
			return
		case <-time.After(diff):
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
		close(cancel)
	}

	go func() {
		defer close(timeout)

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
