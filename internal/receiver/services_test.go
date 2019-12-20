package receiver

import (
	"fmt"
	"github.com/cardil/wathola/internal/config"
	"github.com/cardil/wathola/internal/event"
	"github.com/cardil/wathola/internal/sender"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestReceiverReceive(t *testing.T) {
	// given
	e := sender.NewCloudEvent(event.Step{Number: 42}, event.StepType)
	f := sender.NewCloudEvent(event.Finished{Count: 1}, event.FinishedType)

	instance := New()
	port := freeport.GetPort()
	config.Instance.Receiver.Port = port
	go instance.Receive()
	time.Sleep(time.Second)
	defer Stop()

	// when
	sendEvent(e, port)
	sendEvent(f, port)

	// then
	rr := instance.(*receiver)
	assert.Equal(t, 1, rr.step.Count())
	assert.Equal(t, event.Success, rr.finished.State())
}

func TestMain(m *testing.M) {
	config.Instance.Receiver.Teardown.Duration = 20 * time.Millisecond
	exitcode := m.Run()
	os.Exit(exitcode)
}

func sendEvent(e cloudevents.Event, port int) {
	url := fmt.Sprintf("http://localhost:%v/", port)
	sender.SendEvent(e, url)
}
