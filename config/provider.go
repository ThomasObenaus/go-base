package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Provider is a structure containing the parsed configuration
type Provider struct {
	// config entries are all definitions of config entries that should be regarded
	configEntries []Entry

	// the environment prefix (will be added to all env vars) <envPrefix>_<name of config entry>
	// e.g. assuming the envPrefix is "myApp" and the name of the config entry is "my-entry"
	// then the env var is MYAPP_MY_ENTRY
	envPrefix string

	// instance of pflag, needed to parse command line parameters
	pFlagSet *pflag.FlagSet

	// instance of viper, needed to parse env vars and to read from cfg-file
	*viper.Viper
}

// NewProvider creates a new config provider that is able to parse the command line, env vars and config file based
// on the given entries
func NewProvider(configEntries []Entry, configName, envPrefix string) Provider {

	provider := Provider{
		configEntries: configEntries,
		envPrefix:     envPrefix,
		pFlagSet:      pflag.NewFlagSet(configName, pflag.ContinueOnError),
		Viper:         viper.New(),
	}
	return provider
}
