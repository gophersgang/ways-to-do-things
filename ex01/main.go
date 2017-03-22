package main

import "time"

func main() {
	NewThing()
	select {}
}

type Thing struct{}

func NewThing() *Thing {
	t := &Thing{}
	go t.loop()
	return t
}

func (t *Thing) loop() {
	for range time.Tick(time.Second) {
		println("blep")
	}
}
