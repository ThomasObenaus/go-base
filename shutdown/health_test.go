package shutdown

import (
	"github.com/ThomasObenaus/go-base/stop"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_is_healthy_when_started(t *testing.T) {
	handler := ShutdownHandler{}
	err := handler.IsHealthy()

	assert.NoError(t, err)
}

func Test_is_unhealthy_when_shutting_down(t *testing.T) {
	handler := ShutdownHandler{
		registry: &stop.Registry{},
		logger:   zerolog.Nop(),
	}
	handler.ShutdownSignalReceived()
	err := handler.IsHealthy()

	assert.Error(t, err)
}

func Test_reports_health_status_depending_on_state(t *testing.T) {
	handler := ShutdownHandler{
		registry: &stop.Registry{},
		logger:   zerolog.Nop(),
	}

	status := handler.String()
	assert.Equal(t, "ShutdownHandler (shutdown in progress=false)", status)

	handler.ShutdownSignalReceived()

	status = handler.String()
	assert.Equal(t, "ShutdownHandler (shutdown in progress=true)", status)
}
