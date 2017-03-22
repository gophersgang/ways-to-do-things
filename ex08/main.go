package main

import (
	"errors"
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
	action chan func() error
}

func NewThing(msg string) *Thing {
	t := &Thing{
		msg:    msg,
		action: make(chan func() error),
	}
	go t.loop()
	return t
}

func (t *Thing) Stop() {
	t.action <- func() error {
		return errors.New("stop please")
	}
}

func (t *Thing) loop() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case fn := <-t.action:
			if err := fn(); err != nil {
				println("returning due to:", err.Error())
				return
			}
		case <-ticker.C:
			println(t.msg)
		}
	}
}

func (t *Thing) ChangeMessage(msg string) {
	t.action <- func() error {
		t.msg = msg
		return nil
	}
}

func (t *Thing) EmphasizeMessage(n int) {
	t.action <- func() error {
		t.msg += strings.Repeat("!", n)
		return nil
	}
}
