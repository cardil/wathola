package config

import "time"

// ReceiverTeardownConfig holds config receiver teardown
type ReceiverTeardownConfig struct {
	Duration time.Duration
}

// ReceiverConfig hold configuration for receiver
type ReceiverConfig struct {
	Teardown ReceiverTeardownConfig
	Port     int
}

// SenderConfig hold configuration for sender
type SenderConfig struct {
	Address  string
	Interval time.Duration
	Cooldown time.Duration
}

// ForwarderConfig holds configuration for forwarder
type ForwarderConfig struct {
	Target string
	Port   int
}

// Config hold complete configuration
type Config struct {
	Sender    SenderConfig
	Forwarder ForwarderConfig
	Receiver  ReceiverConfig
}
