package main

import (
	"github.com/cardil/wathola/internal/config"
	"github.com/cardil/wathola/internal/forwarder"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestForwarderMain(t *testing.T) {
	config.Instance.Forwarder.Port = freeport.GetPort()
	go main()
	defer forwarder.Stop()

	time.Sleep(time.Second)

	assert.NotNil(t, instance)
}
