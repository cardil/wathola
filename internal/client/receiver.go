package client

import (
	"context"
	cloudevents "github.com/cloudevents/sdk-go"
	log "github.com/sirupsen/logrus"
)

// Receive events and push then to passed fn
func Receive(port int, cancel *context.CancelFunc, receiveEvent func(e cloudevents.Event)) {
	opt := cloudevents.WithPort(port)
	http, err := cloudevents.NewHTTPTransport(opt)
	if err != nil {
		log.Fatalf("failed to create http transport, %v", err)
	}
	c, err := cloudevents.NewClient(http)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	ctx, ccancel := context.WithCancel(context.Background())
	cancel = &ccancel
	log.Infof("listening for events on port %v", port)
	err = c.StartReceiver(ctx, receiveEvent)
	if err != nil {
		log.Fatal(err)
	}
}
