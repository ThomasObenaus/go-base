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
	at             time.Time
	numErrors      uint
	errorsPerCheck map[string]error
}

// TODO: Add metric for health state

// NewMonitor creates a new health monitor
func NewMonitor(registry *CheckRegistry) (*Monitor, error) {

	if registry == nil {
		return nil, fmt.Errorf("The CheckRegistry must not be nil")
	}

	monitor := &Monitor{
		registry:               registry,
		checkInterval:          time.Second * 5,
		checkEvaluationTimeout: time.Second * 30,
	}

	checkResult := checkEvaluationResult{
		at:             time.Now(),
		numErrors:      0,
		errorsPerCheck: make(map[string]error),
	}
	monitor.latestCheckResult.Store(checkResult)

	return monitor, nil
}

func (m *Monitor) Start() {

	go m.monitor(m.checkInterval)
	m.logger.Info().Msg("Monitor started")
}

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
		at:             at,
		numErrors:      0,
		errorsPerCheck: make(map[string]error),
	}

	for _, check := range m.registry.healthChecks {
		err := check.IsHealthy()
		if err != nil {
			result.errorsPerCheck[check.Name()] = err
			result.numErrors++
		}
	}

	return result
}
