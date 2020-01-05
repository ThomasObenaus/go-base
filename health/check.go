package health

import (
	"fmt"
	"strings"
)

// Check is a simple entity that represents a health check
type Check interface {

	// IsHealthy is called to obtain the health state of the Check.
	// It should return nil if the check is healthy.
	// In case the check is not healthy the according error should be returned
	IsHealthy() error

	// Name shall return the name of the check
	Name() string
}

// CheckRegistry is the container for the health checks
type CheckRegistry struct {
	healthChecks []Check
}

// NewRegistry creates an instance that can be used to register health checks
func NewRegistry() CheckRegistry {
	registry := CheckRegistry{
		healthChecks: make([]Check, 0),
	}
	return registry
}

// Register can be used to register a Check
func (r *CheckRegistry) Register(check Check) error {

	if check == nil {
		return fmt.Errorf("Unable to register a check that is nil")
	}

	if len(strings.TrimSpace(check.Name())) == 0 {
		return fmt.Errorf("Unable to register a check without a name")
	}

	r.healthChecks = append(r.healthChecks, check)
	return nil
}
