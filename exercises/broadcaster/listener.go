package main

type Listener struct {
	dataChan chan interface{}
}

func (l *Listener) Data() <-chan interface{} {
	return l.dataChan
}
