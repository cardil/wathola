package main

import "github.com/cardil/wathola/internal/receiver"

var instance receiver.Receiver

func main() {
	instance = receiver.New()
	instance.Receive()
}
