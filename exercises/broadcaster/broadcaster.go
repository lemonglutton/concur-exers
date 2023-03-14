package main

type MultiListener interface {
	Subscribe()
	Unsubscribe()
	Broadcast()
	Kill()
}

type Listener struct {
	dataStream chan interface{}
	doneStream chan interface{}
}

type Broadcaster struct {
	listeners  map[<-chan interface{}]Listener
	register   chan Listener
	unregister chan (<-chan interface{})
	kill       chan interface{}
	input      chan interface{}
}

func (b *Broadcaster) Subscribe() <-chan interface{} {
	listener := Listener{dataStream: make(chan interface{}), doneStream: make(chan interface{})}
	b.register <- listener

	go func(l Listener) {
		for {
			select {
			case l.dataStream <- b.input:
				return

			case <-l.doneStream:
				close(l.dataStream)
				return

			case <-b.kill:
				close(l.dataStream)
				return
			}

		}
	}(listener)

	return listener.dataStream
}

func (b *Broadcaster) Broadcast(data interface{}) {
	for _, listener := range b.listeners {
		go func(l Listener) {
			l.dataStream <- data
		}(listener)
	}
}

func (b *Broadcaster) Unsubscribe(c <-chan interface{}) {
	b.unregister <- c
}

func (b *Broadcaster) run() {
	defer b.cleanUp()
	for {
		select {
		case listener, ok := <-b.register:
			if !ok {
				return
			}
			b.listeners[listener.dataStream] = listener
		case listenerDataStream := <-b.unregister:
			listener := b.listeners[listenerDataStream]

			close(listener.dataStream)
			delete(b.listeners, listenerDataStream)
			return

		}
	}
}

func (b *Broadcaster) cleanUp() {
	close(b.register)
	close(b.unregister)
}

func (b *Broadcaster) Kill() {
	close(b.kill)
}

func NewBroadcaster() *Broadcaster {
	b := Broadcaster{
		listeners:  nil,
		register:   make(chan Listener),
		unregister: make(chan (<-chan interface{})),
		kill:       make(chan interface{}),
	}
	b.run()

	return &b
}
