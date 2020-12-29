package config

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Provider interface {
	ReadConfig(args []string) error
	Get(key string) interface{}
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetUint(key string) uint
	GetUint32(key string) uint32
	GetUint64(key string) uint64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
	GetIntSlice(key string) []int
	GetStringSlice(key string) []string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	GetStringMapStringSlice(key string) map[string][]string
	GetSizeInBytes(key string) uint
	IsSet(key string) bool
	String() string
}

// providerImpl is a structure containing the parsed configuration
type providerImpl struct {
	// config entries are all definitions of config entries that should be regarded
	configEntries []Entry

	// this special entry is used to specify the name + location of the config file
	configFileEntry Entry

	configName string

	// the environment prefix (will be added to all env vars) <envPrefix>_<name of config entry>
	// e.g. assuming the envPrefix is "myApp" and the name of the config entry is "my-entry"
	// then the env var is MYAPP_MY_ENTRY
	envPrefix string

	// instance of pflag, needed to parse command line parameters
	pFlagSet *pflag.FlagSet

	// instance of viper, needed to parse env vars and to read from cfg-file
	*viper.Viper
}

// ProviderOption represents an option for the Provider
type ProviderOption func(p *providerImpl)

// CfgFile specifies a default value
func CfgFile(parameterName, shortParameterName string) ProviderOption {
	return func(p *providerImpl) {
		p.configFileEntry = NewEntry(parameterName, "Specifies the full path and name of the configuration file", ShortName(shortParameterName))
	}
}

// NewProvider creates a new config provider that is able to parse the command line, env vars and config file based
// on the given entries
func NewProvider(configEntries []Entry, configName, envPrefix string, options ...ProviderOption) Provider {

	defaultConfigFileEntry := NewEntry("config-file", "Specifies the full path and name of the configuration file", Bind(true, true))
	provider := &providerImpl{
		configEntries:   configEntries,
		configName:      configName,
		envPrefix:       envPrefix,
		pFlagSet:        pflag.NewFlagSet(configName, pflag.ContinueOnError),
		Viper:           viper.New(),
		configFileEntry: defaultConfigFileEntry,
	}

	// apply the options
	for _, opt := range options {
		opt(provider)
	}

	// Enable casting to type based on given default values
	// this ensures that viper.Get() returns the casted instance instead of the plain value.
	// That helps for example when a configuration is of type time.Duration.
	// Usually viper.Get() would return a string but now it returns a time.Duration
	provider.Viper.SetTypeByDefaultValue(true)

	return provider
}

func (p *providerImpl) String() string {
	return fmt.Sprintf("%s: %v", p.configName, p.AllSettings())
}
