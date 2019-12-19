package receiver

import (
	"context"
	"github.com/cardil/wathola/internal/event"
	cloudevents "github.com/cloudevents/sdk-go"
	log "github.com/sirupsen/logrus"
)

// New creates new Receiver
func New() Receiver {
	r := newReceiver(
		event.NewRegisterStep(),
		event.NewRegisterFinished())
	return r
}

func (r receiver) Receive() {
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	ctx := context.Background()
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
		step := event.Step{}
		err := e.DataAs(step)
		if err != nil {
			log.Fatal(err)
		}
		r.step.RegisterStep(step)
	}
	if t == event.FinishedType {
		finished := event.Finished{}
		err := e.DataAs(finished)
		if err != nil {
			log.Fatal(err)
		}
		r.finished.RegisterFinished(finished)
	}
}

type receiver struct {
	step     event.RegisterStep
	finished event.RegisterFinished
}

func newReceiver(step event.RegisterStep, finished event.RegisterFinished) *receiver {
	r := &receiver{
		step:     step,
		finished: finished,
	}
	return r
}
