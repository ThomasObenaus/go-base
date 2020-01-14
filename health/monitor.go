package health

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"
)

// Monitor represents a monitor for the health state of a service
type Monitor struct {
	// the registered health checks
	healthChecks []Check

	checkInterval time.Duration

	// in case the check evaluation has not been done within
	// checkEvaluationTimeout the health status will change to unhealthy
	checkEvaluationTimeout time.Duration
	latestCheckResult      atomic.Value

	wg sync.WaitGroup
	// channel used to signal teardown/ stop
	stopChan chan struct{}

	logger zerolog.Logger

	// will be called each time the monitor evaluates the checks
	onCheckCallback OnCheckFun
}

type checkEvaluationResult struct {
	at        time.Time
	numErrors uint
	// a map that contains one entry per check
	// if the check was healthy then the entry (error) is nil
	// if the check was NOT healthy then the entry contains the according error
	checkHealthyness map[string]error
}

// NewMonitor creates a new health monitor
func NewMonitor(options ...Option) (*Monitor, error) {
	monitor := &Monitor{
		healthChecks:           make([]Check, 0),
		checkInterval:          time.Second * 5,
		checkEvaluationTimeout: time.Second * 30,
		stopChan:               make(chan struct{}, 0),
		onCheckCallback:        nil,
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

// Join waits until the monitorhas been stopped
func (m *Monitor) Join() {
	m.wg.Wait()
}

func (m *Monitor) String() string {
	return fmt.Sprintf("HealthMonitor (%d checks)", len(m.healthChecks))
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

	for _, check := range m.healthChecks {
		name := check.String()
		err := check.IsHealthy()
		result.checkHealthyness[name] = err
		if err != nil {
			result.numErrors++
		}
		m.logger.Debug().Msgf("Check - '%s', err=%v", name, err)
	}

	if m.onCheckCallback != nil {
		healthy := true
		if result.numErrors > 0 {
			healthy = false
		}
		m.onCheckCallback(healthy, result.numErrors)
	}

	return result
}

// Register can be used to register a Check
func (m *Monitor) Register(checks ...Check) error {

	for _, check := range checks {
		if check == nil {
			return fmt.Errorf("Unable to register a check that is nil")
		}
		if len(strings.TrimSpace(check.String())) == 0 {
			return fmt.Errorf("Unable to register a check without a name")
		}
	}

	m.healthChecks = append(m.healthChecks, checks...)
	return nil
}
