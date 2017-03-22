package main

import (
	"sync"
	"time"
)

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
	msg  *protectedString
	quit chan chan struct{}
}

func NewThing(msg string) *Thing {
	t := &Thing{
		quit: make(chan chan struct{}),
		msg:  newProtectedString(msg),
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
		case <-ticker.C:
			println(t.msg.String())
		case q := <-t.quit:
			close(q)
			println("done")
			return
		}
	}
}

func (t *Thing) ChangeMessage(msg string) {
	t.msg.set(msg)
}

type protectedString struct {
	mtx sync.RWMutex
	val string
}

func newProtectedString(s string) *protectedString {
	return &protectedString{val: s}
}

func (ps *protectedString) String() string {
	ps.mtx.RLock()
	defer ps.mtx.RUnlock()
	return ps.val
}

func (ps *protectedString) set(s string) {
	ps.mtx.Lock()
	defer ps.mtx.Unlock()
	ps.val = s
}
