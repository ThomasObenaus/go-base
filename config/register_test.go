package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_RegisterEnvParams(t *testing.T) {

	// GIVEN
	var entries []Entry
	entries = append(entries, NewEntry("test", "usage"))
	provider := NewProvider(entries, "configName", "envPrefix")

	// WHEN
	err := provider.registerEnvParams()

	// THEN
	assert.NoError(t, err)
}

func Test_RegisterEnvParamsShouldFail(t *testing.T) {

	// GIVEN
	var entries []Entry
	entries = append(entries, NewEntry("", "usage"))
	provider := NewProvider(entries, "configName", "envPrefix")

	// WHEN
	err := provider.registerEnvParams()

	// THEN
	assert.Error(t, err)

	// GIVEN
	provider = NewProvider(entries, "configName", "envPrefix", CfgFile("", ""))

	// WHEN
	err = provider.registerEnvParams()

	// THEN
	assert.Error(t, err)
}

func Test_RegisterAndParseFlags(t *testing.T) {

	// GIVEN
	var entries []Entry
	entries = append(entries, NewEntry("test1", "usage"))
	entries = append(entries, NewEntry("test2", "usage"))
	provider := NewProvider(entries, "configName", "envPrefix")
	args := []string{"--test1=A"}

	// WHEN
	err := provider.registerAndParseFlags(args)

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, "A", provider.GetString("test1"))
	assert.False(t, provider.IsSet("test2"))
}

func Test_RegisterAndParseFlags_ShouldFail(t *testing.T) {

	// GIVEN - unknown parameter
	var entries []Entry
	entries = append(entries, NewEntry("test1", "usage"))
	provider := NewProvider(entries, "configName", "envPrefix")
	args := []string{"--unkown-param=A"}

	// WHEN
	err := provider.registerAndParseFlags(args)

	// THEN
	assert.Error(t, err)
	assert.False(t, provider.IsSet("test1"))

	// GIVEN - invalid entry
	entries = append(entries, NewEntry("", "usage"))
	provider = NewProvider(entries, "configName", "envPrefix")
	args = []string{}

	// WHEN
	err = provider.registerAndParseFlags(args)

	// THEN
	assert.Error(t, err)
	assert.False(t, provider.IsSet("test1"))
}

func Test_SetDefaults(t *testing.T) {

	// GIVEN
	var entries []Entry
	entries = append(entries, NewEntry("test1", "usage", Default("2h")))
	entries = append(entries, NewEntry("test2", "usage"))
	provider := NewProvider(entries, "configName", "envPrefix")

	// WHEN
	err := provider.setDefaults()

	// THEN
	assert.NoError(t, err)
	assert.Equal(t, time.Hour*2, provider.GetDuration("test1"))
	assert.False(t, provider.IsSet("test2"))
}
