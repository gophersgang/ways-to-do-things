package main

import (
	"strings"
	"time"
)

func main() {
	t := NewThing("blep")
	time.Sleep(time.Second)
	t.ChangeMessage("bork")
	time.Sleep(time.Second)
	t.EmphasizeMessage(1)
	time.Sleep(time.Second)
	t.EmphasizeMessage(3)
	time.Sleep(time.Second)
	t.Stop()
}

type Thing struct {
	msg    string
	action chan func()
	quit   chan chan struct{}
}

func NewThing(msg string) *Thing {
	t := &Thing{
		msg:    msg,
		action: make(chan func()),
		quit:   make(chan chan struct{}),
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
		case fn := <-t.action:
			fn()
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
	t.action <- func() {
		t.msg = msg
	}
}

func (t *Thing) EmphasizeMessage(n int) {
	t.action <- func() {
		t.msg += strings.Repeat("!", n)
	}
}
