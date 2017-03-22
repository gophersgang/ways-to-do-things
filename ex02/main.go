package main

import "time"

func main() {
	t := NewThing()
	time.Sleep(3 * time.Second)
	t.Stop()
}

type Thing struct {
	quit chan struct{}
}

func NewThing() *Thing {
	t := &Thing{
		quit: make(chan struct{}),
	}
	go t.loop()
	return t
}

func (t *Thing) Stop() {
	close(t.quit)
}

func (t *Thing) loop() {
	tick := time.Tick(time.Second)
	for {
		select {
		case <-tick:
			println("blep")
		case <-t.quit:
			println("done")
			return
		}
	}
}
