package main

import (
	"context"
	"time"
)

type MultiListener interface {
	Subscribe() <-chan interface{}
	Unsubscribe(c <-chan interface{})
	Broadcast(data interface{})
	Run(ctx context.Context)
	Kill()
}

type Broadcaster struct {
	listeners  []chan interface{}
	register   chan (chan interface{})
	unregister chan (<-chan interface{})
	input      chan interface{}
}

func (b *Broadcaster) Subscribe() <-chan interface{} {
	newListener := make(chan interface{})
	b.register <- newListener

	return newListener
}

func (b *Broadcaster) Unsubscribe(c <-chan interface{}) {
	b.unregister <- c
}

func (b *Broadcaster) Broadcast(data interface{}) {
	b.input <- data
}

func (b *Broadcaster) Run(ctx context.Context, pulseRate time.Duration) chan<- interface{} {
	heartBeat := make(chan interface{})
	pulse := time.NewTicker(pulseRate)

	go func() {
		sendPulse := func(t time.Time) {
			select {
			case heartBeat <- struct{}{}:
			default:
			}
		}
		for {
			select {
			case listener, ok := <-b.register:
				if !ok {
					return
				}
				b.listeners = append(b.listeners, listener)
			case listenerToRemove, ok := <-b.unregister:
				if !ok {
					return
				}
				for i, listener := range b.listeners {
					if listener == listenerToRemove {
						b.listeners[i] = b.listeners[len(b.listeners)-1]
						b.listeners = b.listeners[:len(b.listeners)-1]
						close(listener)
						break
					}

				}
			case val, ok := <-b.input:
				if !ok {
					return
				}

				for _, listener := range b.listeners {
					select {
					case listener <- val:
					case <-ctx.Done():
						return

					}
				}
			case <-ctx.Done():
				return
			case t := <-pulse.C:
				sendPulse(t)
			}
		}
	}()

	return heartBeat
}

func (b *Broadcaster) Kill() {
	for _, listener := range b.listeners {
		close(listener)
	}
}

func NewBroadcaster() *Broadcaster {
	b := Broadcaster{
		listeners:  nil,
		register:   make(chan (chan interface{})),
		unregister: make(chan (<-chan interface{})),
	}
	return &b
}
