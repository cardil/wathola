package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	nethttp "net/http"
	"os"
	"strconv"
	"time"
)

const (
	// DefaultReceiverPort point to a default port of receiver component, and is
	// unique so that components can be easily run on localhost for easy debugging
	DefaultReceiverPort  = 22111
	// DefaultForwarderPort point to a default port of forwarder component
	DefaultForwarderPort = 22110
)

// Instance holds configuration values
var Instance = defaultValues()

var port = envint("PORT", DefaultReceiverPort)
var forwarderPort = envint("PORT", DefaultForwarderPort)

func envint(envKey string, defaultValue int) int {
	val, ok := os.LookupEnv(envKey)
	if !ok {
		return defaultValue
	}
	result, err := strconv.Atoi(val)
	if err != nil {
		panic(err)
	}
	return result
}

func defaultValues() *Config {
	return &Config{
		Receiver: ReceiverConfig{
			Port: port,
			Teardown: ReceiverTeardownConfig{
				Duration: 3 * time.Second,
			},
			Progress: ReceiverProgressConfig{
				Duration: time.Second,
			},
		},
		Forwarder: ForwarderConfig{
			Target: fmt.Sprintf("http://localhost:%v/", port),
			Port:   forwarderPort,
		},
		Sender: SenderConfig{
			Address:  fmt.Sprintf("http://localhost:%v/", forwarderPort),
			Interval: 10 * time.Millisecond,
			Cooldown: time.Second,
		},
		Readiness: ReadinessConfig{
			Enabled: true,
			URI:     "/healthz",
			Message: "OK",
			Status:  nethttp.StatusOK,
		},
		LogLevel: logrus.InfoLevel,
	}
}
