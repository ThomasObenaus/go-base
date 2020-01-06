package health

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
)

// Monitor represents a monitor for the health state of a service
type Monitor struct {
	registry *CheckRegistry

	checkInterval time.Duration
	// in case the check evaluation has not been done within
	// checkEvaluationTimeout the health status will change to unhealthy
	checkEvaluationTimeout time.Duration
	latestCheckResult      atomic.Value

	wg sync.WaitGroup
	// channel used to signal teardown/ stop
	stopChan chan struct{}

	logger zerolog.Logger
}

type checkEvaluationResult struct {
	at        time.Time
	numErrors uint
	// a map that contains one entry per check
	// if the check was healthy then the entry (error) is nil
	// if the check was NOT healthy then the entry contains the according error
	checkHealthyness map[string]error
}

// TODO: Add metric for health state

// Option represents an option for the Monitor
type Option func(m *Monitor)

// WithLogger specifies the logger that should be used
func WithLogger(logger zerolog.Logger) Option {
	return func(m *Monitor) {
		m.logger = logger
	}
}

// NewMonitor creates a new health monitor
func NewMonitor(registry *CheckRegistry, options ...Option) (*Monitor, error) {

	if registry == nil {
		return nil, fmt.Errorf("The CheckRegistry must not be nil")
	}

	monitor := &Monitor{
		registry:               registry,
		checkInterval:          time.Second * 5,
		checkEvaluationTimeout: time.Second * 30,
		stopChan:               make(chan struct{}, 0),
	}

	checkResult := checkEvaluationResult{
		at:               time.Now(),
		numErrors:        0,
		checkHealthyness: make(map[string]error),
	}
	monitor.latestCheckResult.Store(checkResult)

	// apply the options
	for _, opt := range options {
		opt(monitor)
	}

	return monitor, nil
}

// Start starts the monitoring
func (m *Monitor) Start() {

	go m.monitor(m.checkInterval)
	m.logger.Info().Msg("Monitor started")
}

// Stop stops the monitoring
func (m *Monitor) Stop() error {
	m.logger.Info().Msg("Teardown requested")
	close(m.stopChan)
	return nil
}

func (m *Monitor) monitor(checkInterval time.Duration) {
	m.wg.Add(1)
	defer m.wg.Done()

	checkIntervalTicker := time.NewTicker(checkInterval)

	for {
		select {
		case <-m.stopChan:
			m.logger.Info().Msg("Monitor stopped")
			return
		case <-checkIntervalTicker.C:
			now := time.Now()
			latestCheckResult := m.evaluateChecks(now)
			m.latestCheckResult.Store(latestCheckResult)
		}
	}
}

func (m *Monitor) evaluateChecks(at time.Time) checkEvaluationResult {

	result := checkEvaluationResult{
		at:               at,
		numErrors:        0,
		checkHealthyness: make(map[string]error),
	}

	for _, check := range m.registry.healthChecks {
		name := check.Name()
		err := check.IsHealthy()
		result.checkHealthyness[name] = err
		if err != nil {
			result.numErrors++
		}
		m.logger.Debug().Msgf("Check %s, err=%v", name, err)
	}

	return result
}
