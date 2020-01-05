package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewProvider(t *testing.T) {

	// GIVEN
	var configEntries []Entry
	var args []string
	configName := "testcfg"
	envPrefix := "TST"

	// WHEN
	provider := NewProvider(configEntries, configName, envPrefix)
	err := provider.ReadConfig(args)

	// THEN
	assert.NotNil(t, provider.pFlagSet)
	assert.NotNil(t, provider.Viper)
	assert.Equal(t, envPrefix, provider.envPrefix)
	assert.NoError(t, err)
	assert.Empty(t, provider.AllKeys())
}
