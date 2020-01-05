package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadCfgFile(t *testing.T) {

	// GIVEN
	configFilename := "../test/config.yaml"
	var entries []Entry
	entries = append(entries, NewEntry("test1", "usage"))
	entries = append(entries, NewEntry("test2", "usage"))
	provider := NewProvider(entries, "configName", "envPrefix")

	// WHEN
	err := provider.readCfgFile(configFilename)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, "A", provider.GetString("test1"))
	assert.False(t, provider.IsSet("test2"))
	assert.Equal(t, configFilename, provider.Viper.ConfigFileUsed())
}

func Test_ReadCfgFile_AllowNoCfgFile(t *testing.T) {

	// GIVEN
	configFilename := ""
	var entries []Entry
	entries = append(entries, NewEntry("test1", "usage"))
	provider := NewProvider(entries, "configName", "envPrefix")

	// WHEN
	err := provider.readCfgFile(configFilename)

	// THEN
	assert.NoError(t, err)
	assert.False(t, provider.IsSet("test1"))
	assert.Empty(t, provider.Viper.ConfigFileUsed())
}

func Test_ReadCfgFile_ShouldFail(t *testing.T) {

	// GIVEN
	configFilename := "does_not_exist.yaml"
	var entries []Entry
	entries = append(entries, NewEntry("test1", "usage"))
	provider := NewProvider(entries, "configName", "envPrefix")

	// WHEN
	err := provider.readCfgFile(configFilename)

	// THEN
	assert.Error(t, err)
	assert.False(t, provider.IsSet("test1"))
	assert.Equal(t, configFilename, provider.Viper.ConfigFileUsed())
}

func Test_ReadConfig_ShouldFail(t *testing.T) {

	// GIVEN
	var entries []Entry
	provider := NewProvider(entries, "configName", "envPrefix")
	args := []string{}
	provider.Viper = nil

	// WHEN
	err := provider.ReadConfig(args)

	// THEN
	assert.Error(t, err)

	// GIVEN
	provider = NewProvider(entries, "configName", "envPrefix")
	provider.pFlagSet = nil

	// WHEN
	err = provider.ReadConfig(args)

	// THEN
	assert.Error(t, err)
}
