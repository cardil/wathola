package main

import "github.com/cardil/wathola/internal/sender"

func main() {
	sender.New().SendContinually()
}
