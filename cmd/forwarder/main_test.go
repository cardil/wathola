package main

import (
	"testing"
	"time"

	"github.com/cardil/wathola/internal/config"
	"github.com/cardil/wathola/internal/forwarder"
	"github.com/cardil/wathola/test"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
)

func TestForwarderMain(t *testing.T) {
	config.Instance.Forwarder.Port = freeport.GetPort()
	go main()
	defer forwarder.Stop()
	err := test.WaitUntil(forwarder.IsRunning, 10 * time.Minute)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, instance)
}
