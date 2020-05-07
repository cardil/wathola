package test

import (
	"fmt"
	"time"
)

// WaitUntil a condition is met
func WaitUntil(condition func() bool, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for {
		if condition() {
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("wait for condition exceed timeout of %v", timeout)
		}
		time.Sleep(25 * time.Millisecond)
	}
}
