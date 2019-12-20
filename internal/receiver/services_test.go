package receiver

import (
	"context"
	"fmt"
	"github.com/cardil/wathola/internal/config"
	"github.com/cardil/wathola/internal/event"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestReceiverReceive(t *testing.T) {
	// given
	e := newEvent(event.Step{Number: 42}, event.StepType)
	f := newEvent(event.Finished{Count: 1}, event.FinishedType)

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

func newEvent(data interface{}, typ string) cloudevents.Event {
	e := cloudevents.NewEvent()
	e.SetDataContentType("application/json")
	e.SetDataContentEncoding(cloudevents.Base64)
	e.SetType(typ)
	e.SetSource("test://localhost/TestReceiverReceive")
	e.SetID(randString(16))
	ensureNoError(e.SetData(data))
	ensureNoError(e.Validate())
	return e
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func randStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func randString(length int) string {
	return randStringWithCharset(length, charset)
}


func ensureNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func sendEvent(e cloudevents.Event, port int) {
	url := fmt.Sprintf("http://localhost:%v/", port)

	ht, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget(url),
		cloudevents.WithEncoding(cloudevents.HTTPBinaryV02),
	)
	ensureNoError(err)
	c, err := cloudevents.NewClient(ht)
	ensureNoError(err)
	ctx := context.TODO()
	_, _, err = c.Send(ctx, e)
	ensureNoError(err)
}
