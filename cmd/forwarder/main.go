package main

import "github.com/cardil/wathola/internal/forwarder"

var instance forwarder.Forwarder

func main() {
	instance = forwarder.New()
	instance.Forward()
}
