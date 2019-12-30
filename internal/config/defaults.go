package config

import "time"

// Instance holds configuration values
var Instance = defaultValues()

func defaultValues() *Config {
	return &Config{
		Receiver: ReceiverConfig{
			Port: 22111,
			Teardown: ReceiverTeardownConfig{
				Duration: 5 * time.Second,
			},
		},
		Sender: SenderConfig{
			Address:  "http://localhost:22111/",
			Interval: 10 * time.Millisecond,
			Cooldown: time.Second,
		},
	}
}
