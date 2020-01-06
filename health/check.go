package health

// Check is a simple entity that represents a health check
type Check interface {

	// IsHealthy is called to obtain the health state of the Check.
	// It should return nil if the check is healthy.
	// In case the check is not healthy the according error should be returned
	IsHealthy() error

	// Name shall return the name of the check
	Name() string
}
