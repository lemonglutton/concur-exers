package main

import (
	"context"
	"time"
)

type MultiListener interface {
	Subscribe() Listener
	Unsubscribe(c Listener)
	Broadcast(data interface{})
	Run(ctx context.Context)
	RemoveAllListeners()
	HealthCheck() <-chan interface{}
}

type Broadcaster struct {
	listeners          map[Listener]struct{}
	register           chan (Listener)
	unregister         chan (Listener)
	input              chan interface{}
	removeAllListeners chan interface{}
	heartbeat          chan interface{}
}

func (b *Broadcaster) Subscribe() Listener {
	l := Listener{data: make(chan interface{})}
	b.register <- l

	return l
}

func (b *Broadcaster) Unsubscribe(l Listener) {
	b.unregister <- l
}

func (b *Broadcaster) Broadcast(data interface{}) {
	b.input <- data
}

func (b *Broadcaster) run(pulseRate time.Duration) {
	go func() {
		pulse := time.NewTicker(pulseRate)
		sendPulse := func(t time.Time) {
			select {
			case b.heartbeat <- struct{}{}:
			default:
			}
		}

		for {
			select {
			case listener, ok := <-b.register:
				if !ok {
					return
				}
				b.listeners[listener] = struct{}{}
			case listener, ok := <-b.unregister:
				if !ok {
					return
				}

				delete(b.listeners, listener)
				close(listener.data)

			case val, ok := <-b.input:
				if !ok {
					return
				}

				for listener := range b.listeners {
					listener.data <- val
				}

			case t := <-pulse.C:
				sendPulse(t)

			case <-b.removeAllListeners:
				for listener := range b.listeners {
					close(listener.data)
				}
				b.listeners = make(map[Listener]struct{})
			}
		}
	}()
}
func (b *Broadcaster) RemoveAllListeners() {
	b.removeAllListeners <- struct{}{}
}

func (b *Broadcaster) HealthCheck() <-chan interface{} {
	return b.heartbeat
}

func NewBroadcaster(pulseRate time.Duration) *Broadcaster {
	b := Broadcaster{
		listeners:          make(map[Listener]struct{}),
		register:           make(chan Listener),
		unregister:         make(chan Listener),
		removeAllListeners: make(chan interface{}),
		heartbeat:          make(chan interface{}),
		input:              make(chan interface{}),
	}
	b.run(pulseRate)

	return &b
}
