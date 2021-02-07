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
	assert.NoError(t, err)
}

func Test_NewProviderOverrideCfgFile(t *testing.T) {

	// GIVEN
	var configEntries []Entry
	configName := "testcfg"
	envPrefix := "TST"
	cfgFileLocation := "../test/data/config.yaml"

	// WHEN
	provider := NewProvider(configEntries, configName, envPrefix, CfgFile("cfg-f", "f"))
	err := provider.ReadConfig([]string{"-f=" + cfgFileLocation})

	// THEN
	assert.NoError(t, err)
	cfgFileName := provider.GetString("cfg-f")
	assert.Equal(t, cfgFileLocation, cfgFileName)
}

func ExampleNewProvider() {
	var configEntries []Entry

	configEntries = append(configEntries, NewEntry("port", "the port to listen to", Default(8080), ShortName("p")))
	// no default value for this parameter --> thus it can be treated as a required one
	// to check if it was set by the user one can just call provider.IsSet("db-url")
	configEntries = append(configEntries, NewEntry("db-url", "the address of the data base"))
	configEntries = append(configEntries, NewEntry("db-reconnect", "enable automatic reconnect to the data base", Default(false)))

	provider := NewProvider(configEntries, "my-config", "MY_APP")
	args := []string{"--db-url=http://localhost"}

	err := provider.ReadConfig(args)
	if err != nil {
		panic(err)
	}

	port := provider.GetInt("port")
	// check for mandatory parameter
	if !provider.IsSet("db-url") {
		panic(fmt.Errorf("Parameter '--db-url' is missing"))
	}
	dbURL := provider.GetString("db-url")
	dbReconnect := provider.GetBool("db-reconnect")

	fmt.Printf("port=%d, dbURL=%s, dbReconnect=%t", port, dbURL, dbReconnect)
	// Output:
	// port=8080, dbURL=http://localhost, dbReconnect=false
}

func ExampleNewProvider_withConfigFile() {
	var configEntries []Entry

	configEntries = append(configEntries, NewEntry("port", "the port to listen to", Default(8080), ShortName("p")))

	provider := NewProvider(configEntries, "my-config", "MY_APP")
	args := []string{"--config-file=../test/data/config.yaml"}

	err := provider.ReadConfig(args)
	if err != nil {
		panic(err)
	}

	port := provider.GetInt("port")
	cfgFile := provider.GetString("config-file")

	fmt.Printf("port=%d was read from cfgFile=%s", port, cfgFile)
	// Output:
	// port=12345 was read from cfgFile=../test/data/config.yaml
}
