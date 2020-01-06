package shutdown

import "fmt"

// IsHealthy returns an error in case a shutdown is currently in progress.
// The error is returned to indicate that the service is not healthy any more (can't handle any requests)
func (h Handler) IsHealthy() error {
	if h.isShutdownPending {
		return fmt.Errorf("Shutdown in progress")
	}
	return nil
}

// Name returns the name of this handler
func (h Handler) String() string {
	return fmt.Sprintf("ShutdownHandler (shutdown in progress=%t)", h.isShutdownPending)
}
