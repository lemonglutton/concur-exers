package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	ListeningForMessages(ctx)
}

func ListeningForMessages(c1 context.Context) {
	b := NewBroadcaster(time.Duration(3 * time.Second))
	l1 := b.Subscribe()
	l2 := b.Subscribe()
	l3 := b.Subscribe()

	c2, _ := context.WithCancel(c1)
	go func(ctx context.Context) {
		var cnt int
		t := time.NewTicker(3 * time.Second)
		for {
			select {
			case <-t.C:
				b.Broadcast(fmt.Sprintf("Random message#%d", cnt))
				cnt++
			case <-ctx.Done():
				return
			}
		}

	}(c2)

	go func(ctx context.Context, listeners ...Listener) {
		t := time.NewTicker(12 * time.Second)

		for index, listener := range listeners {
			select {
			case <-t.C:
				log.Printf("Unsubscribing Listener#%d", index+1)
				b.Unsubscribe(listener)
			case <-ctx.Done():
				return
			}
		}
	}(c2, l1, l2, l3)

	b.Broadcast("I love Golang")
	go func() {
		cpl1 := l1.Data()
		cpl2 := l2.Data()
		cpl3 := l3.Data()
		for {
			select {
			case m, ok := <-cpl1:
				if !ok {
					cpl1 = nil
				} else {
					log.Printf("Listener#1 recieved message: %v ", m)
				}

			case m, ok := <-cpl2:
				if !ok {
					cpl2 = nil
				} else {
					log.Printf("Listener#2 recieved message: %v", m)
				}

			case m, ok := <-cpl3:
				if !ok {
					cpl3 = nil
				} else {
					log.Printf("Listener#3 recieved message: %v", m)
				}

			case <-b.HealthCheck():
				log.Printf("Service is alive..")
			case <-c1.Done():
				return
			}
		}
	}()
	<-c2.Done()

}
