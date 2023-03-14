package main

type Listener struct {
	data chan interface{}
}

func (l *Listener) Data() <-chan interface{} {
	return l.data
}
