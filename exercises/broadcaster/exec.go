package main

import (
	"context"
	"log"
	"sync"
	"time"
)

const (
	timeoutForSubscribing  = 1000 * time.Microsecond
	timeoutForBroadcasting = 500 * time.Microsecond
	numberOfListeners      = 100
)

func main() {
	ctx := context.Background()
	b := NewBroadcaster(time.Duration(3 * time.Second))

	listeners := setupSubscribers(ctx, b)
	broadcastMessage(ctx, b, listeners, "Hello all listeners!")
}

func setupSubscribers(ctx context.Context, b *Broadcaster) []Listener {
	wg := sync.WaitGroup{}
	wg.Add(numberOfListeners)

	mu := sync.Mutex{}
	var ls []Listener
	for i := 0; i < numberOfListeners; i++ {
		go func(id int) {
			defer wg.Done()
			childCtx, cancel := context.WithTimeout(ctx, timeoutForSubscribing)
			defer cancel()

			mu.Lock()
			defer mu.Unlock()

			log.Printf("Subscribing Listener$%v", id)
			l, err := b.Subscribe(childCtx)
			if err != nil {
				log.Printf("Listener#%v was too late for audition! %v", id, err)
				return
			}
			ls = append(ls, l)
		}(i)
	}
	wg.Wait()

	return ls
}

func listenForMessages(ctx context.Context, ls []Listener) {
	wg := sync.WaitGroup{}
	wg.Add(len(ls))

	for i, l := range ls {
		go func(id int, l Listener) {
			defer wg.Done()
			select {
			case <-ctx.Done():
				log.Printf("Listener#%v gives up on listening after someone cancelled broadcasting", id)
			case msg, ok := <-l.dataChan:
				if !ok {
					log.Printf("Listeners#%v channel has been closed", id)
					return
				}
				select {
				case <-ctx.Done():
					log.Printf("Recieved a message but context was cancelled")
				default:
				}
				log.Printf("Listener#%v recived message - its: %v", id, msg)
			}
		}(i, l)
	}
	wg.Wait()
}

func broadcastMessage(ctx context.Context, b *Broadcaster, listeners []Listener, msg string) {
	childCtx, cancel := context.WithTimeout(ctx, timeoutForBroadcasting)
	defer cancel()

	b.Broadcast(childCtx, "Let's see if this message gonna get to everyone")
	listenForMessages(childCtx, listeners)
}
