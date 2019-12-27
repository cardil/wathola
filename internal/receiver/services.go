package receiver

import (
	"context"
	"github.com/cardil/wathola/internal/config"
	"github.com/cardil/wathola/internal/event"
	cloudevents "github.com/cloudevents/sdk-go"
	log "github.com/sirupsen/logrus"
)

var cancel context.CancelFunc

// New creates new Receiver
func New() Receiver {
	config.ReadIfPresent()
	stepsStore := event.NewStepsStore()
	finishedStore := event.NewFinishedStore(stepsStore)
	r := newReceiver(stepsStore, finishedStore)
	return r
}

// Stop will stop running receiver if there is one
func Stop() {
	if cancel != nil {
		log.Info("stopping receiver")
		cancel()
		cancel = nil
	}
}

func (r receiver) Receive() {
	opt := cloudevents.WithPort(config.Instance.Receiver.Port)
	http, err := cloudevents.NewHTTPTransport(opt)
	if err != nil {
		log.Fatalf("failed to create http transport, %v", err)
	}
	c, err := cloudevents.NewClient(http)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	var ctx context.Context
	ctx, cancel = context.WithCancel(context.Background())
	log.Info("listening for events")
	err = c.StartReceiver(ctx, r.receiveEvent)
	if err != nil {
		log.Fatal(err)
	}
}

func (r receiver) receiveEvent(e cloudevents.Event) {
	// do something with event.Context and event.Data (via event.DataAs(foo)
	t := e.Context.GetType()
	if t == event.StepType {
		step := &event.Step{}
		err := e.DataAs(step)
		if err != nil {
			log.Fatal(err)
		}
		r.step.RegisterStep(step)
	}
	if t == event.FinishedType {
		finished := &event.Finished{}
		err := e.DataAs(finished)
		if err != nil {
			log.Fatal(err)
		}
		r.finished.RegisterFinished(finished)
	}
}

type receiver struct {
	step     event.StepsStore
	finished event.FinishedStore
}

func newReceiver(step event.StepsStore, finished event.FinishedStore) *receiver {
	r := &receiver{
		step:     step,
		finished: finished,
	}
	return r
}
