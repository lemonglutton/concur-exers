package main

import (
	"sync"
)

type Waiter interface {
	Add(num int)
	Done()
	Wait()
}
type WaitGroup struct {
	cnt   int
	ready sync.Cond
}

func (wg *WaitGroup) Add(num int) {
	wg.cnt += num
}

func (wg *WaitGroup) Done() {
	wg.ready.L.Lock()
	wg.cnt -= 1
	if wg.cnt == 0 {
		wg.ready.Broadcast()
	}
	wg.ready.L.Unlock()
}

func (wg *WaitGroup) Wait() {
	wg.ready.L.Lock()
	wg.ready.Wait()
	wg.ready.L.Unlock()
}

func NewWaitGroup() Waiter {
	return &WaitGroup{cnt: 0, ready: *sync.NewCond(&sync.Mutex{})}
}
