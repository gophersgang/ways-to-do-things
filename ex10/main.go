package main

import (
	"fmt"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	var g Group

	t := NewThing("blep")
	g.Add(func() error {
		t.Loop()
		return nil
	}, func(error) {
		t.Stop()
	})

	cancelEmphasizeActor := make(chan struct{})
	g.Add(func() error {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				t.EmphasizeMessage(1) // BUG!
			case <-cancelEmphasizeActor:
				return nil
			}
		}
	}, func(error) {
		close(cancelEmphasizeActor)
	})

	cancelSignalActor := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelSignalActor:
			return nil
		}
	}, func(error) {
		close(cancelSignalActor)
	})

	err := g.Run()
	fmt.Printf("returned with: %v", err)
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
