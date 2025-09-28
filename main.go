package main

import (
	"fmt"
	"sync"
)

type post struct {
	views int
	mu    sync.RWMutex
}

func (p *post) inc(wg *sync.WaitGroup) {
	defer wg.Done()
	p.mu.Lock()
	p.views++
	p.mu.Unlock()
}

func (p *post) get() int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.views
}

func main() {
	var wg sync.WaitGroup
	p := post{views: 0}

	for range 2 {
		wg.Add(1)
		go p.inc(&wg)
	}

	wg.Wait()
	fmt.Println(p.get())
}
