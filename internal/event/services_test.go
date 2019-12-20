package event

import (
	"github.com/cardil/wathola/internal/config"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestProperEventsPropagation(t *testing.T) {
	// given
	throwns = make([]thrown, 0)
	stepsStore := NewStepsStore()
	finishedStore := NewFinishedStore(stepsStore)

	// when
	stepsStore.RegisterStep(&Step{Number: 1})
	stepsStore.RegisterStep(&Step{Number: 3})
	stepsStore.RegisterStep(&Step{Number: 2})
	finishedStore.RegisterFinished(&Finished{Count: 3})

	// then
	assert.Empty(t, throwns)
}

func TestMissingAndDoubleEvent(t *testing.T) {
	// given
	throwns = make([]thrown, 0)
	stepsStore := NewStepsStore()
	finishedStore := NewFinishedStore(stepsStore)

	// when
	stepsStore.RegisterStep(&Step{Number: 1})
	stepsStore.RegisterStep(&Step{Number: 2})
	stepsStore.RegisterStep(&Step{Number: 2})
	finishedStore.RegisterFinished(&Finished{Count: 3})

	// then
	assert.NotEmpty(t, throwns)
}

func TestDoubleFinished(t *testing.T) {
	// given
	throwns = make([]thrown, 0)
	stepsStore := NewStepsStore()
	finishedStore := NewFinishedStore(stepsStore)

	// when
	stepsStore.RegisterStep(&Step{Number: 1})
	stepsStore.RegisterStep(&Step{Number: 2})
	finishedStore.RegisterFinished(&Finished{Count: 2})
	finishedStore.RegisterFinished(&Finished{Count: 2})

	// then
	assert.NotEmpty(t, throwns)
}

func TestMain(m *testing.M) {
	config.Instance.Receiver.Teardown.Duration = 20 * time.Millisecond
	exitcode := m.Run()
	os.Exit(exitcode)
}
