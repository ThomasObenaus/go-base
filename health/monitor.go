package health

import (
	"fmt"
	"time"
)

// Monitor represents a monitor for the health state of a service
type Monitor struct {
	registry *CheckRegistry

	checkInterval       time.Duration
	statusResetInterval time.Duration
}

// NewMonitor creates a new health monitor
func NewMonitor(registry *CheckRegistry) (*Monitor, error) {

	if registry == nil {
		return nil, fmt.Errorf("The CheckRegistry must not be nil")
	}

	monitor := &Monitor{
		registry:            registry,
		checkInterval:       time.Second * 5,
		statusResetInterval: time.Second * 30,
	}

	return monitor, nil
}

// TODO: Add metric for health state
