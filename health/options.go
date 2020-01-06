package health

import "github.com/rs/zerolog"

// Option represents an option for the Monitor
type Option func(m *Monitor)

// WithLogger specifies the logger that should be used
func WithLogger(logger zerolog.Logger) Option {
	return func(m *Monitor) {
		m.logger = logger
	}
}

// OnCheckFun called each time the monitor evaluates the checks, hence it can provide the state at this point in time
type OnCheckFun func(healthy bool, numErrors uint)

// OnCheck sets the callback that is called each time the monitor evaluates the checks
func OnCheck(fun OnCheckFun) Option {
	return func(m *Monitor) {
		m.onCheckCallback = fun
	}
}
