package shutdown

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func Test_IsHealthyTrue(t *testing.T) {

	// GIVEN
	var logger zerolog.Logger
	h := InstallHandler(nil, logger)

	// WHEN
	err := h.IsHealthy()

	// THEN
	assert.NoError(t, err)
}

func Test_IsHealthyFalse(t *testing.T) {

	// GIVEN
	var logger zerolog.Logger
	h := InstallHandler(nil, logger)

	// WHEN
	h.isShutdownPending.Store(true)
	err := h.IsHealthy()

	// THEN
	assert.Error(t, err)
}
