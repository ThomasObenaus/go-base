package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadConfig(t *testing.T) {

	// GIVEN
	entries := []Entry{}
	provider := NewProvider(entries, "configName", "envPrefix")
	args := []string{}

	// WHEN
	err := provider.ReadConfig(args)

	// THEN
	assert.NoError(t, err)

	// GIVEN
	entries = []Entry{NewEntry("existent", "usage", Default(false))}
	args = []string{"--existent"}
	provider = NewProvider(entries, "configName", "envPrefix")

	// WHEN
	err = provider.ReadConfig(args)

	// THEN
	assert.NoError(t, err)

	// GIVEN
	entries = []Entry{}
	args = []string{"--non-existent"}
	provider = NewProvider(entries, "configName", "envPrefix")

	// WHEN
	err = provider.ReadConfig(args)

	// THEN
	assert.Error(t, err)
}
