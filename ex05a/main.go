package main

import "time"
import "sync/atomic"

func main() {
	t := NewThing("blep")
	time.Sleep(time.Second)
	t.ChangeMessage("bork")
	time.Sleep(time.Second)
	t.ChangeMessage("boop")
	time.Sleep(time.Second)
	t.ChangeMessage("angery")
	time.Sleep(time.Second)
	t.Stop()
}

type Thing struct {
	msg  atomic.Value // hmm
	quit chan chan struct{}
}

func NewThing(msg string) *Thing {
	t := &Thing{
		quit: make(chan chan struct{}),
	}
	t.msg.Store(msg) // hmmmmmm
	go t.loop()
	return t
}

func (t *Thing) Stop() {
	q := make(chan struct{})
	t.quit <- q
	<-q
}

func (t *Thing) loop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			println(t.msg.Load().(string)) // :(
		case q := <-t.quit:
			close(q)
			println("done")
			return
		}
	}
}

func (t *Thing) ChangeMessage(msg string) {
	t.msg.Store(msg) // :\
}
