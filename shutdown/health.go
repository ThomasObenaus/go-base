package shutdown

import (
	"fmt"
)

// IsHealthy returns an error in case a shutdown is currently in progress.
// The error is returned to indicate that the service is not healthy any more (can't handle any requests)
func (h *ShutdownHandler) IsHealthy() error {
	if h.isShutdownPending.Load() {
		return fmt.Errorf("Shutdown in progress")
	}
	return nil
}

func (h *ShutdownHandler) String() string {
	return fmt.Sprintf("ShutdownHandler (shutdown in progress=%t)", h.isShutdownPending.Load())
}
