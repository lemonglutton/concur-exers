package main

import (
	"sync"
)

type Waiter interface {
	Done()
	Wait()
	Add(num int)
}

type WaitGroup struct {
	cnt int
	rdy chan struct{}
	mu  sync.Mutex
}

func (wg *WaitGroup) Add(num int) {
	wg.cnt += num
}

func (wg *WaitGroup) Wait() {
	go func() {
		defer close(wg.rdy)
		for {
			wg.mu.Lock()
			if wg.cnt == 0 {
				wg.rdy <- struct{}{}
				break
			}
			wg.mu.Unlock()
		}
	}()
	<-wg.rdy
}

func (wg *WaitGroup) Done() {
	wg.mu.Lock()
	wg.cnt -= 1
	wg.mu.Unlock()
}

func NewWaitGroup() *WaitGroup {
	return &WaitGroup{cnt: 0, rdy: make(chan struct{})}
}
