package main

import (
	"github.com/cardil/wathola/internal/config"
	"github.com/cardil/wathola/internal/receiver"
	"github.com/phayes/freeport"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReceiverMain(t *testing.T) {
	config.Instance.Receiver.Port = freeport.GetPort()
	go main()
	defer receiver.Stop()

	time.Sleep(time.Second)

	assert.NotNil(t, instance)
}
