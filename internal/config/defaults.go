package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Instance holds configuration values
var Instance = defaultValues()

var port = envint("PORT", 22111)

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
				Duration: 5 * time.Second,
			},
		},
		Sender: SenderConfig{
			Address:  fmt.Sprintf("http://localhost:%v/", port),
			Interval: 10 * time.Millisecond,
			Cooldown: time.Second,
		},
	}
}
