package main

import (
	"fmt"
	"sync"
	"time"
)

type post struct {
	views int
	mu    sync.RWMutex
}

func (p *post) inc(wg *sync.WaitGroup) {

	defer func() {
		wg.Done()
		p.mu.Unlock()
	}()

	p.mu.Lock()
	p.views += 1
}
func (p *post) get(wg *sync.WaitGroup) int {

	defer func() {
		wg.Done()
		p.mu.RUnlock()
	}()

	p.mu.RLock()
	time.Sleep(time.Second * 2)
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
	fmt.Println(p.get(&wg))
}
