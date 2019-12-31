package forwarder

import (
	"context"
	"github.com/cardil/wathola/internal/client"
	"github.com/cardil/wathola/internal/config"
	"github.com/cardil/wathola/internal/sender"
	"github.com/cloudevents/sdk-go"
	log "github.com/sirupsen/logrus"
)

// New creates new forwarder
func New() Forwarder {
	config.ReadIfPresent()
	f := &forwarder{}
	return f
}

// Stop will stop running forwarder if there is one
func Stop() {
	if cancel != nil {
		log.Info("stopping forwarder")
		cancel()
		cancel = nil
	}
}

var cancel context.CancelFunc

func (f *forwarder) Forward() {
	port := config.Instance.Forwarder.Port
	client.Receive(port, &cancel, f.forwardEvent)
}

func (f *forwarder) forwardEvent(e cloudevents.Event) {
	target := config.Instance.Forwarder.Target
	log.Infof("Forwarding event %v to %v", e.ID(), target)
	err := sender.SendEvent(e, target)
	if err != nil {
		log.Error(err)
	}
}

type forwarder struct {
}
