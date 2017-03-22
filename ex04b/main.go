package main

import (
	"context"
	"time"
)

func main() {
	t := NewThing()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	go t.Loop(ctx)
	time.Sleep(3*time.Second + 100*time.Millisecond)
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

func (t *Thing) Loop(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			println("blep")
		case <-ctx.Done():
			println("done")
			return
		}
	}
}
