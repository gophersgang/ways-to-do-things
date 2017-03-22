package main

import "time"

func main() {
	t := NewThing()
	cancel := make(chan struct{})
	go t.Loop(cancel)
	time.Sleep(3 * time.Second)
	close(cancel)
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

func (t *Thing) Loop(cancel <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			println("blep")
		case <-cancel:
			println("done")
			return
		}
	}
}
