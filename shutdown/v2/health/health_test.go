package health

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_is_healthy_when_started(t *testing.T) {
	handler := Health{}
	err := handler.IsHealthy()

	assert.NoError(t, err)
}

func Test_is_unhealthy_when_shutting_down(t *testing.T) {
	handler := Health{}
	handler.ShutdownSignalReceived()
	err := handler.IsHealthy()

	assert.Error(t, err)
}
