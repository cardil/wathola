package main

import (
	"github.com/cardil/wathola/internal/ensure"
	"syscall"
	"testing"
	"time"
)

func TestSenderMain(t *testing.T) {
	p := syscall.Getpid()
	go main()
	time.Sleep(time.Second)
	err := syscall.Kill(p, syscall.SIGINT)
	ensure.NoError(err)
}
