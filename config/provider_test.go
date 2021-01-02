package config

import (
	"fmt"
	"testing"
	"time"

	"github.com/ThomasObenaus/go-base/config/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func toProviderImpl(t *testing.T, pIf interfaces.Provider) *providerImpl {
	p, ok := pIf.(*providerImpl)
	require.True(t, ok)
	require.NotNil(t, p)

	return p
}

func Test_NewProvider(t *testing.T) {

	// GIVEN
	var configEntries []Entry
	var args []string
	configName := "testcfg"
	envPrefix := "TST"

	// WHEN
	provider, err := NewProvider(configEntries, configName, envPrefix)
	require.NoError(t, err)
	err = provider.ReadConfig(args)

	// THEN
	pImpl := toProviderImpl(t, provider)
	assert.NotNil(t, pImpl.pFlagSet)
	assert.NotNil(t, pImpl.Viper)
	assert.Equal(t, envPrefix, pImpl.envPrefix)
	assert.NoError(t, err)
	assert.Len(t, pImpl.AllKeys(), 1)
	assert.Equal(t, "config-file", pImpl.configFileEntry.name)
	assert.Empty(t, pImpl.configFileEntry.flagShortName)
}

func Test_NewProviderOverrideCfgFile(t *testing.T) {

	// GIVEN
	var configEntries []Entry
	configName := "testcfg"
	envPrefix := "TST"

	// WHEN
	provider, err := NewProvider(configEntries, configName, envPrefix, CfgFile("cfg-f", "f"))
	require.NoError(t, err)

	// THEN
	pImpl := toProviderImpl(t, provider)
	assert.NotNil(t, pImpl.pFlagSet)
	assert.NotNil(t, pImpl.Viper)
	assert.Equal(t, envPrefix, pImpl.envPrefix)
	assert.Equal(t, "cfg-f", pImpl.configFileEntry.name)
	assert.Equal(t, "f", pImpl.configFileEntry.flagShortName)
}

func ExampleNewProvider() {
	var configEntries []Entry

	configEntries = append(configEntries, NewEntry("port", "the port to listen to", Default(8080), ShortName("p")))
	// no default value for this parameter --> thus it can be treated as a required one
	// to check if it was set by the user one can just call provider.IsSet("db-url")
	configEntries = append(configEntries, NewEntry("db-url", "the address of the data base"))
	configEntries = append(configEntries, NewEntry("db-reconnect", "enable automatic reconnect to the data base", Default(false)))

	provider, err := NewProvider(configEntries, "my-config", "MY_APP")
	if err != nil {
		panic(err)
	}
	args := []string{"--db-url=http://localhost"}

	err = provider.ReadConfig(args)
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

	provider, err := NewProvider(configEntries, "my-config", "MY_APP")
	if err != nil {
		panic(err)
	}

	args := []string{"--config-file=../test/data/config.yaml"}
	err = provider.ReadConfig(args)
	if err != nil {
		panic(err)
	}

	port := provider.GetInt("port")
	cfgFile := provider.GetString("config-file")

	fmt.Printf("port=%d was read from cfgFile=%s", port, cfgFile)
	// Output:
	// port=12345 was read from cfgFile=../test/data/config.yaml
}

func ExampleNewConfigProvider_primitiveTypes() {

	// The configuration with the annotations needed in order to define how the config should be filled
	type myCfg struct {
		Field1 string        `cfg:"{'name':'field-1','desc':'This is field 1','default':'default value for field 1'}"`
		Field2 int           `cfg:"{'name':'field-2','desc':'This is field 2','default':11}"`
		Field3 bool          `cfg:"{'name':'field-3','desc':'This is field 3','default':false}"`
		Field4 time.Duration `cfg:"{'name':'field-4','desc':'This is field 4','default':'5m'}"`
	}
	cfg := myCfg{}

	// Create a provider based on the given config struct
	provider, err := NewConfigProvider(&cfg, "MyConfig", "MY_APP")
	if err != nil {
		panic(err)
	}

	// As commandline arguments the parameter 'field-1' is missing, hence its default value will be used (see above)
	args := []string{"--field-2=22", "--field-3", "--field-4=10m"}

	// Read the parameters given via commandline into the config struct
	err = provider.ReadConfig(args)
	if err != nil {
		panic(err)
	}

	fmt.Printf("field-1='%s', field-2=%d, field-3=%t, field-4=%s\n", cfg.Field1, cfg.Field2, cfg.Field3, cfg.Field4)
	// Output:
	// field-1='default value for field 1', field-2=22, field-3=true, field-4=10m0s
}

func ExampleNewConfigProvider_slices() {

	// The configuration with the annotations needed in order to define how the config should be filled
	type myCfg struct {
		Field1 []string  `cfg:"{'name':'field-1','desc':'This is field 1','default':['a','b','c']}"`
		Field2 []int     `cfg:"{'name':'field-2','desc':'This is field 2','default':[1,2,3]}"`
		Field3 []bool    `cfg:"{'name':'field-3','desc':'This is field 3','default':[true,false,true]}"`
		Field4 []float64 `cfg:"{'name':'field-4','desc':'This is field 4','default':[1.1,2.2,3.3]}"`
	}
	cfg := myCfg{}

	// Create a provider based on the given config struct
	provider, err := NewConfigProvider(&cfg, "MyConfig", "MY_APP")
	if err != nil {
		panic(err)
	}

	// As commandline arguments the parameter 'field-1' is missing, hence its default value will be used (see above)
	args := []string{
		"--field-2=4,5,6",
		"--field-3=false,true",
		"--field-4=4.4,5.5,6.6",
	}

	// Read the parameters given via commandline into the config struct
	err = provider.ReadConfig(args)
	if err != nil {
		panic(err)
	}

	fmt.Printf("field-1='%v', field-2=%v, field-3=%v\n", cfg.Field1, cfg.Field2, cfg.Field3)
	// Output:
	// field-1='[a b c]', field-2=[4 5 6], field-3=[false true]
}
