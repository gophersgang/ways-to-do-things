package main

import "time"

func main() {
	t := NewThing()
	go t.Loop()
	time.Sleep(3 * time.Second)
	t.Stop()
}

type Thing struct {
	quit chan chan struct{}
}

func NewThing() *Thing {
	t := &Thing{
		quit: make(chan chan struct{}),
	}
	return t
}

func (t *Thing) Stop() {
	q := make(chan struct{})
	t.quit <- q
	<-q
}

func (t *Thing) Loop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			println("blep")
		case q := <-t.quit:
			close(q)
			println("done")
			return
		}
	}
}
