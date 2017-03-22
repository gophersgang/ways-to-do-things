package main

import "time"

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
	msg  string
	msgc chan string
	quit chan chan struct{}
}

func NewThing(msg string) *Thing {
	t := &Thing{
		msg:  msg,
		msgc: make(chan string),
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
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case t.msg = <-t.msgc:
			// thanks, friend!
		case <-ticker.C:
			println(t.msg)
		case q := <-t.quit:
			close(q)
			println("done")
			return
		}
	}
}

func (t *Thing) ChangeMessage(msg string) {
	t.msgc <- msg
}
