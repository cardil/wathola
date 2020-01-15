package receiver

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cardil/wathola/internal/client"
	"github.com/cardil/wathola/internal/config"
	"github.com/cardil/wathola/internal/ensure"
	"github.com/cardil/wathola/internal/event"
	cloudevents "github.com/cloudevents/sdk-go"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var cancel context.CancelFunc

// New creates new Receiver
func New() Receiver {
	config.ReadIfPresent()
	errors := event.NewErrorStore()
	stepsStore := event.NewStepsStore(errors)
	finishedStore := event.NewFinishedStore(stepsStore, errors)
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
	port := config.Instance.Receiver.Port
	client.Receive(port, &cancel, r.receiveEvent, r.reportMiddleware)
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

func (r *receiver) reportMiddleware(next http.Handler) http.Handler {
	return &reportHandler{
		next:     next,
		receiver: r,
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

type reportHandler struct {
	next     http.Handler
	receiver *receiver
}

func (r reportHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.RequestURI == "/report" {
		s := r.receiver.finished.State()
		errs := r.receiver.finished.Thrown()
		events := r.receiver.step.Count()
		sj := &StateJSON{
			State:  stateToString(s),
			Events: events,
			Thrown: errs,
		}
		b, err := json.Marshal(sj)
		ensure.NoError(err)
		rw.Header().Add("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		_, err = rw.Write(b)
		ensure.NoError(err)
	} else {
		r.next.ServeHTTP(rw, req)
	}
}

func stateToString(state event.State) string {
	switch state {
	case event.Active:
		return "active"
	case event.Success:
		return "success"
	case event.Failed:
		return "failed"
	default:
		panic(fmt.Sprintf("unknown state: %v", state))
	}
}

// StateJSON represents state as JSON
type StateJSON struct {
	State  string   `json:"state"`
	Events int      `json:"events"`
	Thrown []string `json:"thrown"`
}
