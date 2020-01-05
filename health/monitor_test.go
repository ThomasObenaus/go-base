package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewMonitor(t *testing.T) {

	// GIVEN
	registry := NewRegistry()

	// WHEN
	monitor, err := NewMonitor(&registry)

	// THEN
	assert.NoError(t, err)
	assert.NotNil(t, monitor)
	assert.NotNil(t, monitor.registry)
}

func Test_NewMonitorShouldFail(t *testing.T) {

	// GIVEN

	// WHEN
	monitor, err := NewMonitor(nil)

	// THEN
	assert.Error(t, err)
	assert.Nil(t, monitor)
}
