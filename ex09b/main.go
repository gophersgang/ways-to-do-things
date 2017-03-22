package main

import (
	_ "net/http/pprof"
	"strings"
	"time"
)

func main() {
	t := NewThing("blep")
	done := make(chan struct{})
	go func() {
		t.Loop()
		close(done)
	}()
	go func() {
		for range time.Tick(time.Second) {
			t.EmphasizeMessage(1)
		}
		t.Stop()
	}()
	<-done
}

type Thing struct {
	msg    string
	action chan func()
	quit   chan struct{}
}

func NewThing(msg string) *Thing {
	t := &Thing{
		msg:    msg,
		action: make(chan func()),
		quit:   make(chan struct{}),
	}
	return t
}

func (t *Thing) Stop() {
	close(t.quit)
}

func (t *Thing) Loop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case fn := <-t.action:
			fn()
		case <-ticker.C:
			println(t.msg)
		case <-t.quit:
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
