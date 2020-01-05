package config

import (
	"fmt"
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

func ExampleNewProvider() {
	var configEntries []Entry
	configEntries = append(configEntries, NewEntry("port", "p", "the port to listen to", 8080))
	configEntries = append(configEntries, NewEntry("db-url", "u", "the address of the data base", ""))
	configEntries = append(configEntries, NewEntry("db-reconnect", "r", "enable automatic reconnect to the data base", false))

	provider := NewProvider(configEntries, "my-config", "MY_APP")
	args := []string{"-p=12000"}
	err := provider.ReadConfig(args)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", provider)

	// Output:
	// my-config: map[db-reconnect:false db-url: port:12000]
}
