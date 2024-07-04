package health

import (
	"fmt"
	"go.uber.org/atomic"
)

type Health struct {
	isShutdownPending atomic.Bool
}

func (h *Health) ShutdownSignalReceived() {
	h.isShutdownPending.Store(true)
}

// IsHealthy returns an error in case a shutdown is currently in progress.
// The error is returned to indicate that the service is not healthy any more (can't handle any requests)
func (h *Health) IsHealthy() error {
	if h.isShutdownPending.Load() {
		return fmt.Errorf("Shutdown in progress")
	}
	return nil
}

func (h *Health) String() string {
	return fmt.Sprintf("ShutdownHandler (shutdown in progress=%t)", h.isShutdownPending.Load())
}
