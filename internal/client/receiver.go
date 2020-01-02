package client

import (
	"context"
	"github.com/cardil/wathola/internal/config"
	"github.com/cardil/wathola/internal/ensure"
	cloudevents "github.com/cloudevents/sdk-go"
	cloudeventshttp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/http"
	log "github.com/sirupsen/logrus"
	nethttp "net/http"
)

// Receive events and push then to passed fn
func Receive(port int, cancel *context.CancelFunc, receiveEvent func(e cloudevents.Event)) {
	portOpt := cloudevents.WithPort(port)
	opts := make([]cloudeventshttp.Option, 0)
	opts = append(opts, portOpt)
	if config.Instance.Readiness.Enabled {
		readyOpt := cloudevents.WithMiddleware(readinessMiddleware)
		opts = append(opts, readyOpt)
	}
	http, err := cloudevents.NewHTTPTransport(opts...)
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

func readinessMiddleware(next nethttp.Handler) nethttp.Handler {
	log.Infof("Using readiness probe: %v", config.Instance.Readiness.URI)
	return &readinessProbe{
		next: next,
	}
}

type readinessProbe struct {
	next nethttp.Handler
}

func (r readinessProbe) ServeHTTP(rw nethttp.ResponseWriter, req *nethttp.Request) {
	if req.RequestURI == config.Instance.Readiness.URI {
		rw.WriteHeader(config.Instance.Readiness.Status)
		_, err := rw.Write([]byte(config.Instance.Readiness.Message))
		ensure.NoError(err)
		log.Debug("Received ready check")
	} else {
		r.next.ServeHTTP(rw, req)
	}
}
