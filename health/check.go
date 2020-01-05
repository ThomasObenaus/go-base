package health

import (
	"fmt"
	"time"
)

type Check interface {
	IsHealthy() error
	Name() string
}

type CheckRegistry struct {
	healthChecks []Check

	checkInterval       time.Duration
	statusResetInterval time.Duration
}

func NewRegistry() CheckRegistry {
	registry := CheckRegistry{
		healthChecks:        make([]Check, 0),
		checkInterval:       time.Second * 5,
		statusResetInterval: time.Second * 30,
	}
	return registry
}

func (r CheckRegistry) Register(check Check) error {

	if check == nil {
		return fmt.Errorf("Unable to register a check that is nil")
	}

	r.healthChecks = append(r.healthChecks, check)
	return nil
}
