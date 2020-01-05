package health

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewRegistry(t *testing.T) {

	// GIVEN

	// WHEN
	registry := NewRegistry()

	// THEN
	assert.NotNil(t, registry.healthChecks)
}

func Test_ShouldRegister(t *testing.T) {

	// GIVEN

	// WHEN
	registry := NewRegistry()

	// THEN
	assert.NotNil(t, registry.healthChecks)
}
