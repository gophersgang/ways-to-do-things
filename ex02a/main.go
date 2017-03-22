package main

import "time"

func main() {
	t := NewThing()
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
	go t.loop()
	return t
}

func (t *Thing) Stop() {
	q := make(chan struct{})
	t.quit <- q
	<-q
}

func (t *Thing) loop() {
	tick := time.Tick(time.Second)
	for {
		select {
		case <-tick:
			println("blep")
		case q := <-t.quit:
			close(q)
			println("done")
			return
		}
	}
}
