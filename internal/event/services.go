package event

import (
	"github.com/cardil/wathola/internal/config"
	log "github.com/sirupsen/logrus"
	"time"
)

// NewRegisterStep creates RegisterStep
func NewRegisterStep() RegisterStep {
	return &stepStore{
		store: make(map[int]int),
	}
}

// NewRegisterFinished creates RegisterFinished
func NewRegisterFinished() RegisterFinished {
	return &finishedStore{
		received: 0,
		count:    -1,
	}
}

func (s stepStore) RegisterStep(step Step) {
	if times, found := s.store[step.Number]; found {
		log.Errorf(
			"event #%d received %d times, but should be received only once",
			step.Number, times+1)
	} else {
		s.store[step.Number] = 0
	}
	s.store[step.Number]++
	log.Infof("event #%d received", step.Number)
}

func (f finishedStore) RegisterFinished(finished Finished) {
	if f.received > 0 {
		log.Errorf(
			"finish event should be received only once, received %d",
			f.received+1)
	}
	f.received++
	f.count = finished.Count
	log.Infof("finish event received, expecting %d event ware propagated", finished.Count)
	d := config.Instance.Receiver.Teardown.Duration
	log.Infof("waiting additional %v to be sure all events came", d)
	time.Sleep(d)
	receivedEvents := len(f.steps.store)
	if receivedEvents != finished.Count {
		log.Errorf("expecting to have %d unique events received, " +
			"but received %d unique events", finished.Count, receivedEvents)
	} else {
		log.Infof("properly received %d unique events", receivedEvents)
	}
}

type stepStore struct {
	store map[int]int
}

type finishedStore struct {
	received int
	count    int
	steps    stepStore
}
