package main

import (
	"errors"
	"sync"
	"time"
)

type Contexter interface {
	Done() <-chan struct{}
}

// This is a simplified version of context struct
type Context struct {
	finish chan interface{}
	closed bool
	mu     sync.Mutex
	err    error
}

type Cancelfunc func()

func Background() Context {
	return Context{}
}

var (
	errContextWasCancelled     = errors.New("Context was cancelled")
	errContextDeadlineExceeded = errors.New("Context didn't meet deadline")
	errContextTimeoutExceeded  = errors.New("Context timeout exceeded")
)

func WithCancel(parent *Context) (*Context, Cancelfunc) {
	cancel := make(chan interface{})
	ctx := Context{cancel, false, sync.Mutex{}, nil}

	cancelFunc := func() {
		ctx.mu.Lock()
		if !ctx.closed {
			close(cancel)
			ctx.closed = true
			ctx.err = errContextWasCancelled
		}
		ctx.mu.Unlock()
	}

	go func() {
		defer func() {
			ctx.mu.Lock()
			if !ctx.closed {
				close(cancel)
				ctx.closed = true
				ctx.err = errContextWasCancelled
			}
			ctx.mu.Unlock()
		}()

		select {
		case <-parent.Done():
			return
		case <-cancel:
			return
		}
	}()

	return &ctx, cancelFunc
}

func WithDeadline(parent *Context, now time.Time, d time.Time) (*Context, Cancelfunc) {
	deadline := make(chan interface{})
	cancel := make(chan struct{})

	if now.After(d) {
		return &Context{}, nil
	}
	diff := d.Sub(now)
	ctx := Context{deadline, false, sync.Mutex{}, nil}

	cancelFunc := func() {
		ctx.mu.Lock()
		if !ctx.closed {
			close(cancel)
			close(deadline)
			ctx.closed = true
			ctx.err = errContextWasCancelled
		}
		ctx.mu.Unlock()
	}

	go func() {
		defer func() {
			ctx.mu.Lock()
			if !ctx.closed {
				close(deadline)
				close(cancel)
				ctx.closed = true
				ctx.err = errContextDeadlineExceeded
			}
			ctx.mu.Unlock()
		}()

		select {
		case <-parent.Done():
			return
		case <-time.After(diff):
		case <-cancel:
			return
		}
	}()
	return &ctx, cancelFunc
}

func (ctx *Context) Done() <-chan interface{} {
	return ctx.finish
}
func (ctx *Context) Err() error {
	return ctx.err
}

func WithTimeout(parent *Context, d time.Duration) (*Context, Cancelfunc) {
	timeout := make(chan interface{})
	cancel := make(chan interface{})
	ctx := Context{timeout, false, sync.Mutex{}, nil}

	cancelFunc := func() {
		ctx.mu.Lock()
		if !ctx.closed {
			close(timeout)
			close(cancel)
			ctx.closed = true
			ctx.err = errContextWasCancelled
		}
		ctx.mu.Unlock()
	}

	go func() {
		defer func() {
			ctx.mu.Lock()
			if !ctx.closed {
				close(timeout)
				close(cancel)
				ctx.closed = true
				ctx.err = errContextTimeoutExceeded
			}
			ctx.mu.Unlock()
		}()

		select {
		case <-parent.Done():
			return
		case <-time.After(d):
			return
		case <-cancel:
			return
		}
	}()

	return &ctx, cancelFunc
}
