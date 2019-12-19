package main

import "github.com/cardil/wathola/internal/receiver"

func main() {
	r := receiver.New()
	r.Receive()
}
