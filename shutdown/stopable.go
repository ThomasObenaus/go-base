package shutdown

// Stopable is a interface for runnable sokar components
type Stopable interface {

	// Stop will be called as soon as the shutdown signal was caught.
	// Hence within this method all teardown actions should be done (e.g. free resources, leave task main loops, ...)
	Stop() error

	// String ... to meet the Stringer interface
	String() string
}
