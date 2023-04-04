package main

import (
	"context"
	"log"
	"time"
)

type MultiListener interface {
	Subscribe(ctx context.Context) Listener
	Unsubscribe(ctx context.Context, l Listener)
	Broadcast(ctx context.Context, data interface{})
	RemoveAllListeners()
	HealthCheck() <-chan interface{}
}

type Broadcaster struct {
	listeners          map[Listener]struct{}
	register           chan (Listener)
	unregister         chan (Listener)
	input              chan Message
	removeAllListeners chan struct{}
	heartbeat          chan interface{}
}

type Message struct {
	ctx  context.Context
	data interface{}
}

func (b *Broadcaster) Subscribe(ctx context.Context) (Listener, error) {
	l := Listener{dataChan: make(chan interface{})}

	select {
	case <-ctx.Done():
		return Listener{}, ctx.Err()
	case b.register <- l:
		return l, nil
	}
}

func (b *Broadcaster) Unsubscribe(ctx context.Context, l Listener) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case b.unregister <- l:
		return nil
	}
}

func (b *Broadcaster) Broadcast(ctx context.Context, data interface{}) {
	m := Message{ctx, data}
	b.input <- m
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

	monitorLoop:
		for {
			select {
			case listener, ok := <-b.register:
				if !ok {
					log.Printf("Register channel has been closed, cannot add any more listeners")
					return
				}
				b.listeners[listener] = struct{}{}
				log.Printf("Listener has been added")

			case listener, ok := <-b.unregister:
				if !ok {
					log.Printf("Unregister channel has been closed, cannot remove any more listeners")
					return
				}

				delete(b.listeners, listener)
				close(listener.dataChan)
				log.Printf("Listener has been unsubscribed")

			case msg, ok := <-b.input:
				if !ok {
					log.Printf("Input channel has been closed, cannot recieve any more inputs")
					return
				}
				for listener := range b.listeners {
					select {
					case listener.dataChan <- msg.data:
					case <-msg.ctx.Done():
						log.Printf("Broadcasting has been cancelled. Stop sending")
						continue monitorLoop
					}
				}

			case t := <-pulse.C:
				sendPulse(t)

			case <-b.removeAllListeners:
				log.Printf("RemoveAllListeners functionality has been triggered. Deleting all listeners")

				for listener := range b.listeners {
					close(listener.dataChan)
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
		removeAllListeners: make(chan struct{}),
		heartbeat:          make(chan interface{}),
		input:              make(chan Message),
	}
	b.run(pulseRate)

	return &b
}
