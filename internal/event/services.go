package event

import (
	"github.com/cardil/wathola/internal/config"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

var throwns = make([]thrown, 0)
var mutex = sync.RWMutex{}
var lastProgressReport = time.Now()

// NewStepsStore creates StepsStore
func NewStepsStore() StepsStore {
	return &stepStore{
		store: make(map[int]int),
	}
}

// NewFinishedStore creates FinishedStore
func NewFinishedStore(steps StepsStore) FinishedStore {
	return &finishedStore{
		received: 0,
		state:    Active,
		count:    -1,
		steps:    steps,
	}
}

func (s *stepStore) RegisterStep(step *Step) {
	mutex.Lock()
	if times, found := s.store[step.Number]; found {
		throw(
			"event #%d received %d times, but should be received only once",
			step.Number, times+1)
	} else {
		s.store[step.Number] = 0
	}
	s.store[step.Number]++
	mutex.Unlock()
	log.Debugf("event #%d received", step.Number)
	s.reportProgress()
}

func (s *stepStore) Count() int {
	return len(s.store)
}

func (f *finishedStore) RegisterFinished(finished *Finished) {
	if f.received > 0 {
		throw(
			"finish event should be received only once, received %d",
			f.received+1)
	}
	f.received++
	f.count = finished.Count
	log.Infof("finish event received, expecting %d event ware propagated", finished.Count)
	d := config.Instance.Receiver.Teardown.Duration
	log.Infof("waiting additional %v to be sure all events came", d)
	time.Sleep(d)
	receivedEvents := f.steps.Count()
	if receivedEvents != finished.Count {
		throw("expecting to have %v unique events received, "+
			"but received %v unique events", finished.Count, receivedEvents)
		f.reportViolations(finished)
		f.state = Failed
	} else {
		log.Infof("properly received %d unique events", receivedEvents)
		f.state = Success
	}
}

func (f *finishedStore) State() State {
	return f.state
}

func (f *finishedStore) reportViolations(finished *Finished) {
	steps := f.steps.(*stepStore)
	for eventNo := 1; eventNo <= finished.Count; eventNo++ {
		times, ok := steps.store[eventNo]
		if !ok {
			times = 0
		}
		if times != 1 {
			throw("event #%v should be received once, but was received %v times",
				eventNo, times)
		}
	}
}

func (s *stepStore) reportProgress() {
	if lastProgressReport.Add(config.Instance.Receiver.Progress.Duration).Before(time.Now()) {
		lastProgressReport = time.Now()
		log.Infof("collected %v unique events", s.Count())
	}
}

func throw(format string, args ...interface{}) {
	t := thrown{
		format: format,
		args:   args,
	}
	throwns = append(throwns, t)
	log.Errorf(t.format, t.args...)
}

type stepStore struct {
	store map[int]int
}

type finishedStore struct {
	received int
	count    int
	state    State
	steps    StepsStore
}

type thrown struct {
	format string
	args   []interface{}
}
