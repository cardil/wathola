package main

import (
	"github.com/cardil/wathola/internal/config"
	"github.com/cardil/wathola/internal/receiver"
	"github.com/cardil/wathola/test"
	"github.com/phayes/freeport"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReceiverMain(t *testing.T) {
	config.Instance.Receiver.Port = freeport.GetPort()
	go main()
	defer receiver.Stop()
	err := test.WaitUntil(receiver.IsRunning, 10 * time.Minute)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, instance)
}
