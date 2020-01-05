package health

import (
	"fmt"
)

// Monitor represents a monitor for the health state of a service
type Monitor struct {
	registry *CheckRegistry
}

// NewMonitor creates a new health monitor
func NewMonitor(registry *CheckRegistry) (*Monitor, error) {

	if registry == nil {
		return nil, fmt.Errorf("The CheckRegistry must not be nil")
	}

	monitor := &Monitor{
		registry: registry,
	}

	return monitor, nil
}
